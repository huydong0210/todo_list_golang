package main

import (
	"context"
	"net/http"
)

type contextKey struct {
	name string
}

var userContextKey = &contextKey{"UserPrincipal"}

type UserPrincipal struct {
	Username string `json:"username"`
	Role     string `json:"role"`
	Email    string `json:"email"`
}

var AnonymousUser = &UserPrincipal{}

func (app *application) contextSetUser(r *http.Request, user *UserPrincipal) *http.Request {
	ctx := context.WithValue(r.Context(), userContextKey, user)
	return r.WithContext(ctx)
}

func (app *application) contextGetUser(r *http.Request) *UserPrincipal {
	user, ok := r.Context().Value(userContextKey).(*UserPrincipal)
	if !ok {
		panic("missing user value in request context")
	}
	return user
}
func (user *UserPrincipal) isAnonymousUser() bool {
	return user == AnonymousUser
}
