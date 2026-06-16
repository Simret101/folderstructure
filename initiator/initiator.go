package initiator

import (
	"context"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"folderstructure/docs"
	_ "folderstructure/docs"
	"folderstructure/internal/handler/middleware"
	"folderstructure/internal/model/persistencedb"
	"folderstructure/platform/workerpool"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
	"go.uber.org/zap"
)

func Initiate() {
	ctx := context.Background()

	log, err := zap.NewProduction()
	if err != nil {
		log.Fatal("unable to start logger")
	}

	log.Info("initializing config ")
	configName := "config"
	if os.Getenv("CONFIG_NAME") != "" {
		configName = os.Getenv("CONFIG_NAME")
	}
	languageConfig := "language"

	err = InitConfig(Config{
		Names:  []string{configName, languageConfig},
		Path:   "config",
		Logger: log,
	})

	if err != nil {
		log.Fatal("unable to start config", zap.Error(err))
	}

	log.Info("initializing config completed")
	docs.SwaggerInfo.Title = "Loans API"
	docs.SwaggerInfo.Description = "API documentation for Loans"
	docs.SwaggerInfo.Version = "1.0"
	docs.SwaggerInfo.BasePath = "/"

	logger := InitLogger()
	log.Info("initializing databases")
	pgxPool, mongoClient, mongoDB := initDatabases(logger)

	log.Info("initializing persistence layer ")
	persistenceDB := persistencedb.New(pgxPool, mongoClient, mongoDB, logger)

	redisPool := initRedis(viper.GetString("redis.url"), logger)
	log.Info("redis connection initialized")

	logger.Info(ctx, "fayda initialized")

	wp := workerpool.New(viper.GetInt("workerpool.max_workers"), viper.GetInt("workerpool.task_buffer"))
	wp.Start()

	persistence := initPersistence(&persistenceDB, logger, redisPool)
	logger.Info(ctx, "done initializing persistence layer")

	logger.Info(ctx, "initializing module layer")
	


	service := initService( persistence, logger, wp)
	logger.Info(ctx, "done initializing service layer")

	logger.Info(ctx, "initializing transport layer ")
	transport := initTransport(service, logger)
	logger.Info(ctx, "done initializing transport layer")

	logger.Info(ctx, "initializing http server")
	server := gin.New()
	server.Use(middleware.GinLogger(logger))
	server.Use(middleware.CORS())
	server.Use(middleware.ErrorHandler())
	swaggerUser := viper.GetString("swagger.username")
	swaggerPass := viper.GetString("swagger.password")

	swagger := server.Group("/swagger")

	swagger.Use(middleware.SwaggerBasicAuth(swaggerUser, swaggerPass))
	swagger.GET("/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	logger.Info(ctx, "done initializing server")
	api := server.Group("")

	crypto, err := middleware.NewCryptoMiddleware(viper.GetString("crypto.key"))
	if err != nil {
		logger.Fatal(context.Background(), "invalid crypto key config", zap.Error(err))
	}
	api.Use(crypto.Handler())

	logger.Info(ctx, "initializing route")
	initRoute(api, transport, service, logger, viper.GetString("jwt.secret"))
	logger.Info(ctx, "done initializing route")

	logger.Info(ctx, "initializing server")
	srv := &http.Server{
		Addr:              fmt.Sprintf("%s:%d", viper.GetString("app.host"), viper.GetInt("app.port")),
		Handler:           server,
		ReadHeaderTimeout: viper.GetDuration("app.timeout"),
		IdleTimeout:       30 * time.Minute,
	}
	go func() {
		sigint := make(chan os.Signal, 1)
		signal.Notify(sigint, syscall.SIGINT)
		<-sigint
		log.Fatal("HTTP server Shutdown")
	}()

	logger.Info(ctx, fmt.Sprintf("http server listening on port : %d", viper.GetInt("app.port")))

	err = srv.ListenAndServe()
	if err != nil {
		logger.Fatal(ctx, fmt.Sprintf("Could not start HTTP server: %s", err))
	}
}
