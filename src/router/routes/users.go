package routes

import (
	"api/src/controllers"
	"net/http"
)

var userRoutes = []Route{
	{
		Uri:       "/users",
		Method:    http.MethodPost,
		Function:  controllers.UserCreate,
		Protected: false,
	},
	{
		Uri:       "/users",
		Method:    http.MethodGet,
		Function:  controllers.UserGetAll,
		Protected: false,
	},
	{
		Uri:       "/users/{id}",
		Method:    http.MethodGet,
		Function:  controllers.UserGetOne,
		Protected: false,
	},
	{
		Uri:       "/users/{id}",
		Method:    http.MethodPut,
		Function:  controllers.UserUpdate,
		Protected: false,
	},
	{
		Uri:       "/users/{id}",
		Method:    http.MethodDelete,
		Function:  controllers.UserDelete,
		Protected: false,
	},
}
