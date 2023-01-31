// Обработчики POST запросов
package handlers

import "net/http"

// Registration new user
func RegisterUserHandler(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("This is Post Register!"))
}

// Authentication user
func AuthUserHandler(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("This is Post Auth User!"))
}

// Uploading the numder order for calculation
func UploadNumberOrderHandler(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("This is Post Upload Number!"))
}

// Withdrow points
func WithdrawPointsHandler(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("This is Post Withdraw Points!"))
}
