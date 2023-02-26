package customerr

import (
	"errors"

	"github.com/jackc/pgx/v5"
)

var (
	ErrFieldEmpty            error = errors.New("login or password cannot be empty")
	ErrLoginIsExists         error = errors.New("login is exists")
	ErrLoginOrPassIncorrect  error = errors.New("login or password incorrect")
	ErrNoRows                error = pgx.ErrNoRows
	ErrNoCookie              error = errors.New("token not found")
	ErrSigningMethod         error = errors.New("unexpected signing method")
	ErrBadToken              error = errors.New("bad token")
	ErrTokenExp              error = errors.New("token is either expired or not active yet")
	ErrInvalidNumOrder       error = errors.New("invalid number order")
	ErrOrderIsExists         error = errors.New("the order has already been uploaded by another user")
	ErrTooManyRequests       error = errors.New("too many requests")
	ErrInternalBlackboxError error = errors.New("internal blackbox service error")
	ErrInsufficientFounds    error = errors.New("insufficient founds")
)
