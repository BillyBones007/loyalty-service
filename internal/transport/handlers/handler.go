package handlers

import (
	"github.com/BillyBones007/loyalty-service/internal/db"
)

// Handler type
type Handler struct {
	Storage db.Store
	Key     string
	Accrual string
}

// TODO: Здесь разместить функцию инициализации
func InitHandler(storage db.Store, skey string, addrAccrual string) *Handler {
	return &Handler{Storage: storage, Key: skey, Accrual: addrAccrual}
}
