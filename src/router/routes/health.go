package routes

import (
	"api/src/controllers"
	"net/http"
)

var healthRoutes = []Route{
	{
		Uri:       "/health",
		Method:    http.MethodGet,
		Function:  controllers.HealthCheck,
		Protected: false,
	},
}
