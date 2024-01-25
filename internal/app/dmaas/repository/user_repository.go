package repository

import (
	"dmaas/internal/app/dmaas/database"
	"dmaas/internal/app/dmaas/entity"
)

func FindUserByID(id int) (entity.User, error) {
	var user entity.User
	if err := database.DB.Where(entity.User{ID: id}).First(&user).Error; err != nil {
		return entity.User{}, err
	}

	return user, nil
}

func FindUserByUsername(username string) (entity.User, error) {
	var user entity.User
	if err := database.DB.Where(entity.User{Username: username}).First(&user).Error; err != nil {
		return entity.User{}, err
	}

	return user, nil
}
