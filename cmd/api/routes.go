package main

import (
	"github.com/julienschmidt/httprouter"
	"net/http"
)

func (app *application) routes() *httprouter.Router {
	router := httprouter.New()

	router.HandlerFunc(http.MethodGet, "/v1/healthcheck", app.healthcheckHandler)
	//router.HandlerFunc(http.MethodPost, "/auth/login", app.loginHandler)
	router.HandlerFunc(http.MethodGet, "/auth/:id", app.showUserHandler) // im gonna delete it later it is used just for simple tests


	router.HandlerFunc(http.MethodPost, "/auth/login", app.loginHandler)
	router.HandlerFunc(http.MethodPost, "/users", app.registerUserHandler)
	router.HandlerFunc(http.MethodPost, "/torcedores", app.registerTorcedorHandler)


	protectedCampeonatos := app.requireAuthenticatedUser(http.HandlerFunc(app.listCampeonatosHandler))



	router.Handler(http.MethodGet, "/v1/campeonatos", protectedCampeonatos)
	router.Handler(http.MethodPost, "/broadcast", app.requireAuthenticatedUser(app.requireAdminUser(app.broadcastHandler)))
	router.HandlerFunc(http.MethodGet, "/v1/campeonatos/:id/partidas", app.listPartidasHandler)



	
	return router
}
