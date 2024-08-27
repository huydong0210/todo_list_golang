package main

import (
	"net/http"
	"todo_list_be/internal/model"
)

func (app *application) createTodoItem(w http.ResponseWriter, r *http.Request) {
	type input struct {
		Name  string `json:"name"`
		State string `json:"state"`
	}
	var in input
	err := app.readJSON(w, r, &in)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	user, err := model.FindUserByUsername(app.db, app.contextGetUser(r).Username)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	item := model.TodoItem{
		Name:   in.Name,
		State:  in.State,
		UserId: user.ID,
	}
	err = model.CreateTodoItem(app.db, &item)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}
	err = app.writeJSON(w, http.StatusCreated, envelope{
		"message": "create todo item successfully",
	}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
func (app *application) deleteTodoItem(w http.ResponseWriter, r *http.Request) {
	id, err := app.readIDParam(r)
	user, err := model.FindUserByUsername(app.db, app.contextGetUser(r).Username)

	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	item, err := model.FindTodoItemById(app.db, int(id))
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	if item.UserId != user.ID {
		app.errorResponse(w, r, http.StatusUnauthorized, "You are not authorized to delete this item")
		return
	}
	err = model.DeleteTodoItem(app.db, int(id))
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	err = app.writeJSON(w, http.StatusOK, envelope{
		"message": "delete todo item successfully",
	}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
func (app *application) updateTodoItem(w http.ResponseWriter, r *http.Request) {

	type input struct {
		Name  string `json:"name"`
		State string `json:"state"`
	}
	var in input
	err := app.readJSON(w, r, &in)
	if err != nil {
		app.badRequestResponse(w, r, err)
		return
	}

	id, err := app.readIDParam(r)
	user, err := model.FindUserByUsername(app.db, app.contextGetUser(r).Username)

	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	item, err := model.FindTodoItemById(app.db, int(id))
	if err != nil {
		app.notFoundResponse(w, r)
		return
	}
	if item.UserId != user.ID {
		app.errorResponse(w, r, http.StatusUnauthorized, "You are not authorized to delete this item")
		return
	}
	item.Name = in.Name
	item.State = in.State
	err = model.UpdateTodoItem(app.db, int(id), &item)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
	err = app.writeJSON(w, http.StatusOK, envelope{
		"message": "update todo item successfully",
	}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}
}
