package handlers

import (
	"fmt"
	"net/http"
	"sort"
	"strings"
	"time"

	"github.com/gorilla/mux"
	"github.com/lucacoratu/disertatie/api/config"
	"github.com/lucacoratu/disertatie/api/data"
	"github.com/lucacoratu/disertatie/api/database"
	"github.com/lucacoratu/disertatie/api/logging"

	response "github.com/lucacoratu/disertatie/api/data/response"
)

type LogsHandler struct {
	logger        logging.ILogger
	configuration config.Configuration
	dbConnection  database.IConnection
}

// Creates a new handler that will hold the functions necessary for registering proxies
func NewLogsHandler(logger logging.ILogger, configuration config.Configuration, dbConnection database.IConnection) *LogsHandler {
	return &LogsHandler{logger: logger, configuration: configuration, dbConnection: dbConnection}
}

// Handler to get all the logs from the database for a specific agent
func (lh *LogsHandler) GetLogsShort(rw http.ResponseWriter, r *http.Request) {
	//Get the agent uuid from the URL
	vars := mux.Vars(r)
	uuid := vars["uuid"]
	//Check if the uuid is available
	if uuid == "" {
		//Send an error message
		rw.WriteHeader(http.StatusBadRequest)
		apiErr := data.APIError{Code: data.REQUEST_ERROR, Message: "uuid missing"}
		apiErr.ToJSON(rw)
		return
	}
	//Get the logs of an agent
	logs, err := lh.dbConnection.GetAgentLogsShort(uuid)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		apiErr := data.APIError{Code: data.DATABASE_ERROR, Message: err.Error()}
		apiErr.ToJSON(rw)
		return
	}

	//Sort the logs by timestamp descending (the most recent ones should be first)
	sort.Slice(logs[:], func(i, j int) bool {
		return logs[i].Timestamp > logs[j].Timestamp
	})

	lh.logger.Debug(logs)

	//Send the logs back to the client
	respData := response.LogsGetResponse{Logs: logs}
	rw.WriteHeader(http.StatusOK)
	respData.ToJSON(rw)
}

// Handler for getting methods metrics
func (lh *LogsHandler) GetLogsMethodMetrics(rw http.ResponseWriter, r *http.Request) {
	//Get the agent uuid from the URL
	vars := mux.Vars(r)
	uuid := vars["uuid"]
	//Check if the uuid is available
	if uuid == "" {
		//Send an error message
		rw.WriteHeader(http.StatusBadRequest)
		apiErr := data.APIError{Code: data.REQUEST_ERROR, Message: "uuid missing"}
		apiErr.ToJSON(rw)
		return
	}
	//Get the logs of an agent
	logs, err := lh.dbConnection.GetAgentLogsShort(uuid)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		apiErr := data.APIError{Code: data.DATABASE_ERROR, Message: err.Error()}
		apiErr.ToJSON(rw)
		return
	}

	//Create a map which will hold all the request methods and the number of occurences
	var methodsOccurencesMap map[string]int64 = make(map[string]int64)

	//Parse the logs and find how many times each request method appears
	for _, log := range logs {
		method := strings.Split(log.RequestPreview, " ")[0]
		_, ok := methodsOccurencesMap[method]
		if !ok {
			methodsOccurencesMap[method] = 0
		}
		methodsOccurencesMap[method] += 1
	}

	//Log the occurences of each method for easy debugging
	//lh.logger.Debug(methodsOccurencesMap)

	responseData := make([]data.MethodsMetrics, 0)
	//Prepare the response for the client
	//Counter will be used as an ID
	var counter int64 = 0
	//Loop through all the keys and values in the occurences dictionary
	for key, value := range methodsOccurencesMap {
		//Create the structure which will be returned to the client
		responseData = append(responseData, data.MethodsMetrics{Id: counter, Method: key, Count: value})
		//Increment the counter
		counter += 1
	}
	//Sort the response data based on count descending
	sort.Slice(responseData[:], func(i, j int) bool {
		return responseData[i].Count > responseData[j].Count
	})

	//Create the response structure
	resp := response.MethodMetricsResponse{Metrics: responseData}

	//Send the response to the client
	rw.WriteHeader(http.StatusOK)
	resp.ToJSON(rw)
}

