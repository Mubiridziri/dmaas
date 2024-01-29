package context

import (
	"dmaas/internal/app/dmaas/config"
	"dmaas/internal/app/dmaas/dto"
	"dmaas/internal/app/dmaas/repository"
	"dmaas/internal/app/dmaas/usecase"
	"gorm.io/gorm"
)

type ApplicationContext struct {
	Config                *config.Config
	SourceChan            *chan dto.SourceChan
	SourceUseCase         usecase.SourceUseCaseInterface
	DictionaryUseCase     usecase.DictionaryUseCaseInterface
	DictionaryDataUseCase usecase.DictionaryDataUseCaseInterface
	UserUseCase           usecase.UserUseCaseInterface
	TableUseCase          usecase.TableUseCaseInterface
	TableDataUseCase      usecase.TableDataUseCaseInterface
}

func New(cfg *config.Config, db *gorm.DB, sourceChan *chan dto.SourceChan) *ApplicationContext {
	return &ApplicationContext{
		Config:     cfg,
		SourceChan: sourceChan,
		SourceUseCase: &usecase.SourceUseCase{
			DB:               db,
			SourceRepository: &repository.SourceRepository{DB: db},
		},
		DictionaryUseCase: &usecase.DictionaryUseCase{
			DictionaryRepository: &repository.DictionaryRepository{DB: db},
		},
		DictionaryDataUseCase: &usecase.DictionaryDataUseCase{
			DictionaryDataRepository: &repository.DictionaryDataRepository{DB: db},
		},
		UserUseCase: &usecase.UserUseCase{
			UserRepository: &repository.UserRepository{DB: db},
		},
		TableUseCase: &usecase.TableUseCase{
			TableRepository: &repository.TableRepository{DB: db},
		},
		TableDataUseCase: &usecase.TableDataUseCase{
			TableDataRepository: &repository.TableDataRepository{DB: db},
		},
	}
}
