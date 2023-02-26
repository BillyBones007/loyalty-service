package jwttoken

import (
	"errors"
	"log"

	"github.com/BillyBones007/loyalty-service/internal/customerr"
	"github.com/golang-jwt/jwt/v4"
)

type CurrentToken struct {
	Err         error
	ClaimsToken jwt.MapClaims
}

// Return token. Accepts two parameters: secret key and user uuid
func GetTokenString(key []byte, uuid string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"user": uuid,
	})
	tokenString, err := token.SignedString(key)
	if err != nil {
		log.Printf("error from GetToken function: %s\n", err)
		return "", err
	}
	return tokenString, nil
}

// Parse token. Return pointer to CurrentToken.
func ParseToken(key []byte, tokenString string) *CurrentToken {
	currToken := CurrentToken{}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, customerr.ErrSigningMethod
		}
		return key, nil
	})
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		currToken.Err = nil
		currToken.ClaimsToken = claims
		return &currToken

	} else if errors.Is(err, jwt.ErrTokenMalformed) {
		currToken.Err = customerr.ErrBadToken

	} else if errors.Is(err, jwt.ErrTokenExpired) || errors.Is(err, jwt.ErrTokenNotValidYet) {
		currToken.Err = customerr.ErrTokenExp
	}
	currToken.Err = err
	return &currToken
}
