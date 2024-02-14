package handlers

import (
	"net/http"

	"github.com/go-playground/validator/v10"
	"github.com/gorilla/mux"
	"github.com/lucacoratu/disertatie/api/config"
	"github.com/lucacoratu/disertatie/api/data"
	request "github.com/lucacoratu/disertatie/api/data/request"
	response "github.com/lucacoratu/disertatie/api/data/response"
	"github.com/lucacoratu/disertatie/api/database"
	"github.com/lucacoratu/disertatie/api/logging"
)

type AgentsHandler struct {
	logger        logging.ILogger
	configuration config.Configuration
	dbConnection  database.IConnection
}

// Creates a new handler that will hold the functions necessary for registering proxies
func NewAgentsHandler(logger logging.ILogger, configuration config.Configuration, dbConnection database.IConnection) *AgentsHandler {
	return &AgentsHandler{logger: logger, configuration: configuration, dbConnection: dbConnection}
}

// Handler for registerning a new agent
func (ah *AgentsHandler) RegisterAgent(rw http.ResponseWriter, r *http.Request) {
	//Parse the JSON data from the body into the correct structure
	requestData := request.AgentInformation{}
	err := requestData.FromJSON(r.Body)
	if err != nil {
		ah.logger.Error("Could not parse JSON data from body in the coresponding structure", err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		retErr := data.APIError{Code: data.PARSE_ERROR, Message: err.Error()}
		retErr.ToJSON(rw)
		return
	}

	//Validate the data
	//Initialize the validator of the json data
	validate := validator.New(validator.WithRequiredStructEnabled())
	//Validate the fields of the struct
	err = validate.Struct(requestData)
	//Check if an error occured when validating the data
	if err != nil {
		ah.logger.Error("Could not add agent, data validation error", err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		retErr := data.APIError{Code: data.REQUEST_ERROR, Message: err.Error()}
		retErr.ToJSON(rw)
		return
	}

	ah.logger.Debug(requestData)

	//Check if the machine information received from the agent doesn't corespond to an existing machine
	machine_id, err := ah.dbConnection.CheckMachineExists(requestData.MachineInfo.OS, requestData.MachineInfo.Hostname)
	if err != nil {
		ah.logger.Error("Could not check if the machine exists", err.Error())
		//Create the custom error message and return to client
		rw.WriteHeader(http.StatusInternalServerError)
		retErr := data.APIError{Code: data.DATABASE_ERROR, Message: err.Error()}
		retErr.ToJSON(rw)
		return
	}

	//If the machine does not exist in the database add it
	if machine_id == "" {
		machine_id, err = ah.dbConnection.InsertMachine(requestData.MachineInfo.OS, requestData.MachineInfo.Hostname, requestData.MachineInfo.IPAddresses)
		if err != nil {
			ah.logger.Error("Could not insert the machine in the database", err.Error())
			//Create the custom error message and return to client
			rw.WriteHeader(http.StatusInternalServerError)
			retErr := data.APIError{Code: data.DATABASE_ERROR, Message: err.Error()}
			retErr.ToJSON(rw)
			return
		}
	}

	//Insert the data into the database
	uuid, err := ah.dbConnection.InsertAgent(requestData.Protocol, requestData.IPAddress, requestData.Port, requestData.WebServerProtocol, requestData.WebServerIP, requestData.WebServerPort, machine_id)
	if err != nil {
		ah.logger.Error("Could not insert agent instance in the database", err.Error())
		//Create the custom error message and return to client
		rw.WriteHeader(http.StatusInternalServerError)
		retErr := data.APIError{Code: data.DATABASE_ERROR, Message: err.Error()}
		retErr.ToJSON(rw)
		return
	}
	//Send the response back to the client
	responseData := response.RegisterProxyResponse{Uuid: uuid}
	rw.WriteHeader(http.StatusOK)
	responseData.ToJSON(rw)
}

// Handler for receiving logs from the agent
func (ah *AgentsHandler) AddLog(rw http.ResponseWriter, r *http.Request) {
	//Get data from the request body
	logData := data.LogData{}
	err := logData.FromJSON(r.Body)
	//Check if an error occured when parsing the JSON body
	if err != nil {
		//Create the custom error message and return to client
		rw.WriteHeader(http.StatusInternalServerError)
		retErr := data.APIError{Code: data.PARSE_ERROR, Message: err.Error()}
		retErr.ToJSON(rw)
		return
	}
	//Insert the log in the database
	_, err = ah.dbConnection.InsertLog(logData)
	//Check if an error occured when inserting the log in the database
	if err != nil {
		//Create the custom error message and return to client
		rw.WriteHeader(http.StatusInternalServerError)
		retErr := data.APIError{Code: data.DATABASE_ERROR, Message: err.Error()}
		retErr.ToJSON(rw)
		return
	}

	//Send the success message
	rw.WriteHeader(http.StatusOK)
	//Create the success message
	message := data.SuccessMessage{Message: "log has been added"}
	message.ToJSON(rw)
}

// Function that handles GET request on /api/v1/agents (returns all the agents)
func (ah *AgentsHandler) GetAgents(rw http.ResponseWriter, r *http.Request) {
	//Get the agents from the database
	agents, err := ah.dbConnection.GetAgents()
	//Check if an error occured when getting the agents from the database
	if err != nil {
		//Create the error structure that will describe the error
		rw.WriteHeader(http.StatusInternalServerError)
		retErr := data.APIError{Code: data.DATABASE_ERROR, Message: err.Error()}
		retErr.ToJSON(rw)
		return
	}
	//Complete the agents data with the machine information
	for index := range agents {
		//Get the machine information
		machine, err := ah.dbConnection.GetMachine(agents[index].MachineId)
		if err != nil {
			ah.logger.Warning("could not get information about machine", agents[index].MachineId, "where agent", agents[index].ID, "is deployed")
			continue
		}
		agents[index].MachineOS = machine.OS
		agents[index].MachineHostname = machine.Hostname
		agents[index].MachineIPAddreses = machine.IPAddresses
	}

	agResponse := response.AgentsGetResponse{Agents: agents}
	ah.logger.Debug(agResponse)
	//Return the agents
	rw.WriteHeader(http.StatusOK)
	agResponse.ToJSON(rw)
}

// Function for getting a single agent
func (ah *AgentsHandler) GetAgent(rw http.ResponseWriter, r *http.Request) {
	//Get the agent uuid from mux vars
	vars := mux.Vars(r)
	agent_uuid := vars["uuid"]
	//Check if the agent UUID has been specified
	if agent_uuid == "" {
		//Create the custom error message and return to client
		ah.logger.Error("Error occured when getting an agent details, missing UUID")
		rw.WriteHeader(http.StatusInternalServerError)
		retErr := data.APIError{Code: data.REQUEST_ERROR, Message: "missing agent uuid"}
		retErr.ToJSON(rw)
		return
	}
	//Get the agent details from the database
	agent, err := ah.dbConnection.GetAgent(agent_uuid)
	//Check if an error occured when getting the agent from the database
	if err != nil {
		ah.logger.Error("Error occured when retrieving the agent from the database", err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		retErr := data.APIError{Code: data.DATABASE_ERROR, Message: err.Error()}
		retErr.ToJSON(rw)
		return
	}
	//Return the agent back to the client
	rw.WriteHeader(http.StatusOK)
	agent.ToJSON(rw)
}

// Handler for modifying an agent
func (ah *AgentsHandler) ModifyAgent(rw http.ResponseWriter, r *http.Request) {
	//Get the agent uuid from mux vars
	vars := mux.Vars(r)
	agent_uuid := vars["uuid"]
	//Check if the agent UUID has been specified
	if agent_uuid == "" {
		//Create the custom error message and return to client
		ah.logger.Error("Error occured when modifying an agent details, missing UUID")
		rw.WriteHeader(http.StatusInternalServerError)
		retErr := data.APIError{Code: data.REQUEST_ERROR, Message: "missing agent uuid"}
		retErr.ToJSON(rw)
		return
	}
	//Get the data from the request body
	agent := data.UpdateAgent{}
	err := agent.FromJSON(r.Body)
	//Check if an error occured when parsing request data from body
	if err != nil {
		//Create the custom error message and return to client
		ah.logger.Error("Error occured when modifying an agent details, failed to parse request body from JSON", err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		retErr := data.APIError{Code: data.REQUEST_ERROR, Message: "failed to parse body from JSON"}
		retErr.ToJSON(rw)
		return
	}

	//Update the agent in the database
	err = ah.dbConnection.ModifyAgent(agent_uuid, agent)
	//Check if an erorr occured
	if err != nil {
		ah.logger.Error("Error occured when modifying an agent details, failed to update agent row in the database", err.Error())
		rw.WriteHeader(http.StatusInternalServerError)
		retErr := data.APIError{Code: data.DATABASE_ERROR, Message: "failed to update row in the database"}
		retErr.ToJSON(rw)
		return
	}

	rw.WriteHeader(http.StatusOK)
}
