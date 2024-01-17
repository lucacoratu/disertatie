package handlers

import (
	"net/http"

	"github.com/lucacoratu/disertatie/api/config"
	"github.com/lucacoratu/disertatie/api/logging"
)

// Structure that holds data used by the HealthCheck routes
type HealthCheckHandler struct {
	logger        logging.ILogger
	configuration config.Configuration
}

// Creates a new instance of the HealthCheckHandler
func NewHealthCheckHandler(logger logging.ILogger, configuration config.Configuration) *HealthCheckHandler {
	return &HealthCheckHandler{logger: logger, configuration: configuration}
}

// Handler function for the endpoint that the proxies can access to test the connection to the collector
func (healthCheck *HealthCheckHandler) HealthCheck(rw http.ResponseWriter, r *http.Request) {
	rw.WriteHeader(http.StatusOK)
	response := "{\"status\":\"alive\"}"
	rw.Write([]byte(response))
}
