package router

import (
	"api/src/router/routes"

	"github.com/gorilla/mux"
	httpSwagger "github.com/swaggo/http-swagger"
)

func GenerateRouter() *mux.Router {
	r := mux.NewRouter()
	r.PathPrefix("/swagger/").Handler(httpSwagger.Handler(
		httpSwagger.URL("/swagger/doc.json"),
		httpSwagger.DeepLinking(true),
		httpSwagger.DocExpansion("none"),
	))
	apiRouter := r.PathPrefix("/api").Subrouter()
	routes.ConfigRoutes(apiRouter)
	return r
}
