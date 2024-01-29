package dto

import "dmaas/internal/app/dmaas/entity"

type UserRequest struct {
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
}

type UserUpdateRequest struct {
	Name     string `json:"name" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password"`
}

func (dto *UserRequest) ToUser() entity.User {
	return entity.User{
		Name:     dto.Name,
		Username: dto.Username,
		Password: dto.Password,
	}
}

func (dto *UserUpdateRequest) ToUser(user *entity.User) {
	user.Name = dto.Name
	user.Username = dto.Username
	user.Password = dto.Password

}
