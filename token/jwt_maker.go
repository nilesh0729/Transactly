package token

import (
	"errors"
	"fmt"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

const minSecretKeySize = 32

type JWTMAKER struct {
	secretkey string
}

func NewJWTMAKER(secretkey string) (Maker, error) {
	if len(secretkey) < minSecretKeySize {
		return nil, fmt.Errorf("invalid size of the secret key: must be atleast %d characters", minSecretKeySize)
	}
	return &JWTMAKER{secretkey}, nil
}

func (maker *JWTMAKER) CreateToken(username string, duration time.Duration) (string, error) {
	payload, err := NewPayload(username, duration)
	if err != nil {
		return "", err
	}
	jwtToken := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return jwtToken.SignedString([]byte(maker.secretkey))

}
func (maker *JWTMAKER) VerifyToken(token string) (*payload, error) {
	Keyfunc := func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return nil, ErrInvalidToken
		}
		return []byte(maker.secretkey), nil
	}
	jwtToken, err := jwt.ParseWithClaims(token, &payload{}, Keyfunc)
	if err != nil {
		if errors.Is(err, jwt.ErrTokenExpired){
			return nil, ErrExpiredToken
		}
		return nil, ErrInvalidToken
	}
	payload, ok := jwtToken.Claims.(*payload)
	if !ok {
		return nil, ErrInvalidToken
	}

	return payload, nil
}
