package handlers

import (
	"net/http"

	"github.com/lucacoratu/disertatie/api/config"
	"github.com/lucacoratu/disertatie/api/data"
	response "github.com/lucacoratu/disertatie/api/data/response"
	"github.com/lucacoratu/disertatie/api/database"
	"github.com/lucacoratu/disertatie/api/logging"
)

type MachinesHandler struct {
	logger        logging.ILogger
	configuration config.Configuration
	dbConnection  database.IConnection
}

// Creates a new handler that will hold the functions necessary for registering proxies
func NewMachinesHandler(logger logging.ILogger, configuration config.Configuration, dbConnection database.IConnection) *MachinesHandler {
	return &MachinesHandler{logger: logger, configuration: configuration, dbConnection: dbConnection}
}

// Handler to get all the machines from the database
func (mh *MachinesHandler) GetMachines(rw http.ResponseWriter, r *http.Request) {
	//Get the machines from the database
	machines, err := mh.dbConnection.GetMachines()
	mh.logger.Debug(machines)
	//Check if an error occured when getting the machines from the database
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		retErr := data.APIError{Code: data.DATABASE_ERROR, Message: err.Error()}
		retErr.ToJSON(rw)
		return
	}
	//Return the list of machines to the client
	responseData := response.MachinesGetResponse{Machines: machines}
	rw.WriteHeader(http.StatusOK)
	responseData.ToJSON(rw)
}
