// Обработчики GET запросов
package handlers

import (
	"encoding/json"
	"net/http"

	"github.com/BillyBones007/loyalty-service/internal/db/models"
	"github.com/BillyBones007/loyalty-service/internal/tools/convert"
	"github.com/BillyBones007/loyalty-service/internal/tools/jwttoken"
)

// For bad requests
func (h *Handler) BadHandler(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusBadRequest)
	rw.Write([]byte("Hey! Motherfucker! You know, you don't belong here! Request incorrect!"))
}

// Get a list with information on user orders
func (h *Handler) GetOrdersInfoHandler(rw http.ResponseWriter, r *http.Request) {
	token := r.Context().Value(Tkn).(*jwttoken.CurrentToken)
	if token.Err != nil {
		http.Error(rw, token.Err.Error(), http.StatusUnauthorized)
		return
	}
	userID := token.ClaimsToken["user"]
	listOrders, err := h.Storage.Order().GetOrdersInfo(userID)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusNoContent)
		return
	}
	listOrders = convert.TimeInRC3339ForListOrders(listOrders)
	resp, err := json.Marshal(listOrders)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(http.StatusOK)
	rw.Write(resp)
}

// Get the current balance and withdrawal of user
func (h *Handler) GetCurrentBalanceHandler(rw http.ResponseWriter, r *http.Request) {
	token := r.Context().Value(Tkn).(*jwttoken.CurrentToken)
	if token.Err != nil {
		http.Error(rw, token.Err.Error(), http.StatusUnauthorized)
		return
	}
	userID := token.ClaimsToken["user"]
	currBalance := models.CurrentBalance{}
	err := h.Storage.Order().GetCurrentBalance(&currBalance, userID)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	currBalance.Current = convert.ConvToFloatBalance(currBalance.IntCurrent)
	currBalance.Withdrawn = convert.ConvToFloatBalance(currBalance.IntWithdrawn)

	resp, err := json.Marshal(currBalance)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(http.StatusOK)
	rw.Write(resp)
}

// Get a list with information for withdrawals user
func (h *Handler) GetWithdrawalsHandler(rw http.ResponseWriter, r *http.Request) {
	token := r.Context().Value(Tkn).(*jwttoken.CurrentToken)
	if token.Err != nil {
		http.Error(rw, token.Err.Error(), http.StatusUnauthorized)
		return
	}
	userID := token.ClaimsToken["user"]
	listWithdrawals, err := h.Storage.Order().GetWithdrawalsInfo(userID)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusNoContent)
		return
	}
	listWithdrawals = convert.TimeInRC3339ForWithdrawals(listWithdrawals)
	resp, err := json.Marshal(listWithdrawals)
	if err != nil {
		http.Error(rw, err.Error(), http.StatusInternalServerError)
		return
	}
	rw.Header().Set("Content-Type", "application/json; charset=utf-8")
	rw.WriteHeader(http.StatusOK)
	rw.Write(resp)
}
