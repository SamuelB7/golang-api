package controllers

import (
	"api/src/database"
	"api/src/responses"
	"context"
	"net/http"
	"time"
)

// HealthCheck godoc
// @Summary Health Check
// @Description Get the health status of the application
// @Tags healthcheck
// @Accept json
// @Produce json
// @Success 200
// @Failure 503
// @Router /health [get]
func HealthCheck(w http.ResponseWriter, r *http.Request) {
	status := map[string]interface{}{
		"status":    "UP",
		"timestamp": time.Now().Format(time.RFC3339),
	}

	db, err := database.Connect()
	if err != nil {
		status["status"] = "DOWN"
		status["database"] = map[string]string{
			"status": "DOWN",
			"error":  "Failed to connect to database",
		}
		responses.JsonResponse(w, http.StatusServiceUnavailable, status)
		return
	}
	defer db.Close(context.Background())

	ctx, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	if err := db.Ping(ctx); err != nil {
		status["status"] = "DOWN"
		status["database"] = map[string]string{
			"status": "DOWN",
			"error":  "Database ping failed",
		}
		responses.JsonResponse(w, http.StatusServiceUnavailable, status)
		return
	}

	var dbVersion string
	err = db.QueryRow(ctx, "SELECT version()").Scan(&dbVersion)

	dbStatus := map[string]string{
		"status": "UP",
	}

	if err != nil {
		dbStatus["version_error"] = "Could not retrieve database version"
	} else {
		dbStatus["version"] = dbVersion
	}

	status["database"] = dbStatus

	responses.JsonResponse(w, http.StatusOK, status)
}
