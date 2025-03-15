package routes

import (
	"api/src/middlewares"
	"net/http"

	"github.com/gorilla/mux"
)

type Method string

const (
	GET    Method = http.MethodGet
	POST   Method = http.MethodPost
	PUT    Method = http.MethodPut
	DELETE Method = http.MethodDelete
)

type Route struct {
	Uri       string
	Method    Method
	Function  func(http.ResponseWriter, *http.Request)
	Protected bool
}

func ConfigRoutes(r *mux.Router) *mux.Router {
	routes := userRoutes
	routes = append(routes, authRoutes...)
	routes = append(routes, postRoutes...)

	for _, route := range routes {
		if route.Protected {
			r.HandleFunc(route.Uri, middlewares.Auth(route.Function)).Methods(string(route.Method))
		} else {
			r.HandleFunc(route.Uri, route.Function).Methods(string(route.Method))
		}
	}

	return r
}
