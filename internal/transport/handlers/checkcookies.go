package handlers

import (
	"context"
	"net/http"

	"github.com/BillyBones007/loyalty-service/internal/customerr"
	"github.com/BillyBones007/loyalty-service/internal/tools/jwttoken"
)

// Type for pass to context
type AuthKey string

const Tkn AuthKey = "token"

// Check cookie
func (h *Handler) CheckCookies(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		cookie, err := r.Cookie("token")
		if err != nil {
			currToken := jwttoken.CurrentToken{Err: customerr.ErrNoCookie}
			r = r.WithContext(context.WithValue(r.Context(), Tkn, &currToken))
			next.ServeHTTP(w, r)
			return
		}
		currToken := jwttoken.ParseToken([]byte(h.Key), cookie.Value)
		r = r.WithContext(context.WithValue(r.Context(), Tkn, currToken))
		next.ServeHTTP(w, r)
	})
}
