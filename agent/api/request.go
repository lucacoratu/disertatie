package api

import (
	"bytes"
	"encoding/json"
	"errors"
	"net/http"
	"strconv"

	"github.com/lucacoratu/disertatie/agent/config"
	"github.com/lucacoratu/disertatie/agent/data"
	"github.com/lucacoratu/disertatie/agent/logging"
)

// Holds data necessary in order to communicate with collector
type APIHandler struct {
	logger        logging.ILogger
	configuration config.Configuration
}

// Creates a new instance of the APIHandler struct
func NewAPIHandler(logger logging.ILogger, configuration config.Configuration) *APIHandler {
	return &APIHandler{logger: logger, configuration: configuration}
}

// Registers the agent to the api
func (api *APIHandler) RegisterAgent(apiBaseUrl string, agentInfo data.AgentInformation) (string, error) {
	//Parse the data into a JSON
	bodyData, err := json.Marshal(agentInfo)
	//Check if an error occured when transforming machine info into JSON
	if err != nil {
		return "", errors.New("could not transform the machine info into JSON")
	}
	//Send the data to the api
	resp, err := http.Post(apiBaseUrl+"/registeragent", "application/json", bytes.NewBuffer(bodyData))
	//Check if an error occured when sending the request to the collector
	if err != nil {
		return "", errors.New("could not register the agent to api, " + err.Error())
	}
	//Check if the response says the operation was successful and get the uuid from the body
	responseData := data.RegisterProxyResponse{}
	err = responseData.FromJSON(resp.Body)
	if err != nil {
		return "", errors.New("could not parse register response, " + err.Error())
	}
	//Return the uuid of the proxy send by the server
	return responseData.Uuid, nil
}

// Sends logs to the API
func (api *APIHandler) SendLog(apiBaseUrl string, logData data.LogData) (bool, error) {
	api.logger.Info("Sending log to the API")
	//Parse the data into a JSON
	bodyData, err := json.Marshal(logData)
	//Check if an error occured when transforming machine info into JSON
	if err != nil {
		return false, errors.New("could not transform the log data into JSON")
	}
	//Send the data to the api
	resp, err := http.Post(apiBaseUrl+"/addlog", "application/json", bytes.NewBuffer(bodyData))
	//Check if an error occured when sending the request to the collector
	if err != nil {
		return false, errors.New("could not send the logs to api, " + err.Error())
	}
	//Check the status code of the response
	if resp.StatusCode != 200 {
		apiErr := data.APIError{}
		//Parse the error response from the API
		err := apiErr.FromJSON(resp.Body)
		//Check if an error occured when parsing the api error response
		if err != nil {
			return false, errors.New("could not parse error message from API, " + err.Error())
		}
		//Log the error message
		return false, errors.New("error on the server, code " + strconv.Itoa(int(apiErr.Code)) + ", message: " + apiErr.Message)
	}
	//The log has been added in the database
	return true, nil
}
