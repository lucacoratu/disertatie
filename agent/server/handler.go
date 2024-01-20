package server

import (
	"bytes"
	b64 "encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"time"

	"github.com/lucacoratu/disertatie/agent/api"
	"github.com/lucacoratu/disertatie/agent/config"
	"github.com/lucacoratu/disertatie/agent/data"
	code "github.com/lucacoratu/disertatie/agent/detection/code"
	rules "github.com/lucacoratu/disertatie/agent/detection/rules"
	"github.com/lucacoratu/disertatie/agent/logging"
	"github.com/lucacoratu/disertatie/agent/utils"
)

type AgentHandler struct {
	logger        logging.ILogger
	apiBaseURL    string
	configuration config.Configuration
	checkers      []code.IValidator
	rules         []rules.YamlRule
}

func NewAgentHandler(logger logging.ILogger, apiBaseURL string, configuration config.Configuration, checkers []code.IValidator, rules []rules.YamlRule) *AgentHandler {
	return &AgentHandler{logger: logger, apiBaseURL: apiBaseURL, configuration: configuration, checkers: checkers, rules: rules}
}

func (agentHandler *AgentHandler) forwardRequest(req *http.Request) (*http.Response, error) {
	// we need to buffer the body if we want to read it here and send it
	// in the request.
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, errors.New("could not send the request to the target web server, " + err.Error())
	}

	// you can reassign the body if you need to parse it as multipart
	req.Body = io.NopCloser(bytes.NewReader(body))

	// create a new url from the raw RequestURI sent by the client
	url := fmt.Sprintf("%s://%s%s", agentHandler.configuration.ForwardServerProtocol, agentHandler.configuration.ForwardServerAddress+":"+agentHandler.configuration.ForwardServerPort, req.RequestURI)

	proxyReq, err := http.NewRequest(req.Method, url, bytes.NewReader(body))
	if err != nil {
		return nil, errors.New("could not create the new request to forward to target web server")
	}

	proxyReq.Header = make(http.Header)
	for h, val := range req.Header {
		proxyReq.Header[h] = val
	}

	httpClient := &http.Client{}

	resp, err := httpClient.Do(proxyReq)
	if err != nil {
		return nil, errors.New("could not send the request to the target web server, " + err.Error())
	}

	return resp, nil
}

func (agentHandler *AgentHandler) forwardResponse(rw http.ResponseWriter, response *http.Response) {
	//Send the status code
	rw.WriteHeader(response.StatusCode)
	//Send the headers
	for name, values := range response.Header {
		val := ""
		for _, value := range values {
			val += value
			if len(values) > 1 {
				val += ";"
			}
		}
		rw.Header().Set(name, val)
	}
	//Send the body
	body, err := io.ReadAll(response.Body)
	//agent.logger.Debug("Body:", body)

	if err != nil {
		rw.Write([]byte("error"))
		return
	}
	rw.Write(body)
}

func (agentHandler *AgentHandler) HandleRequest(rw http.ResponseWriter, r *http.Request) {
	agentHandler.logger.Debug("Received request on", r.URL.Path)
	//Create the raw HTTP request string
	rawRequest, _ := utils.DumpHTTPRequest(r)
	//agent.logger.Debug(string(rawRequest))

	//Create the list of findings
	requestFindings := make([]data.FindingData, 0)

	//Run all the validators to check if the request seems valid
	for _, valid := range agentHandler.checkers {
		//Call the validate method of the validator for the request
		findingsRequest, err := valid.ValidateRequest(r)
		//Check if there was an error when looking for malicious input in the request
		if err != nil {
			agentHandler.logger.Error("Error occured when trying to find malicious input in the request", err.Error())
		}
		//Check if there was a finding returned by the validator
		if findingsRequest != nil {
			//Add the findings to the list of findings for the request
			requestFindings = append(requestFindings, findingsRequest...)
			//Log the findings
			for _, finding := range findingsRequest {
				agentHandler.logger.Debug(finding)
			}
		}
	}

	//Log request findings
	agentHandler.logger.Debug("Request findings", requestFindings)

	//Forward the request to the destination web server
	response, err := agentHandler.forwardRequest(r)
	if err != nil {
		agentHandler.logger.Error(err.Error())
		return
	}

	//Create the structure which will hold response findings
	responseFindings := make([]data.FindingData, 0)

	//Run all the validators to check if the response seems valid
	for _, valid := range agentHandler.checkers {
		//Call the validate method of the validator for the response
		findingsResponse, err := valid.ValidateResponse(response)
		//Check if an error occured when looking for malicious input in the response
		if err != nil {
			agentHandler.logger.Error("Error occured when trying to find malicious input in the response", err.Error())
		}
		//Check if there was any finding
		if findingsResponse != nil {
			//Add the finding to the list of response findings
			responseFindings = append(responseFindings, findingsResponse...)
			//Log the findings
			for _, finding := range findingsResponse {
				agentHandler.logger.Debug(finding)
			}
		}
	}

	//Log response findings
	agentHandler.logger.Debug("Response findings", responseFindings)

	//Dump the response as raw string
	rawResponse, err := utils.DumpHTTPResponse(response)
	//Check if an error occured when dumping the response as raw string
	if err != nil {
		agentHandler.logger.Error(err.Error())
		return
	}
	//Convert raw request to base64
	b64RawRequest := b64.StdEncoding.EncodeToString(rawRequest)
	//Convert raw response to base64
	b64RawResponse := b64.StdEncoding.EncodeToString(rawResponse)

	//Add all the findings from all the validators to a list which will be sent to the API
	allFindings := make([]data.Finding, 0)
	//Add all request findings
	for index, finding := range requestFindings {
		if index < len(responseFindings) {
			allFindings = append(allFindings, data.Finding{Request: finding, Response: responseFindings[index]})
		} else {
			allFindings = append(allFindings, data.Finding{Request: finding, Response: data.FindingData{}})
		}
	}

	//Add the response findings
	for index, finding := range responseFindings {
		//If the index is less than the length of the all findings list then complete the index structure with the response findings
		if index < len(allFindings) {
			allFindings[index].Response = finding
		} else {
			//Otherwise add a new structure to the list of all findings which will have the Request empty
			allFindings = append(allFindings, data.Finding{Request: data.FindingData{}, Response: finding})
		}
	}

	//Create the log structure that should be sent to the API
	//logData := data.LogData{AgentId: agent.runtimeData.Uuid, RemoteIP: r.RemoteAddr, Timestamp: time.Now().Unix(), Request: b64RawRequest, Response: b64RawResponse, Findings: make([]data.Finding, 0)}
	logData := data.LogData{AgentId: agentHandler.configuration.UUID, RemoteIP: r.RemoteAddr, Timestamp: time.Now().Unix(), Request: b64RawRequest, Response: b64RawResponse, Findings: allFindings}
	agentHandler.logger.Debug(logData)
	//Send log information to the API
	apiHandler := api.NewAPIHandler(agentHandler.logger, agentHandler.configuration)
	_, err = apiHandler.SendLog(agentHandler.apiBaseURL, logData)
	//Check if an error occured when sending log to the API
	if err != nil {
		agentHandler.logger.Error(err.Error())
		return
	}

	//Send the response from the web server back to the client
	agentHandler.forwardResponse(rw, response)
}
