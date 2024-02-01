package context

import (
	"dmaas/internal/config"
	"dmaas/internal/dto"
	"dmaas/internal/repository"
	"dmaas/internal/usecase/dictionaries"
	"dmaas/internal/usecase/sources"
	"dmaas/internal/usecase/users"
	"gorm.io/gorm"
)

type ApplicationContext struct {
	Config                *config.Config
	SourceSender          chan dto.SourceMessage
	SourceUseCase         sources.SourceUseCaseInterface
	DictionaryUseCase     dictionaries.DictionaryUseCaseInterface
	DictionaryDataUseCase dictionaries.DictionaryDataUseCaseInterface
	UserUseCase           users.UserUseCaseInterface
	TableUseCase          sources.TableUseCaseInterface
	TableDataUseCase      sources.TableDataUseCaseInterface
}

func New(cfg *config.Config, db *gorm.DB, sourceSender chan dto.SourceMessage) *ApplicationContext {
	return &ApplicationContext{
		Config:       cfg,
		SourceSender: sourceSender,
		SourceUseCase: &sources.SourceUseCase{
			DB:               db,
			SourceRepository: &repository.SourceRepository{DB: db},
			SourceSender:     sourceSender,
		},
		DictionaryUseCase: &dictionaries.DictionaryUseCase{
			DictionaryRepository: &repository.DictionaryRepository{DB: db},
		},
		DictionaryDataUseCase: &dictionaries.DictionaryDataUseCase{
			DictionaryDataRepository: &repository.DictionaryDataRepository{DB: db},
		},
		UserUseCase: &users.UserUseCase{
			UserRepository: &repository.UserRepository{DB: db},
		},
		TableUseCase: &sources.TableUseCase{
			TableRepository: &repository.TableRepository{DB: db},
		},
		TableDataUseCase: &sources.TableDataUseCase{
			TableDataRepository: &repository.TableDataRepository{DB: db},
		},
	}
}
