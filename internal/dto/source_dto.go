package dto

import (
	"dmaas/internal/entity"
	"errors"
	"strings"
)

type SourceRequest struct {
	Title    string `json:"title" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Type     string `json:"type" binding:"required"`
	Host     string `json:"host" binding:"required"`
	Port     int    `json:"port" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Schema   string `json:"schema" binding:"required"`
}

type SourceUpdateRequest struct {
	Title    string `json:"title" binding:"required"`
	Name     string `json:"name" binding:"required"`
	Host     string `json:"host" binding:"required"`
	Port     int    `json:"port" binding:"required"`
	Username string `json:"username" binding:"required"`
	Password string `json:"password" binding:"required"`
	Schema   string `json:"schema" binding:"required"`
}

func (dto *SourceRequest) Validate() error {
	if strings.Contains(dto.Name, ";") {
		return errors.New("name contains is forbidden symbol")
	}

	if strings.Contains(dto.Host, ";") {
		return errors.New("host contains is forbidden symbol")
	}
	if strings.Contains(dto.Username, ";") {
		return errors.New("username contains is forbidden symbol")
	}
	if strings.Contains(dto.Password, ";") {
		return errors.New("password contains is forbidden symbol")
	}
	if strings.Contains(dto.Schema, ";") {
		return errors.New("schema contains is forbidden symbol")
	}

	return nil
}

// Validate как переиспользовать этот код вместо копирования?
func (dto *SourceUpdateRequest) Validate() error {
	if strings.Contains(dto.Name, ";") {
		return errors.New("name contains is forbidden symbol")
	}

	if strings.Contains(dto.Host, ";") {
		return errors.New("host contains is forbidden symbol")
	}
	if strings.Contains(dto.Username, ";") {
		return errors.New("username contains is forbidden symbol")
	}
	if strings.Contains(dto.Password, ";") {
		return errors.New("password contains is forbidden symbol")
	}
	if strings.Contains(dto.Schema, ";") {
		return errors.New("schema contains is forbidden symbol")
	}

	return nil
}

func (dto *SourceRequest) ToSource() entity.Source {
	return entity.Source{
		Title:    dto.Title,
		Name:     dto.Name,
		Type:     dto.Type,
		Host:     dto.Host,
		Port:     dto.Port,
		Username: dto.Username,
		Password: dto.Password,
		Schema:   dto.Schema,
	}
}

func (dto *SourceUpdateRequest) ToSource(source *entity.Source) {
	source.Title = dto.Title
	source.Name = dto.Name
	source.Host = dto.Host
	source.Port = dto.Port
	source.Username = dto.Username
	source.Password = dto.Password
	source.Schema = dto.Schema
}
