package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lucacoratu/disertatie/api/config"
	"github.com/lucacoratu/disertatie/api/data"
	request "github.com/lucacoratu/disertatie/api/data/request"
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

// Handler to register a new machine
func (mh *MachinesHandler) RegisterMachine(rw http.ResponseWriter, r *http.Request) {
	//Get the data from the request
	bodyData := request.RegisterMachineRequest{}
	err := bodyData.FromJSON(r.Body)
	//Check if an error occured when parsing the request body in the struct
	if err != nil {
		mh.logger.Error("Erorr occured when parsing the body for registering a machine", err.Error())
		//The request is bad so send an error message
		rw.WriteHeader(http.StatusBadRequest)
		retErr := data.APIError{Code: data.PARSE_ERROR, Message: err.Error()}
		retErr.ToJSON(rw)
		return
	}
	//Check the connection to the SSH
	//Save the machine in the database
	ipAddresses := make([]string, 0)
	ipAddresses = append(ipAddresses, bodyData.IPAddress)
	uuid, err := mh.dbConnection.InsertMachine(bodyData.OS, bodyData.Hostname, ipAddresses)
	//Check if an error occured when saving the machine in the database
	if err != nil {
		mh.logger.Error("Error occured when saving a new machine in the database", err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		retErr := data.APIError{Code: data.DATABASE_ERROR, Message: err.Error()}
		retErr.ToJSON(rw)
		return
	}
	//Return the UUID of the machine
	retData := response.MachinesRegisterResponse{Uuid: uuid}
	rw.WriteHeader(http.StatusOK)
	retData.ToJSON(rw)
}

// Handler to delete a machine
func (mh *MachinesHandler) DeleteMachine(rw http.ResponseWriter, r *http.Request) {
	//Get the machine uuid from mux vars
	vars := mux.Vars(r)
	machine_uuid := vars["machineuuid"]
	//Check if the machine uuid exists
	if machine_uuid == "" {
		mh.logger.Error("Error occured when deleting a machine, missing uuid")
		rw.WriteHeader(http.StatusBadRequest)
		retErr := data.APIError{Code: data.REQUEST_ERROR, Message: "machine uuid is mandatory for this operation"}
		retErr.ToJSON(rw)
		return
	}

	err := mh.dbConnection.DeleteMachine(machine_uuid)
	if err != nil {
		mh.logger.Error("Error occured when deleting machine from the database", err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		retErr := data.APIError{Code: data.DATABASE_ERROR, Message: err.Error()}
		retErr.ToJSON(rw)
		return
	}

	//Delete the machine from the database
	rw.WriteHeader(http.StatusOK)
	rw.Write([]byte("machine deleted"))
}
