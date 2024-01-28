package driver

import "dmaas/internal/app/dmaas/entity"

type DriverInterface interface {
	ImportDatabase(source entity.Source, localSchemaName string) error
	DropForeignServer(source entity.Source) error
}
