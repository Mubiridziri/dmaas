package users

import (
	"dmaas/internal/dto"
	"dmaas/internal/entity"
	"dmaas/internal/repository"
)

type UserUseCaseInterface interface {
	CreateUser(user *entity.User) error
	UpdateUser(user *entity.User) error
	RemoveUser(user *entity.User) error
	ListUsers(pagination dto.Query) ([]entity.User, error)
	GetUserById(id int) (entity.User, error)
	GetUserByUsername(username string) (entity.User, error)
	GetCount() int64
}

type UserUseCase struct {
	UserRepository repository.UserRepositoryInterface
}

func (u UserUseCase) CreateUser(user *entity.User) error {
	return u.UserRepository.CreateUser(user)
}

func (u UserUseCase) UpdateUser(user *entity.User) error {
	return u.UserRepository.UpdateUser(user)
}

func (u UserUseCase) RemoveUser(user *entity.User) error {
	return u.UserRepository.RemoveUser(user)
}

func (u UserUseCase) ListUsers(pagination dto.Query) ([]entity.User, error) {
	return u.UserRepository.ListUsers(pagination)
}

func (u UserUseCase) GetUserById(id int) (entity.User, error) {
	return u.UserRepository.GetUserById(id)
}

func (u UserUseCase) GetUserByUsername(username string) (entity.User, error) {
	return u.UserRepository.GetUserByUsername(username)
}

func (u UserUseCase) GetCount() int64 {
	return u.UserRepository.GetCount()
}
