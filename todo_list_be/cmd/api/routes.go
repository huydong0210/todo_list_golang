package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() http.Handler {
	router := httprouter.New()
	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodPost, "/api/login", app.loginHandler)
	router.HandlerFunc(http.MethodPost, "/api/sign-up", app.signUpHandler)
	router.HandlerFunc(http.MethodGet, "/api/account", app.requireRole("USER", app.test))
	router.HandlerFunc(http.MethodGet, "/api/list-users", app.requireRole("ADMIN", app.findAllUsers))
	router.HandlerFunc(http.MethodPost, "/api/todo-item", app.requireRole("USER", app.createTodoItem))
	router.HandlerFunc(http.MethodDelete, "/api/todo-item/:id", app.requireRole("USER", app.deleteTodoItem))
	router.HandlerFunc(http.MethodPut, "/api/todo-item/:id", app.requireRole("USER", app.updateTodoItem))
	return app.authenticate(router)
}

func (app *application) test(w http.ResponseWriter, r *http.Request) {

	err := app.writeJSON(w, http.StatusOK, envelope{
		"data": "oke roi nay",
	}, nil)
	if err != nil {
		app.serverErrorResponse(w, r, err)
	}

}
