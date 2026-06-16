package initiator

import (
	"folderstructure/platform/logger"

	"folderstructure/platform/workerpool"
)

type Service struct {
}

func initService(persistence *Persistence, log logger.Logger, pool *workerpool.WorkerPool) *Service {

	return &Service{}
}