// Handler for getting number of requests each day
func (lh *LogsHandler) GetLogsCountPerDay(rw http.ResponseWriter, r *http.Request) {
	//Get the agent uuid from the URL
	vars := mux.Vars(r)
	uuid := vars["uuid"]
	//Check if the uuid is available
	if uuid == "" {
		//Send an error message
		rw.WriteHeader(http.StatusBadRequest)
		apiErr := data.APIError{Code: data.REQUEST_ERROR, Message: "uuid missing"}
		apiErr.ToJSON(rw)
		return
	}
	//Get the logs of an agent
	logs, err := lh.dbConnection.GetAgentLogsShort(uuid)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		apiErr := data.APIError{Code: data.DATABASE_ERROR, Message: err.Error()}
		apiErr.ToJSON(rw)
		return
	}

	//Sort the logs by timestamp ascending (the most recent ones should be last)
	sort.Slice(logs[:], func(i, j int) bool {
		return logs[i].Timestamp < logs[j].Timestamp
	})

	//Create a map which will have the occurences per day
	var mapOccurencesPerDay map[string]int64 = make(map[string]int64)

	//Count the number of logs each day
	for _, log := range logs {
		//Convert the unix timestamp to time structure
		logTime := time.Unix(log.Timestamp, 0)
		//Get the year, month and day of the request
		year, month, day := logTime.Date()
		//Create the map key
		mapKey := fmt.Sprintf("%d %s %d", day, month.String(), year)
		//Check if the map key exists
		_, ok := mapOccurencesPerDay[mapKey]
		//If the map key does not exist then initialize it with 0
		if !ok {
			mapOccurencesPerDay[mapKey] = 0
		}
		//Increment the number of requests that day
		mapOccurencesPerDay[mapKey] += 1
		//lh.logger.Debug(mapKey)
	}
	//lh.logger.Debug(mapOccurencesPerDay)

	//Create the data structure which will be returned to the client
	responseData := make([]data.DayMetrics, 0)
	//Create the counter which will act as the id of the data structure
	var counter int64 = 0
	//Loop through all the keys and values of the occurences map
	for key, value := range mapOccurencesPerDay {
		responseData = append(responseData, data.DayMetrics{Id: counter, Date: key, Count: value})
		//Increment the counter
		counter += 1
	}

	//lh.logger.Debug(responseData)

	//Create the response for the client
	resp := response.DaysMetricsResponse{Metrics: responseData}
	rw.WriteHeader(http.StatusOK)
	resp.ToJSON(rw)
}

// Handler for getting response status codes metrics
func (lh *LogsHandler) GetResponseStatusCodesMetrics(rw http.ResponseWriter, r *http.Request) {
	//Get the agent uuid from the URL
	vars := mux.Vars(r)
	uuid := vars["uuid"]
	//Check if the uuid is available
	if uuid == "" {
		//Send an error message
		rw.WriteHeader(http.StatusBadRequest)
		apiErr := data.APIError{Code: data.REQUEST_ERROR, Message: "uuid missing"}
		apiErr.ToJSON(rw)
		return
	}
	//Get the logs of an agent
	logs, err := lh.dbConnection.GetAgentLogsShort(uuid)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		apiErr := data.APIError{Code: data.DATABASE_ERROR, Message: err.Error()}
		apiErr.ToJSON(rw)
		return
	}
	//Create a map which will hold all the response status codes and the number of occurences
	var statusCodesOccurencesMap map[string]int64 = make(map[string]int64)

	//Parse the logs and find how many times each response status code appears
	for _, log := range logs {
		statusCode := strings.Split(log.ResponsePreview, " ")[1]
		_, ok := statusCodesOccurencesMap[statusCode]
		if !ok {
			statusCodesOccurencesMap[statusCode] = 0
		}
		statusCodesOccurencesMap[statusCode] += 1
	}
	//Log the occurences of each method for easy debugging
	//lh.logger.Debug(statusCodesOccurencesMap)

	responseData := make([]data.StatusCodesMetrics, 0)
	//Prepare the response for the client
	//Counter will be used as an ID
	var counter int64 = 0
	//Loop through all the keys and values in the occurences dictionary
	for key, value := range statusCodesOccurencesMap {
		//Create the structure which will be returned to the client
		responseData = append(responseData, data.StatusCodesMetrics{Id: counter, StatusCode: key, Count: value})
		//Increment the counter
		counter += 1
	}

	//Sort the response data based on count descending
	sort.Slice(responseData[:], func(i, j int) bool {
		return responseData[i].Count > responseData[j].Count
	})

	//Create the response structure
	resp := response.StatusCodesMetricsResponse{Metrics: responseData}

	//Send the response to the client
	rw.WriteHeader(http.StatusOK)
	resp.ToJSON(rw)
}

