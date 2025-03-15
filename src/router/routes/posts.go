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
		Protected: true,
	},
	{
		Uri:       "/posts-by-user",
		Method:    http.MethodGet,
		Function:  controllers.PostGetAllByUserId,
		Protected: true,
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
		Protected: true,
	},
	{
		Uri:       "/posts/{id}",
		Method:    http.MethodDelete,
		Function:  controllers.PostDelete,
		Protected: true,
	},
}
