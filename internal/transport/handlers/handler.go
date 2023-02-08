package handlers

import "github.com/BillyBones007/loyalty-service/internal/db"

// Handler type
type Handler struct {
	Storage db.Store
}
