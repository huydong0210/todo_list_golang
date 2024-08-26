package main

import (
	"net/http"
	"todo_list_be/internal/model"
)

func (app *application) loginHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	user, err := model.FindUserByUsername(app.db, input.Username)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	if !checkPasswordHash(input.Password, user.Password) {
		app.errorResponse(w, r, http.StatusUnauthorized, "incorrect password")
		return
	}
	token, err := app.GenerateToken(&user)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{
		"access-token": token,
	}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
}
func (app *application) signUpHandler(w http.ResponseWriter, r *http.Request) {
	var input struct {
		Username string `json:"username"`
		Password string `json:"password"`
		Email    string `json:"email"`
	}
	err := app.readJSON(w, r, &input)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	hashPass, err := HashPassword(input.Password)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	user := model.User{
		Username: input.Username,
		Password: hashPass,
		Email:    input.Email,
	}
	err = model.CreateUser(app.db, &user)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = model.InsertUserRoles(app.db, user.ID, "USER")
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{
		"message": "sign up successfully",
	}, nil)
}
