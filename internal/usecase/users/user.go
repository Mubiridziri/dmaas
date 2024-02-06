package users

import (
	"dmaas/internal/entity"
	"time"
)

type UserView struct {
	Name      string    `json:"name"`
	Username  string    `json:"username"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
}

type CreateOrUpdateUserView struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Password string `json:"password"`
}

type Repository interface {
	CreateUser(user *entity.User) error
	UpdateUser(user *entity.User) error
	RemoveUser(user *entity.User) error
	ListUsers(page, limit int) ([]entity.User, error)
	GetUserById(id int) (entity.User, error)
	GetUserByUsername(username string) (entity.User, error)
	GetUsersCount() int64
}

type PaginatedUsersList struct {
	Total   int64      `json:"total"`
	Entries []UserView `json:"entries"`
}

type Controller struct {
	Repository
}

func NewController(repo Repository) *Controller {
	return &Controller{Repository: repo}
}

func (c Controller) CreateUser(input CreateOrUpdateUserView) (UserView, error) {
	user := entity.User{
		Name:     input.Name,
		Username: input.Username,
		Password: input.Password,
	}
	if err := c.Repository.CreateUser(&user); err != nil {
		return UserView{}, err
	}

	return fromDBUser(&user), nil
}

func (c Controller) UpdateUser(input CreateOrUpdateUserView) (UserView, error) {
	user := entity.User{
		Name:     input.Name,
		Username: input.Username,
	}

	if input.Password != "" {
		user.Password = input.Password
	}

	if err := c.Repository.UpdateUser(&user); err != nil {
		return UserView{}, err
	}

	return fromDBUser(&user), nil
}

func (c Controller) RemoveUser(id int) (UserView, error) {

	user, err := c.Repository.GetUserById(id)

	if err != nil {
		return UserView{}, err
	}

	if err := c.Repository.RemoveUser(&user); err != nil {
		return UserView{}, err
	}

	return fromDBUser(&user), nil
}

func (c Controller) ListUsers(page, limit int) (PaginatedUsersList, error) {

	users, err := c.Repository.ListUsers(page, limit)

	if err != nil {
		return PaginatedUsersList{}, err
	}

	var userViews []UserView

	for _, user := range users {
		userViews = append(userViews, fromDBUser(&user))
	}

	return PaginatedUsersList{
		Total:   c.Repository.GetUsersCount(),
		Entries: userViews,
	}, nil
}

func (c Controller) GetUserById(id int) (UserView, error) {

	user, err := c.Repository.GetUserById(id)

	if err != nil {
		return UserView{}, err
	}

	return fromDBUser(&user), nil
}

func (c Controller) GetUserByUsername(username string) (UserView, error) {
	user, err := c.Repository.GetUserByUsername(username)

	if err != nil {
		return UserView{}, err
	}

	return fromDBUser(&user), nil
}

func fromDBUser(user *entity.User) UserView {
	return UserView{
		Name:      user.Name,
		Username:  user.Username,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
	}
}
