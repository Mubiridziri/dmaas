package entity

import "gorm.io/gorm"

type Repository struct {
	db *gorm.DB
	userRepository
	sourceRepository
	tableRepository
	tableDataRepository
	dictionaryRepository
	dictionaryDataRepository
	// ...
}

func NewRepository(db *gorm.DB) Repository {
	return Repository{
		db: db,
		userRepository: userRepository{
			db: db,
		},
		sourceRepository: sourceRepository{
			db: db,
		},
		tableRepository: tableRepository{
			db: db,
		},
		tableDataRepository: tableDataRepository{
			db: db,
		},
		dictionaryRepository: dictionaryRepository{
			db: db,
		},
		dictionaryDataRepository: dictionaryDataRepository{
			db: db,
		},
	}
}
