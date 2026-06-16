package persistencedb

import (
	"folderstructure/internal/model/db"
	"folderstructure/platform/logger"

	"github.com/jackc/pgx/v4/pgxpool"
	"go.mongodb.org/mongo-driver/mongo"
)

type PersistenceDB struct {
	*db.Queries
	Mongo   *mongo.Client
	MongoDB *mongo.Database
	pool    *pgxpool.Pool
	log     logger.Logger
}

func New(pool *pgxpool.Pool, mongoClient *mongo.Client, mongoDB *mongo.Database, log logger.Logger) PersistenceDB {
	return PersistenceDB{
		Queries: db.New(pool),
		Mongo:   mongoClient,
		MongoDB: mongoDB,
		pool:    pool,
		log:     log,
	}
}
