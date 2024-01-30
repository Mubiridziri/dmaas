package dto

import "dmaas/internal/app/dmaas/entity"

const (
	ImportDatabaseAction = "ImportDatabaseAction"
	RemoveDatabaseAction = "RemoveDatabaseAction"
)

type SourceMessage struct {
	Source *entity.Source
	Action string
}
