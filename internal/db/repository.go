package db

import "github.com/BillyBones007/loyalty-service/internal/db/models"

type UniversalUserRepository interface {
	Create(model *models.User) error
	UserIsExists(model *models.User) bool
	FindByLogin(model *models.User) error
}

type UniversalOrderRepository interface {
}
