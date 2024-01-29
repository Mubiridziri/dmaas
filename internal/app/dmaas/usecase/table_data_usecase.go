package usecase

import (
	"dmaas/internal/app/dmaas/dto"
	"dmaas/internal/app/dmaas/entity"
	"dmaas/internal/app/dmaas/repository"
	"fmt"
)

type TableDataUseCaseInterface interface {
	ListTableData(source entity.Source, table entity.Table, pagination dto.Query) ([]map[string]interface{}, error)
	GetCount(source entity.Source, table entity.Table) int64
}

type TableDataUseCase struct {
	TableDataRepository repository.TableDataRepositoryInterface
}

func (t TableDataUseCase) ListTableData(source entity.Source, table entity.Table, pagination dto.Query) ([]map[string]interface{}, error) {
	localSchemaName := t.getLocalSchemaName(source)
	return t.TableDataRepository.ListTableData(localSchemaName, table, pagination)
}

func (t TableDataUseCase) GetCount(source entity.Source, table entity.Table) int64 {
	localSchemaName := t.getLocalSchemaName(source)
	return t.TableDataRepository.GetCount(localSchemaName, table)
}

// getLocalSchemaName BAD! ALREADY EXISTS IN SourceUseCase.getLocalSchemaName
func (t TableDataUseCase) getLocalSchemaName(source entity.Source) string {
	return fmt.Sprintf("import_schema_%v", source.ID)
}
