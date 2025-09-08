package token

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/google/uuid"
)

var (
	ErrInvalidToken = errors.New("token is invalid")
	ErrExpiredToken = errors.New("token has expired")
)

type payload struct {
	Username string `json:"username"`
	jwt.RegisteredClaims
}

func NewPayload(username string, duration time.Duration) (*payload, error) {
	TokenId, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}

	Payload := &payload{
		Username: username,
		RegisteredClaims: jwt.RegisteredClaims{
			ID:        TokenId.String(),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(duration)),
		},
	}

	return Payload, nil
}
func (Payload *payload) valid() error {
	if time.Now().After(Payload.RegisteredClaims.ExpiresAt.Time) {
		return  ErrExpiredToken
	}
	return nil
}
