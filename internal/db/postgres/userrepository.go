package postgres

import (
	"context"
	"errors"
	"log"

	"github.com/BillyBones007/loyalty-service/internal/customerr"
	"github.com/BillyBones007/loyalty-service/internal/db/models"
	"github.com/BillyBones007/loyalty-service/internal/tools/convert"
)

// Repository for work to user
type UserRepository struct {
	store *Storage
}

// Create new user
func (u *UserRepository) Create(model *models.User) error {
	// implement encrypting the password before INSERT request
	if err := model.EncryptPassword(); err != nil {
		return err
	}
	tx, err := u.store.Pool.Begin(context.TODO())
	if err != nil {
		return err
	}
	defer func() {
		if err != nil {
			tx.Rollback(context.TODO())
		}
		tx.Commit(context.TODO())
	}()
	qUsers := "INSERT INTO users (uuid, login, encrypted_password) VALUES (uuid_generate_v4(), $1, $2) RETURNING uuid;"
	if err := tx.QueryRow(context.Background(), qUsers, model.Login, model.EncryptedPassword).Scan(&model.UUID); err != nil {
		log.Printf("error from Create function: %s\n", err)
		return err
	}
	qBalace := "INSERT INTO balance (user_id) VALUES ($1);"
	_, err = tx.Exec(context.Background(), qBalace, model.UUID)
	if err != nil {
		log.Printf("error from Create function: %s\n", err)
	}

	return nil
}

// User is exist?
func (u *UserRepository) UserIsExists(model *models.User) bool {
	var flag bool
	q := "SELECT EXISTS(SELECT login FROM users WHERE login = $1);"
	if err := u.store.Pool.QueryRow(context.Background(), q, model.Login).Scan(&flag); err != nil {
		log.Printf("error in function UserIsExists: %s\n", err)
		return flag
	}
	return flag
}

// Find by login
func (u *UserRepository) FindByLogin(model *models.User) error {
	var uuid [16]byte
	var encrPass string
	q := "SELECT uuid, encrypted_password FROM users WHERE login = $1;"
	if err := u.store.Pool.QueryRow(context.Background(), q, model.Login).Scan(&uuid, &encrPass); err != nil {
		if errors.Is(err, customerr.ErrNoRows) {
			return customerr.ErrLoginOrPassIncorrect
		}
		return err
	}
	model.UUID = convert.UUID(uuid).String()
	model.EncryptedPassword = encrPass
	return nil
}
