package dictionaries

import (
	"dmaas/internal/dto"
	"dmaas/internal/entity"
	"dmaas/internal/repository"
)

type DictionaryUseCaseInterface interface {
	CreateDictionary(dictionary *entity.Dictionary) error
	UpdateDictionary(dictionary *entity.Dictionary) error
	RemoveDictionary(dictionary *entity.Dictionary) error
	ListDictionaries(pagination dto.Query) ([]entity.Dictionary, error)
	GetDictionaryById(id int) (entity.Dictionary, error)
	GetCount() int64
}

type DictionaryUseCase struct {
	DictionaryRepository repository.DictionaryRepositoryInterface
}

func (useCase DictionaryUseCase) CreateDictionary(dictionary *entity.Dictionary) error {
	return useCase.DictionaryRepository.CreateDictionary(dictionary)
}

func (useCase DictionaryUseCase) UpdateDictionary(dictionary *entity.Dictionary) error {
	return useCase.UpdateDictionary(dictionary)
}

func (useCase DictionaryUseCase) RemoveDictionary(dictionary *entity.Dictionary) error {
	return useCase.DictionaryRepository.RemoveDictionary(dictionary)
}

func (useCase DictionaryUseCase) ListDictionaries(pagination dto.Query) ([]entity.Dictionary, error) {
	return useCase.DictionaryRepository.ListDictionaries(pagination)
}

func (useCase DictionaryUseCase) GetDictionaryById(id int) (entity.Dictionary, error) {
	return useCase.DictionaryRepository.GetDictionaryById(id)
}

func (useCase DictionaryUseCase) GetCount() int64 {
	return useCase.DictionaryRepository.GetCount()
}
