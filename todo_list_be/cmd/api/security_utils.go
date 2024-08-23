package main

import (
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	_ "golang.org/x/crypto/bcrypt"
	"time"
)

type UserPrincipal struct {
	id       int
	username string
	role     string
}

type SecurityContext struct {
	user *UserPrincipal
}

func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
func (app *application) GenerateToken(user *UserPrincipal) (string, error) {
	claims := jwt.MapClaims{
		"username": user.username,
		"role":     user.role,
		"iat":      time.Now().Unix(),
		"exp":      time.Now().Add(app.config.security.exp).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString(app.config.security.jwtSecret)
	return tokenString, err
}
func (app *application) ParseToken(tokenString string) bool {

}
