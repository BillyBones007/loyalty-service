// Обработчики POST запросов
package handlers

import (
	"encoding/json"
	"errors"
	"io"
	"net/http"

	"github.com/BillyBones007/loyalty-service/internal/customerr"
	"github.com/BillyBones007/loyalty-service/internal/db/models"
	"github.com/BillyBones007/loyalty-service/internal/tools/convert"
	"github.com/BillyBones007/loyalty-service/internal/tools/jwttoken"
	"github.com/BillyBones007/loyalty-service/internal/tools/luhn"
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
		return
	}
	// check correct the request
	if len(user.Login) == 0 || len(user.Password) == 0 {
		http.Error(rw, customerr.ErrFieldEmpty.Error(), http.StatusBadRequest)
		return
	}
	// if login is busy UserIsExists returns true
	if h.Storage.User().UserIsExists(user) {
		http.Error(rw, customerr.ErrLoginIsExists.Error(), http.StatusConflict)
		return
	}
	err = h.Storage.User().Create(user)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	// set cookie with JWT
	tokenStr, err := jwttoken.GetTokenString([]byte(h.Key), user.UUID)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	cookie := &http.Cookie{
		Name:  "token",
		Value: tokenStr,
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
		return
	}
	// check correct the request
	if len(user.Login) == 0 || len(user.Password) == 0 {
		http.Error(rw, customerr.ErrFieldEmpty.Error(), http.StatusBadRequest)
		return
	}

	err = h.Storage.User().FindByLogin(user)
	if err != nil {
		http.Error(rw, customerr.ErrLoginOrPassIncorrect.Error(), http.StatusUnauthorized)
		return
	}
	if !user.ComparePassword() {
		http.Error(rw, customerr.ErrLoginOrPassIncorrect.Error(), http.StatusUnauthorized)
		return
	}

	// set cookie with JWT
	tokenStr, err := jwttoken.GetTokenString([]byte(h.Key), user.UUID)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	cookie := &http.Cookie{
		Name:  "token",
		Value: tokenStr,
	}
	http.SetCookie(rw, cookie)
	rw.WriteHeader(http.StatusOK)
}

// Uploading the numder order for calculation
func (h *Handler) UploadNumberOrderHandler(rw http.ResponseWriter, r *http.Request) {
	token := r.Context().Value(Tkn).(*jwttoken.CurrentToken)
	if token.Err != nil {
		http.Error(rw, token.Err.Error(), http.StatusUnauthorized)
		return
	}
	userID := token.ClaimsToken["user"]
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	numOrder := string(body)
	if !luhn.LuhnValid(numOrder) {
		http.Error(rw, customerr.ErrInvalidNumOrder.Error(), http.StatusUnprocessableEntity)
		return
	}
	id := h.Storage.Order().CheckOrder(numOrder)
	if id != "" {
		if id == userID {
			rw.WriteHeader(http.StatusOK)
			return
		}
		http.Error(rw, customerr.ErrOrderIsExists.Error(), http.StatusConflict)
		return
	}
	newOrder := models.OrderInfo{Order: numOrder, Status: models.New}
	err = h.Storage.Order().InsertNewOrder(newOrder, userID)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.WriteHeader(http.StatusAccepted)
}

// Withdrawal points
func (h *Handler) WithdrawPointsHandler(rw http.ResponseWriter, r *http.Request) {
	token := r.Context().Value(Tkn).(*jwttoken.CurrentToken)
	if token.Err != nil {
		http.Error(rw, token.Err.Error(), http.StatusUnauthorized)
		return
	}
	userID := token.ClaimsToken["user"]
	body, err := io.ReadAll(r.Body)
	defer r.Body.Close()
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	modelWithdraw := models.BalanceWithdraw{}
	err = json.Unmarshal(body, &modelWithdraw)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusBadRequest)
		return
	}
	// check order number
	id := h.Storage.Order().CheckOrder(modelWithdraw.Order)
	if id != "" {
		http.Error(rw, customerr.ErrInvalidNumOrder.Error(), http.StatusUnprocessableEntity)
		return
	}
	// try insert order withdraw
	modelWithdraw.IntSum = convert.ConvToIntBalance(modelWithdraw.Sum)
	err = h.Storage.Order().WithdrawalPoints(modelWithdraw, userID)
	if err != nil {
		if errors.Is(err, customerr.ErrInsufficientFounds) {
			http.Error(rw, err.Error(), http.StatusPaymentRequired)
			return
		}
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}

	rw.WriteHeader(http.StatusOK)
}
