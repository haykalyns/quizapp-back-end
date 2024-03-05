package config

import (
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

var JWT_KEY = []byte("afdafadjfaho1o34141adfafa1351")

type JWTClaim struct {
	Nama string
	Role string
	jwt.RegisteredClaims
}

var (
	AdminJWTKey = []byte("admin_secret_key")
	UserJWTKey  = []byte("user_secret_key")
)

func (c *JWTClaim) Valid() error {
	if time.Now().Unix() > c.ExpiresAt.Time.Unix() {
		return errors.New("token is expired")
	}
	return nil
}
