package routes

import (
	"api/src/controllers"
	"net/http"
)

var authRoutes = []Route{
	{
		Uri:       "/login",
		Method:    http.MethodPost,
		Function:  controllers.Login,
		Protected: false,
	},
	{
		Uri:       "/sign-in",
		Method:    http.MethodPost,
		Function:  controllers.SignIn,
		Protected: false,
	},
}
