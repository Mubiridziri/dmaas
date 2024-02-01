package sources

import (
	"dmaas/internal/dto"
	"dmaas/internal/entity"
	"dmaas/internal/repository"
)

type TableUseCaseInterface interface {
	ListTables(sourceId int, pagination dto.Query) ([]entity.Table, error)
	GetTableById(id int) (entity.Table, error)
	GetCount() int64
}

type TableUseCase struct {
	TableRepository repository.TableRepositoryInterface
}

func (t TableUseCase) ListTables(sourceId int, pagination dto.Query) ([]entity.Table, error) {
	return t.TableRepository.ListTables(sourceId, pagination)
}

func (t TableUseCase) GetTableById(id int) (entity.Table, error) {
	return t.TableRepository.GetTableById(id)
}

func (t TableUseCase) GetCount() int64 {
	return t.TableRepository.GetCount()
}
