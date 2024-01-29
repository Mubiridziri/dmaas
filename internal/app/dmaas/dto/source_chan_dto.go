package dto

import "dmaas/internal/app/dmaas/entity"

const (
	ImportDatabaseAction = "ImportDatabaseAction"
	RemoveDatabaseAction = "RemoveDatabaseAction"
)

type SourceChan struct {
	Source entity.Source
	Action string
}
