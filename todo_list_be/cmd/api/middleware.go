package main

import (
	"fmt"
	"github.com/golang-jwt/jwt/v5"
	_ "github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	_ "golang.org/x/crypto/bcrypt"
	"net/http"
	"strings"
	"time"
	"todo_list_be/internal/model"
)

type CustomClaims struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	Email    string `json:"email"`
	jwt.RegisteredClaims
}

func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

func checkPasswordHash(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}
func (app *application) GenerateToken(user *model.User) (string, error) {
	roles, err := model.FindRolesByUserId(app.db, user.ID)
	var roleNames string
	for _, role := range roles {
		roleNames += role.Name + " "
	}
	roleNames = strings.TrimSpace(roleNames)
	claims := CustomClaims{
		Username: user.Username,
		Role:     roleNames,
		Email:    user.Email,
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(time.Hour * 24)),
			IssuedAt:  jwt.NewNumericDate(time.Now()),
			Issuer:    "your-app-name",
		},
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)

	tokenString, err := token.SignedString([]byte(app.config.security.jwtSecret))
	return tokenString, err
}
func (app *application) ParseToken(tokenString string) (*jwt.Token, error) {
	secretKey := []byte(app.config.security.jwtSecret)
	var result CustomClaims
	token, err := jwt.ParseWithClaims(tokenString, &result, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}
		return secretKey, nil
	})
	return token, err
}
func (app *application) requireRole(requiredRole string, next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user := app.contextGetUser(r)
		if user.IsAnonymous() {
			app.authenticationRequiredResponse(w, r)
			return
		}
		roles, err := model.FindRolesByUserId(app.db, user.ID)
		if err != nil {
			app.serverErrorResponse(w, r, err)
		}
		hashRole := false
		for _, role := range roles {
			if role.Name == requiredRole {
				hashRole = true
				break
			}
		}
		if !hashRole {
			app.authenticationRequiredResponse(w, r)
		}
		next.ServeHTTP(w, r)
	}
}
