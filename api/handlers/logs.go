package handlers

import (
	"net/http"
	"sort"

	"github.com/gorilla/mux"
	"github.com/lucacoratu/disertatie/api/config"
	"github.com/lucacoratu/disertatie/api/data"
	"github.com/lucacoratu/disertatie/api/database"
	"github.com/lucacoratu/disertatie/api/logging"
	"github.com/lucacoratu/disertatie/api/utils"

	b64 "encoding/base64"

	response "github.com/lucacoratu/disertatie/api/data/response"
)

type LogsHandler struct {
	logger            logging.ILogger
	configuration     config.Configuration
	dbConnection      database.IConnection
	elasticConnection database.IElasticConnection
}

// Creates a new handler that will hold the functions necessary for registering proxies
func NewLogsHandler(logger logging.ILogger, configuration config.Configuration, dbConnection database.IConnection, elasticConnection database.IElasticConnection) *LogsHandler {
	return &LogsHandler{logger: logger, configuration: configuration, dbConnection: dbConnection, elasticConnection: elasticConnection}
}

func (lh *LogsHandler) GetTotalLogsCount(rw http.ResponseWriter, r *http.Request) {
	//Get the data from the elasticsearch database
	count, err := lh.elasticConnection.GetTotalCountLogs()
	//Check if an error occured when getting the total logs count
	if err != nil {
		//Send an error message
		rw.WriteHeader(http.StatusBadRequest)
		apiErr := data.APIError{Code: data.DATABASE_ERROR, Message: "could not get total logs count"}
		apiErr.ToJSON(rw)
		return
	}
	//Create the response structure
	resp := response.TotalLogCountResponse{Count: count}
	rw.WriteHeader(http.StatusOK)
	resp.ToJSON(rw)
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

	//Get the page from the request GET params
	page := r.URL.Query().Get("page")

	next_page, logs, err := lh.dbConnection.GetAgentLogsShortPaginated(uuid, page)

	// //Get the logs of an agent
	// logs, err := lh.dbConnection.GetAgentLogsShort(uuid)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		apiErr := data.APIError{Code: data.DATABASE_ERROR, Message: err.Error()}
		apiErr.ToJSON(rw)
		return
	}

	// //Sort the logs by timestamp descending (the most recent ones should be first)
	// sort.Slice(logs[:], func(i, j int) bool {
	// 	return logs[i].Timestamp > logs[j].Timestamp
	// })

	// //lh.logger.Debug(logs)

	//Send the logs back to the client
	respData := response.LogsGetResponse{Logs: logs, NextPage: next_page}
	rw.WriteHeader(http.StatusOK)
	respData.ToJSON(rw)
}

// Handler to get recent logs (10) from elasticsearch
func (lh *LogsHandler) GetRecentLogsElastic(rw http.ResponseWriter, r *http.Request) {
	//Get the recent logs from the elasticseach database
	recentLogs := lh.elasticConnection.GetRecentLogs()

	//Check if the logs could be pulled from the elasticsearch database
	if recentLogs == nil {
		rw.WriteHeader(http.StatusInternalServerError)
		apiErr := data.APIError{Code: data.DATABASE_ERROR, Message: "Failed to get recent logs"}
		apiErr.ToJSON(rw)
		return
	}

	//Add the agent name to the log struct
	for i, recentLog := range recentLogs {
		agent, err := lh.dbConnection.GetAgent(recentLog.AgentId)
		//Check if an error occured
		if err != nil {
			lh.logger.Error("Error occured when getting agent details from database", err.Error())
			continue
		}
		recentLog.AgentName = agent.Name
		recentLogs[i] = recentLog
	}

	//Send the logs back to the client
	respData := response.LogsGetResponseElastic{Logs: recentLogs}
	rw.WriteHeader(http.StatusOK)
	respData.ToJSON(rw)
}

// Handler to get recent classified logs (10) from elasticsearch
func (lh *LogsHandler) GetRecentClassifiedLogsElastic(rw http.ResponseWriter, r *http.Request) {
	//Get the recent logs from the elasticseach database
	recentLogs := lh.elasticConnection.GetRecentRuleClassifiedLogs()

	//Check if the logs could be pulled from the elasticsearch database
	if recentLogs == nil {
		rw.WriteHeader(http.StatusInternalServerError)
		apiErr := data.APIError{Code: data.DATABASE_ERROR, Message: "Failed to get recent logs"}
		apiErr.ToJSON(rw)
		return
	}

	//Add the agent name to the log struct
	for i, recentLog := range recentLogs {
		agent, err := lh.dbConnection.GetAgent(recentLog.AgentId)
		//Check if an error occured
		if err != nil {
			lh.logger.Error("Error occured when getting agent details from database", err.Error())
			continue
		}
		recentLog.AgentName = agent.Name
		recentLogs[i] = recentLog
	}

	//Send the logs back to the client
	respData := response.LogsGetResponseElastic{Logs: recentLogs}
	rw.WriteHeader(http.StatusOK)
	respData.ToJSON(rw)
}

