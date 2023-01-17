package usecase

import "github.com/diyliv/storage/internal/storage"

type storageUC struct {
	postgresRepo storage.PostgresRepository
}

func NewStorageUC(postgresRepo storage.PostgresRepository) *storageUC {
	return &storageUC{postgresRepo: postgresRepo}
}
