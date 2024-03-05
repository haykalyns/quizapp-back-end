package middleware

import (
	"net/http"
	"quiz-project/src/config"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt/v4"
)

func AuthenticateAdmin() gin.HandlerFunc {
	return func(c *gin.Context) {

		tokenString, err := c.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				response := map[string]string{"message": "unauthorized"}
				c.JSON(http.StatusUnauthorized, response)
				c.Abort()
				return
			}
		}

		claim := &config.JWTClaim{}
		token, err := jwt.ParseWithClaims(tokenString, claim, func(t *jwt.Token) (interface{}, error) {
			// Menentukan kunci berdasarkan peran pengguna (admin atau user)
			var key []byte
			if claim.Role == "admin" {
				key = config.AdminJWTKey
			} else {
				key = config.UserJWTKey
			}
			return key, nil
		})

		if err != nil {
			v, _ := err.(*jwt.ValidationError)
			switch v.Errors {
			case jwt.ValidationErrorSignatureInvalid:
				response := map[string]string{"message": "unauthorized"}
				c.JSON(http.StatusUnauthorized, response)
				c.Abort()
				return
			case jwt.ValidationErrorExpired:
				response := map[string]string{"message": "unauthorized, token expired!"}
				c.JSON(http.StatusUnauthorized, response)
				c.Abort()
				return
			default:
				response := map[string]string{"message": "unauthorized"}
				c.JSON(http.StatusUnauthorized, response)
				c.Abort()
				return
			}
		}

		if !token.Valid {
			response := map[string]string{"message": "unauthorized"}
			c.JSON(http.StatusUnauthorized, response)
			c.Abort()
			return
		}

		// Check if the role is admin
		if claim.Role != "admin" {
			response := map[string]string{"message": "unauthorized, user is not admin"}
			c.JSON(http.StatusUnauthorized, response)
			c.Abort()
			return
		}

		c.Next()
	}
}

func AuthenticateUser() gin.HandlerFunc {
	return func(c *gin.Context) {

		tokenString, err := c.Cookie("token")
		if err != nil {
			if err == http.ErrNoCookie {
				response := map[string]string{"message": "unauthorized"}
				c.JSON(http.StatusUnauthorized, response)
				c.Abort()
				return
			}
		}

		claim := &config.JWTClaim{}
		token, err := jwt.ParseWithClaims(tokenString, claim, func(t *jwt.Token) (interface{}, error) {
			// Menentukan kunci berdasarkan peran pengguna (admin atau user)
			var key []byte
			if claim.Role == "admin" {
				key = config.AdminJWTKey
			} else {
				key = config.UserJWTKey
			}
			return key, nil
		})

		if err != nil {
			v, _ := err.(*jwt.ValidationError)
			switch v.Errors {
			case jwt.ValidationErrorSignatureInvalid:
				response := map[string]string{"message": "unauthorized"}
				c.JSON(http.StatusUnauthorized, response)
				c.Abort()
				return
			case jwt.ValidationErrorExpired:
				response := map[string]string{"message": "unauthorized, token expired!"}
				c.JSON(http.StatusUnauthorized, response)
				c.Abort()
				return
			default:
				response := map[string]string{"message": "unauthorized"}
				c.JSON(http.StatusUnauthorized, response)
				c.Abort()
				return
			}
		}

		if !token.Valid {
			response := map[string]string{"message": "unauthorized"}
			c.JSON(http.StatusUnauthorized, response)
			c.Abort()
			return
		}

		if claim.Role != "user" {
			response := map[string]string{"message": "unauthorized, user is not user"}
			c.JSON(http.StatusUnauthorized, response)
			c.Abort()
			return
		}

		c.Next()
	}
}