// Handler to get all the logs from elasticsearch
func (lh *LogsHandler) GetAllLogs(rw http.ResponseWriter, r *http.Request) {
	allLogs := lh.elasticConnection.GetAllLogs()

	//Check if the logs could be pulled from elasticsearch
	if allLogs == nil {
		rw.WriteHeader(http.StatusInternalServerError)
		apiErr := data.APIError{Code: data.DATABASE_ERROR, Message: "Failed to get recent logs"}
		apiErr.ToJSON(rw)
		return
	}

	//Send the logs back to the client
	respData := response.LogsGetResponseElastic{Logs: allLogs}
	rw.WriteHeader(http.StatusOK)
	respData.ToJSON(rw)
}

// Handler to get all the classified logs from elasticsearch
func (lh *LogsHandler) GetAllClassifiedLogs(rw http.ResponseWriter, r *http.Request) {
	allLogs := lh.elasticConnection.GetAllClassifiedLogs()

	//Check if the logs could be pulled from elasticsearch
	if allLogs == nil {
		rw.WriteHeader(http.StatusInternalServerError)
		apiErr := data.APIError{Code: data.DATABASE_ERROR, Message: "Failed to get recent logs"}
		apiErr.ToJSON(rw)
		return
	}

	//Send the logs back to the client
	respData := response.LogsGetResponseElastic{Logs: allLogs}
	rw.WriteHeader(http.StatusOK)
	respData.ToJSON(rw)
}

