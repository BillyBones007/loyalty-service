// Обработчики GET запросов
package handlers

import "net/http"

func (h *Handler) BadHandler(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusBadRequest)
	rw.Write([]byte("Hey! Motherfucker! You know you don't belong here! Request incorrect!"))
}

// Get a list with information on user orders
func (h *Handler) GetInfoHandler(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("This is Get Info Handler!"))
}

// Get a current balance of points on user
func (h *Handler) GetCurrentBalanceHandler(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("This is Get Balance Handler!"))
}

// Get information for withdrawals user
func (h *Handler) GetWithdrawalsHandler(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("This is Get Withdraw Handler!"))
}
