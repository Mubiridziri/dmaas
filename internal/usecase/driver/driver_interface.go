package driver

import "dmaas/internal/entity"

type DriverInterface interface {
	ImportDatabase(source entity.Source, localSchemaName string) error
	DropForeignServer(source entity.Source) error
}