// Handler to get all the logs from elasticsearch for a specific agent
func (lh *LogsHandler) GetLogsShortElastic(rw http.ResponseWriter, r *http.Request) {
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

	// //Get the page from the request GET params
	// page := r.URL.Query().Get("page")

	// next_page, logs, err := lh.dbConnection.GetAgentLogsShortPaginated(uuid, page)

	logs := lh.elasticConnection.GetLogsPaginated(uuid)
	var err error = nil

	// //Get the logs of an agent
	// logs, err := lh.dbConnection.GetAgentLogsShort(uuid)
	if err != nil {
		rw.WriteHeader(http.StatusInternalServerError)
		apiErr := data.APIError{Code: data.DATABASE_ERROR, Message: err.Error()}
		apiErr.ToJSON(rw)
		return
	}

	// //Sort the logs by timestamp descending (the most recent ones should be first)
	// sort.Slice(logs[:], func(i, j int) bool {
	// 	return logs[i].Timestamp > logs[j].Timestamp
	// })

	// //lh.logger.Debug(logs)

	//Send the logs back to the client
	respData := response.LogsGetResponseElastic{Logs: logs}
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

	//Define the list of possible methods in the HTTP request (https://developer.mozilla.org/en-US/docs/Web/HTTP/Methods)
	methods := []string{"GET", "HEAD", "POST", "PUT", "DELETE", "CONNECT", "OPTIONS", "TRACE", "PATCH"}
	//Create a map which will hold all the request methods and the number of occurences
	var methodsOccurencesMap map[string]int64 = make(map[string]int64)

	for _, method := range methods {
		//Get the count for the method from database
		count, err := lh.dbConnection.GetLogsMethodCount(uuid, method)
		if err != nil {
			lh.logger.Error("Error occured when getting count for method", method, err.Error())
			//Send an error message back to the client
		}
		methodsOccurencesMap[method] = count
	}

	responseData := make([]data.MethodsMetrics, 0)
	//Prepare the response for the client
	//Counter will be used as an ID
	var counter int64 = 0
	//Loop through all the keys and values in the occurences dictionary
	for key, value := range methodsOccurencesMap {
		//Add only the methods which have a count > 0
		if value > 0 {
			//Create the structure which will be returned to the client
			responseData = append(responseData, data.MethodsMetrics{Id: counter, Method: key, Count: value})
			//Increment the counter
			counter += 1
		}
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

	mapOccurencesPerDay, _ := lh.dbConnection.GetRequestsPerDay(uuid)

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

	statusCodesOccurencesMap, _ := lh.dbConnection.GetStatusCodeCounts(uuid)

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

	ipAddressesOccurencesMap, _ := lh.dbConnection.GetIPAddressesCounts(uuid)

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

	//Check if there are more than 20 ip addresses in the list of metrics, if there are then truncate it to first 20
	if len(responseData) > 20 {
		responseData = responseData[:20]
	}

	//Create the response structure
	resp := response.IPAddressMetricsResponse{Metrics: responseData}

	//Send the response to the client
	rw.WriteHeader(http.StatusOK)
	resp.ToJSON(rw)
}

// Handler for getting IP addresses metrics
func (lh *LogsHandler) GetAllIPAddressesMetrics(rw http.ResponseWriter, r *http.Request) {
	ipAddressesOccurencesMap, _ := lh.dbConnection.GetAllIPAddressesCounts()

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

	//Check if there are more than 20 ip addresses in the list of metrics, if there are then truncate it to first 20
	if len(responseData) > 20 {
		responseData = responseData[:20]
	}

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

	//Get the log rule findings
	ruleFindings, err := lh.dbConnection.GetLogRuleFindings(log_uuid)
	//Check if an error occured
	if err != nil {
		//Send an error message
		rw.WriteHeader(http.StatusBadRequest)
		apiErr := data.APIError{Code: data.DATABASE_ERROR, Message: "could not retrieve the log rule findings"}
		apiErr.ToJSON(rw)
		return
	}

	log.Findings = append(log.Findings, findings...)
	log.RuleFindings = append(log.RuleFindings, ruleFindings...)

	//Send the log back to the client
	resp := response.LogGetResponse{Log: log}
	lh.logger.Debug(log.RuleFindings)
	rw.WriteHeader(http.StatusOK)
	resp.ToJSON(rw)
}

// Handler for getting the exploit code of the log
func (lh *LogsHandler) GetLogExploitPythonCode(rw http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	//Get the log id from the URL
	log_uuid := vars["loguuid"]
	if log_uuid == "" {
		//Send an error message
		rw.WriteHeader(http.StatusBadRequest)
		apiErr := data.APIError{Code: data.REQUEST_ERROR, Message: "log uuid missing"}
		apiErr.ToJSON(rw)
		return
	}

	//Check if the exploit code exists for the log in the database
	exploitExists, err := lh.dbConnection.CheckExploitCodeExists(log_uuid)
	if err != nil {
		//Send an error message
		rw.WriteHeader(http.StatusBadRequest)
		apiErr := data.APIError{Code: data.DATABASE_ERROR, Message: "could not retrieve the number of exploit codes from database"}
		apiErr.ToJSON(rw)
		return
	}

	//If the exploit code does not exist then create it
	if !exploitExists {
		//Get the request from the database
		raw_request, err := lh.dbConnection.GetLogRequest(log_uuid)
		if err != nil {
			//Send an error message
			rw.WriteHeader(http.StatusBadRequest)
			apiErr := data.APIError{Code: data.DATABASE_ERROR, Message: "could not retrieve the log request from database, " + err.Error()}
			apiErr.ToJSON(rw)
			return
		}
		exploitCode, err := utils.CreatePythonExploitCode(raw_request, lh.configuration.ExploitTemplatePath)
		if err != nil {
			//Send an error message
			rw.WriteHeader(http.StatusBadRequest)
			apiErr := data.APIError{Code: data.DATABASE_ERROR, Message: "could not create the exploit code, " + err.Error()}
			apiErr.ToJSON(rw)
			return
		}

		//lh.logger.Debug(exploitCode)
		encExploitCode := b64.StdEncoding.EncodeToString([]byte(exploitCode))
		exploitResponse := response.ExploitResponse{Exploit: encExploitCode}

		rw.WriteHeader(http.StatusOK)
		exploitResponse.ToJSON(rw)
	}
	// //Save the newly created exploit code to the database

	// if exploitExists {

	// }
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

func (lh *LogsHandler) GetLogsRuleFindingsMetrics(rw http.ResponseWriter, r *http.Request) {
	//Get the metrics from elasticsearch
	metrics, err := lh.elasticConnection.GetRuleFindingsStats()
	if err != nil {
		//Send an error message
		rw.WriteHeader(http.StatusBadRequest)
		apiErr := data.APIError{Code: data.DATABASE_ERROR, Message: "could not retrieve the metrics for rule findings"}
		apiErr.ToJSON(rw)
		return
	}
	//Send the metrics back to the client
	resp := response.FindingsMetricsResponse{Metrics: metrics}
	rw.WriteHeader(http.StatusOK)
	resp.ToJSON(rw)
}

func (lh *LogsHandler) GetLogsRuleIdMetrics(rw http.ResponseWriter, r *http.Request) {
	//Get the metrics from elasticsearch
	metrics, err := lh.elasticConnection.GetRuleIdStats()
	if err != nil {
		//Send an error message
		rw.WriteHeader(http.StatusBadRequest)
		apiErr := data.APIError{Code: data.DATABASE_ERROR, Message: "could not retrieve the metrics for rule ids"}
		apiErr.ToJSON(rw)
		return
	}
	//Send the metrics back to the client
	resp := response.FindingsMetricsResponse{Metrics: metrics}
	rw.WriteHeader(http.StatusOK)
	resp.ToJSON(rw)
}

func (lh *LogsHandler) GetFindingsCount(rw http.ResponseWriter, r *http.Request) {
	//Get the metrics from elasticsearch
	metrics, err := lh.elasticConnection.GetFindingsStats()
	if err != nil {
		//Send an error message
		rw.WriteHeader(http.StatusBadRequest)
		apiErr := data.APIError{Code: data.DATABASE_ERROR, Message: "could not retrieve the findings count metrics"}
		apiErr.ToJSON(rw)
		return
	}
	//Send the metrics back to the client
	resp := response.FindingsCountMetricsResponse{Metrics: metrics}
	rw.WriteHeader(http.StatusOK)
	resp.ToJSON(rw)
}

// Handler for getting classification metrics
func (lh *LogsHandler) GetClassificationMetrics(rw http.ResponseWriter, r *http.Request) {
	//Get the findings count from the database
	findingsMetrics, err := lh.elasticConnection.GetFindingsStats()
	if err != nil {
		//Send an error message
		rw.WriteHeader(http.StatusBadRequest)
		apiErr := data.APIError{Code: data.DATABASE_ERROR, Message: "could not retrieve the findings count metrics"}
		apiErr.ToJSON(rw)
		return
	}
	var classifiedCount int64 = findingsMetrics.FindingsCount + findingsMetrics.RuleFindingsCount
	totalLogs, err := lh.elasticConnection.GetTotalCountLogs()
	if err != nil {
		//Send an error message
		rw.WriteHeader(http.StatusBadRequest)
		apiErr := data.APIError{Code: data.DATABASE_ERROR, Message: "could not retrieve the total logs count"}
		apiErr.ToJSON(rw)
		return
	}

	totalLogs = totalLogs - classifiedCount
	metrics := data.ClassificationMetrics{ClassifiedCount: classifiedCount, UnclassifiedCount: totalLogs}
	response := response.ClassificationMetricsResponse{Metrics: metrics}
	rw.WriteHeader(http.StatusOK)
	response.ToJSON(rw)
}

// Handler for getting agent log counts metrics
func (lh *LogsHandler) GetAgentsMetrics(rw http.ResponseWriter, r *http.Request) {
	//Get the agent metrics from elasticsearch
	metrics, err := lh.elasticConnection.GetAgentsStatistics()
	if err != nil {
		//Send an error message
		rw.WriteHeader(http.StatusBadRequest)
		apiErr := data.APIError{Code: data.DATABASE_ERROR, Message: "could not retrieve the agents logs count metrics"}
		apiErr.ToJSON(rw)
		return
	}

	for index, metric := range metrics {
		//Add the agent name to the structure from cassandra
		agent, err := lh.dbConnection.GetAgent(metric.AgentId)
		if err != nil {
			//Log an error message
			lh.logger.Error("Could not get the agent name from cassandra for agent", metric.AgentId)
			continue
		}
		//Add the name to the metric structure
		metric.AgentName = agent.Name
		metrics[index] = metric
	}
	lh.logger.Debug(metrics)

	//Send the response to the client
	response := response.AgentsMetricsResponse{Metrics: metrics}
	rw.WriteHeader(http.StatusOK)
	response.ToJSON(rw)
}

func (lh *LogsHandler) ExportAgentLogs(rw http.ResponseWriter, r *http.Request) {
	//Get the format to be exported into
	format := r.URL.Query().Get("format")
	if format == "" {
		//Default format is json
		format = "json"
	}

	//Check if the format specified is correct
	if format != "json" || format != "csv" {
		//Send an error message
		rw.WriteHeader(http.StatusBadRequest)
		apiErr := data.APIError{Code: data.PARSE_ERROR, Message: "format should be json or csv"}
		apiErr.ToJSON(rw)
		return
	}

	//Get the agent id from the mux vars
	vars := mux.Vars(r)
	agentId := vars["uuid"]

	if agentId == "" {
		//Send an error message
		rw.WriteHeader(http.StatusBadRequest)
		apiErr := data.APIError{Code: data.PARSE_ERROR, Message: "agentId is required"}
		apiErr.ToJSON(rw)
		return
	}

	//Get the logs from the database
	agentLogs := lh.elasticConnection.GetAllAgentLogs(agentId)
	logs := response.LogsGetResponseElastic{Logs: agentLogs}

	rw.WriteHeader(http.StatusOK)
	rw.Header().Set("Content-Disposition", "attachment; filename=export."+format)
	rw.Header().Set("Content-Type", "application/"+format)
	logs.ToJSON(rw)
}
