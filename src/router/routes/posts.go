package routes

import (
	"api/src/controllers"
	"net/http"
)

var postRoutes = []Route{
	{
		Uri:       "/posts",
		Method:    http.MethodPost,
		Function:  controllers.PostCreate,
		Protected: false,
	},
	{
		Uri:       "/posts",
		Method:    http.MethodGet,
		Function:  controllers.PostGetAllByUserId,
		Protected: false,
	},
	{
		Uri:       "/posts/{id}",
		Method:    http.MethodGet,
		Function:  controllers.PostGetOne,
		Protected: false,
	},
	{
		Uri:       "/posts/{id}",
		Method:    http.MethodPut,
		Function:  controllers.PostUpdate,
		Protected: false,
	},
	{
		Uri:       "/posts/{id}",
		Method:    http.MethodDelete,
		Function:  controllers.PostDelete,
		Protected: false,
	},
}
