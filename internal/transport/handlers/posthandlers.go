// Обработчики POST запросов
package handlers

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/BillyBones007/loyalty-service/internal/customerr"
	"github.com/BillyBones007/loyalty-service/internal/db/models"
)

// Registration new user
func (h *Handler) RegisterUserHandler(rw http.ResponseWriter, r *http.Request) {
	user := models.NewUser()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	err = json.Unmarshal(body, user)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}
	// check correct the request
	if len(user.Login) == 0 || len(user.Password) == 0 {
		http.Error(rw, customerr.ErrFieldEmpty.Error(), http.StatusBadRequest)
	}
	// if login is busy UserIsExists returns true
	if h.Storage.User().UserIsExists(user) {
		http.Error(rw, customerr.ErrLoginIsExists.Error(), http.StatusConflict)
	}
	err = h.Storage.User().Create(user)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
	}

	// set cookie with user.UUID
	cookie := &http.Cookie{
		Name:  "uuid",
		Value: user.UUID,
	}
	http.SetCookie(rw, cookie)
	rw.WriteHeader(http.StatusOK)
}

// Authentication user
func (h *Handler) AuthUserHandler(rw http.ResponseWriter, r *http.Request) {
	user := models.NewUser()
	body, err := io.ReadAll(r.Body)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	defer r.Body.Close()
	err = json.Unmarshal(body, user)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
	}
	// check correct the request
	if len(user.Login) == 0 || len(user.Password) == 0 {
		http.Error(rw, customerr.ErrFieldEmpty.Error(), http.StatusBadRequest)
	}

	h.Storage.User().FindByLogin(user)

	// set cookie with user.UUID
	cookie := &http.Cookie{
		Name:  "uuid",
		Value: user.UUID,
	}
	http.SetCookie(rw, cookie)
	rw.WriteHeader(http.StatusOK)
}

// Uploading the numder order for calculation
func (h *Handler) UploadNumberOrderHandler(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("This is Post Upload Number!"))
}

// Withdrow points
func (h *Handler) WithdrawPointsHandler(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("This is Post Withdraw Points!"))
}
