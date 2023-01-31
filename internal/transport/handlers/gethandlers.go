// Обработчики GET запросов
package handlers

import "net/http"

func CoreHandler(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("Hallo my little hacker!!!"))
}

// Get a list with information on user orders
func GetInfoHandler(rw http.ResponseWriter, r *http.Request) {
}

// Get a current ballance of points on user
func GetCurrentBallanceHandler(rw http.ResponseWriter, r *http.Request) {
}

// Get information for withdrawals user
func GetWithdrawalsHandler(rw http.ResponseWriter, r *http.Request) {
}
