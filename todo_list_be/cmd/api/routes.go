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
	router.HandlerFunc(http.MethodGet, "/api/account", app.requireRole("ADMIN", app.test))

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
