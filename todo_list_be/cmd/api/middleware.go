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
		userPrincipal := app.contextGetUser(r)
		if userPrincipal.isAnonymousUser() {
			app.authenticationRequiredResponse(w, r)
			return
		}
		roles := strings.Split(userPrincipal.Role, " ")
		hashRole := false
		for _, role := range roles {
			if role == requiredRole {
				hashRole = true
				break
			}
		}
		if !hashRole {
			app.authenticationRequiredResponse(w, r)
			return
		}
		next.ServeHTTP(w, r)
	}
}

func (app *application) authenticate(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Add("Vary", "Authorization")

		authorizationHeader := r.Header.Get("Authorization")

		if authorizationHeader == "" {
			r = app.contextSetUser(r, AnonymousUser)
			next.ServeHTTP(w, r)
			return
		}

		headerParts := strings.Split(authorizationHeader, " ")
		if len(headerParts) != 2 || headerParts[0] != "Bearer" {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}

		token := headerParts[1]
		parsedToken, err := app.ParseToken(token)
		if err != nil {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}
		claims, oke := parsedToken.Claims.(*CustomClaims)
		if !oke || !parsedToken.Valid {
			app.invalidAuthenticationTokenResponse(w, r)
			return
		}
		userPrincipal := &UserPrincipal{
			Username: claims.Username,
			Role:     claims.Role,
			Email:    claims.Email,
		}
		r = app.contextSetUser(r, userPrincipal)
		next.ServeHTTP(w, r)

	})
}
