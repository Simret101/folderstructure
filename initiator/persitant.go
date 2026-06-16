package initiator

import (
	"folderstructure/internal/model/persistencedb"

	"folderstructure/platform/logger"

	"github.com/go-redis/redis/v8"
)

type Persistence struct {
}

func initPersistence(persistencdb *persistencedb.PersistenceDB, log logger.Logger, redisClient *redis.Client) *Persistence {

	return &Persistence{}
}
