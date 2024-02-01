package dto

import "dmaas/internal/entity"

const (
	ImportDatabaseAction = "ImportDatabaseAction"
	RemoveDatabaseAction = "RemoveDatabaseAction"
)

type SourceMessage struct {
	Source *entity.Source
	Action string
}
