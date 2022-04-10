package services

import "userCrudApp/models"

type UserService interface {
	CreateUser(user *models.User) error
	GetUser(string2 *string) (*models.User, error)
	GetAll() ([]*models.User, error)
	UpdateUser(user *models.User) error
	DeleteUser(string2 *string) error
}
