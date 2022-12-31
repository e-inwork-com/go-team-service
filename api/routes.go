// Copyright 2022, e-inwork.com. All rights reserved.

package api

import (
	"expvar"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func (app *Application) routes() http.Handler {
	router := httprouter.New()

	router.NotFound = http.HandlerFunc(app.notFoundResponse)
	router.MethodNotAllowed = http.HandlerFunc(app.methodNotAllowedResponse)

	router.HandlerFunc(http.MethodGet, "/service/teams/health", app.healthcheckHandler)
	router.HandlerFunc(http.MethodPost, "/service/teams", app.requireAuthenticated(app.createTeamHandler))
	router.HandlerFunc(http.MethodGet, "/service/teams/me", app.requireAuthenticated(app.getTeamHandler))
	router.HandlerFunc(http.MethodPatch, "/service/teams/:id", app.requireAuthenticated(app.patchTeamHandler))
	router.HandlerFunc(http.MethodGet, "/service/teams/pictures/:file", app.getProfilePictureHandler)

	router.Handler(http.MethodGet, "/service/profiles/debug/vars", expvar.Handler())

	return app.metrics(app.recoverPanic(app.enableCORS(app.rateLimit(app.authenticate(router)))))
}