// Handler for getting IP addresses metrics
func (lh *LogsHandler) GetIPAddressesMetrics(rw http.ResponseWriter, r *http.Request) {
	//Get the agent uuid from the URL
	vars := mux.Vars(r)
	uuid := vars["uuid"]
	//Check if the uuid is available
	if uuid == "" {
		//Send an error message
		rw.WriteHeader(http.StatusBadRequest)
		apiErr := data.APIError{Code: data.REQUEST_ERROR, Message: "uuid missing"}
		apiErr.ToJSON(rw)
		return
	}
	//Get the logs of an agent
	logs, err := lh.dbConnection.GetAgentLogsShort(uuid)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		apiErr := data.APIError{Code: data.DATABASE_ERROR, Message: err.Error()}
		apiErr.ToJSON(rw)
		return
	}
	//Create a map which will hold all the response status codes and the number of occurences
	var ipAddressesOccurencesMap map[string]int64 = make(map[string]int64)

	//Parse the logs and find how many times each response status code appears
	for _, log := range logs {
		ipAddress := strings.Split(log.RemoteIP, ":")[0]
		_, ok := ipAddressesOccurencesMap[ipAddress]
		if !ok {
			ipAddressesOccurencesMap[ipAddress] = 0
		}
		ipAddressesOccurencesMap[ipAddress] += 1
	}
	//Log the occurences of each method for easy debugging
	//lh.logger.Debug(ipAddressesOccurencesMap)

	responseData := make([]data.IPMetrics, 0)
	//Prepare the response for the client
	//Counter will be used as an ID
	var counter int64 = 0
	//Loop through all the keys and values in the occurences dictionary
	for key, value := range ipAddressesOccurencesMap {
		//Create the structure which will be returned to the client
		responseData = append(responseData, data.IPMetrics{Id: counter, IPAddress: key, Count: value})
		//Increment the counter
		counter += 1
	}

	//Sort the response data based on count descending
	sort.Slice(responseData[:], func(i, j int) bool {
		return responseData[i].Count > responseData[j].Count
	})

	//Create the response structure
	resp := response.IPAddressMetricsResponse{Metrics: responseData}

	//Send the response to the client
	rw.WriteHeader(http.StatusOK)
	resp.ToJSON(rw)
}

// Handler for getting full date from a log
func (lh *LogsHandler) GetLog(rw http.ResponseWriter, r *http.Request) {
	//Get the agent id and the log id from the vars
	//Get the agent uuid from the URL
	vars := mux.Vars(r)
	// uuid := vars["uuid"]
	// //Check if the uuid is available
	// if uuid == "" {
	// 	//Send an error message
	// 	rw.WriteHeader(http.StatusBadRequest)
	// 	apiErr := data.APIError{Code: data.REQUEST_ERROR, Message: "uuid missing"}
	// 	apiErr.ToJSON(rw)
	// 	return
	// }
	//Get the log id from the URL
	log_uuid := vars["loguuid"]
	if log_uuid == "" {
		//Send an error message
		rw.WriteHeader(http.StatusBadRequest)
		apiErr := data.APIError{Code: data.REQUEST_ERROR, Message: "log uuid missing"}
		apiErr.ToJSON(rw)
		return
	}

	//Get the full log from the database
	log, err := lh.dbConnection.GetLog(log_uuid)
	//Check if an error occured
	if err != nil {
		//Send an error message
		rw.WriteHeader(http.StatusBadRequest)
		apiErr := data.APIError{Code: data.DATABASE_ERROR, Message: "could not retrieve the log"}
		apiErr.ToJSON(rw)
		return
	}

	//Get the log findings
	findings, err := lh.dbConnection.GetLogFindings(log_uuid)
	//Check if an error occured
	if err != nil {
		//Send an error message
		rw.WriteHeader(http.StatusBadRequest)
		apiErr := data.APIError{Code: data.DATABASE_ERROR, Message: "could not retrieve the log findings"}
		apiErr.ToJSON(rw)
		return
	}

	log.Findings = append(log.Findings, findings...)

	//Send the log back to the client
	resp := response.LogGetResponse{Log: log}
	rw.WriteHeader(http.StatusOK)
	resp.ToJSON(rw)
}

// Handler for getting finding string
func (lh *LogsHandler) GetFindingsClassificationString(rw http.ResponseWriter, r *http.Request) {
	findingsString := make([]data.FindingClassificationString, 0)
	//Create the return object
	for key, value := range data.ClassificationsMap {
		findingsString = append(findingsString, data.FindingClassificationString{IntegerFormat: key, StringFormat: value, Description: data.ClassificationDescriptionMap[key]})
	}
	respData := response.FindingClassificationStringResponse{FindingsString: findingsString}
	rw.WriteHeader(http.StatusOK)
	respData.ToJSON(rw)
}
