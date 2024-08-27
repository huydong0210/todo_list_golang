package main

import (
	"net/http"
	"todo_list_be/internal/model"
)

func (app *application) findAllUsers(w http.ResponseWriter, r *http.Request) {
	users, err := model.FindAllUsers(app.db)
	if err != nil {
		app.serverErrorResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusOK, envelope{
		"users": users,
	}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
