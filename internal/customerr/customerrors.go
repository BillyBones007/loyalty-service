package customerr

import (
	"errors"

	"github.com/jackc/pgx/v5"
)

var (
	ErrFieldEmpty           error = errors.New("login or password cannot be empty")
	ErrLoginIsExists        error = errors.New("login is exists")
	ErrLoginOrPassIncorrect error = errors.New("login or password incorrect")
	ErrNoRows               error = pgx.ErrNoRows
)
