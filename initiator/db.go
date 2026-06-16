package initiator

import (
	"context"
	"fmt"
	"time"

	"folderstructure/platform/logger"

	"github.com/jackc/pgx/v4/pgxpool"
	"github.com/spf13/viper"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"go.uber.org/zap"
)

func initDatabases(log logger.Logger) (*pgxpool.Pool, *mongo.Client, *mongo.Database) {
	pgPool := initPostgres("db.url", "_", log)

	mongoClient, mongoDB := initMongo(log)

	return pgPool, mongoClient, mongoDB
}

func initPostgres(configKey, dbName string, log logger.Logger) *pgxpool.Pool {
	url := viper.GetString(configKey)
	if url == "" {
		log.Fatal(context.Background(), "database url is empty")
	}

	config, err := pgxpool.ParseConfig(url)
	if err != nil {
		log.Fatal(context.Background(), "unable to parse postgres config")
	}

	config.MaxConnIdleTime = viper.GetDuration("database.idle_conn_timeout")
	if config.MaxConnIdleTime == 0 {
		config.MaxConnIdleTime = 4 * time.Minute
	}

	pool, err := pgxpool.ConnectConfig(context.Background(), config)
	if err != nil {
		log.Fatal(context.Background(), fmt.Sprintf("failed to connect to postgres (%s)", dbName))
	}

	log.Info(context.Background(), "PostgreSQL connected successfully")
	return pool
}

func initMongo(log logger.Logger) (*mongo.Client, *mongo.Database) {
	mongoURI := viper.GetString("mongo.url")
	if mongoURI == "" {
		log.Fatal(context.Background(), "mongo.url is not configured")
	}

	poolSize := viper.GetUint64("mongo.pool_size")
	if poolSize == 0 {
		poolSize = 50
	}

	clientOptions := options.Client().
		ApplyURI(mongoURI).
		SetMaxPoolSize(poolSize)

	client, err := mongo.Connect(context.Background(), clientOptions)
	if err != nil {
		log.Fatal(context.Background(), "failed to connect to MongoDB")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	if err := client.Ping(ctx, nil); err != nil {
		log.Fatal(context.Background(), "MongoDB ping failed")
	}

	dbName := viper.GetString("mongo.database")
	if dbName == "" {
		dbName = "ussd"
	}

	log.Info(context.Background(), "MongoDB connected successfully",
		zap.String("database", dbName))

	return client, client.Database(dbName)
}
