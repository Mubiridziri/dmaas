package repository

import (
	"dmaas/internal/app/dmaas/dto"
	"dmaas/internal/app/dmaas/entity"
	"gorm.io/gorm"
)

type UserRepositoryInterface interface {
	CreateUser(user *entity.User) error
	UpdateUser(user *entity.User) error
	RemoveUser(user *entity.User) error
	ListUsers(pagination dto.Query) ([]entity.User, error)
	GetUserById(id int) (entity.User, error)
	GetUserByUsername(username string) (entity.User, error)
	GetCount() int64
}

type UserRepository struct {
	DB *gorm.DB
}

func (repository *UserRepository) GetUserById(id int) (entity.User, error) {
	var user entity.User
	if err := repository.DB.Where(entity.User{ID: id}).First(&user).Error; err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (repository *UserRepository) GetUserByUsername(username string) (entity.User, error) {
	var user entity.User
	if err := repository.DB.Where(entity.User{Username: username}).First(&user).Error; err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func (repository *UserRepository) CreateUser(user *entity.User) error {
	return repository.DB.Create(user).Error
}

func (repository *UserRepository) UpdateUser(user *entity.User) error {
	return repository.DB.Save(user).Error
}

func (repository *UserRepository) RemoveUser(user *entity.User) error {
	return repository.DB.Delete(user).Error
}

func (repository *UserRepository) ListUsers(pagination dto.Query) ([]entity.User, error) {
	var users []entity.User
	offset := (pagination.Page - 1) * pagination.Limit
	if err := repository.DB.Offset(offset).Limit(pagination.Limit).Find(&users).Error; err != nil {
		return []entity.User{}, err
	}

	return users, nil
}

func (repository *UserRepository) GetCount() int64 {
	var count int64
	repository.DB.Model(&entity.User{}).Count(&count)
	return count
}
