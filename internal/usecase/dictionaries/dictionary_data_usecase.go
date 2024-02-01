package dictionaries

import (
	"dmaas/internal/dto"
	"dmaas/internal/entity"
	"dmaas/internal/repository"
)

type DictionaryDataUseCaseInterface interface {
	ListDictionaryData(dictionary entity.Dictionary, pagination dto.Query) ([]entity.DictionaryRow, error)
	GetCount(dictionary entity.Dictionary) int64
}

type DictionaryDataUseCase struct {
	DictionaryDataRepository repository.DictionaryDataRepositoryInterface
}

func (d DictionaryDataUseCase) ListDictionaryData(dictionary entity.Dictionary, pagination dto.Query) ([]entity.DictionaryRow, error) {
	return d.DictionaryDataRepository.ListDictionaryData(dictionary, pagination)
}

func (d DictionaryDataUseCase) GetCount(dictionary entity.Dictionary) int64 {
	return d.DictionaryDataRepository.GetCount(dictionary)
}
