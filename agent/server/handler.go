package server

import (
	"bytes"
	b64 "encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"

	"github.com/lucacoratu/disertatie/agent/api"
	"github.com/lucacoratu/disertatie/agent/config"
	"github.com/lucacoratu/disertatie/agent/data"
	code "github.com/lucacoratu/disertatie/agent/detection/code"
	rules "github.com/lucacoratu/disertatie/agent/detection/rules"
	"github.com/lucacoratu/disertatie/agent/logging"
	"github.com/lucacoratu/disertatie/agent/utils"
	"github.com/lucacoratu/disertatie/agent/websocket"
)

/*
 * Structure which holds all the information needed by the handler for the HTTP requests
 */
type AgentHandler struct {
	logger        logging.ILogger                   //The logger interface
	apiBaseURL    string                            //The API base URL
	configuration config.Configuration              //The configuration structure
	checkers      []code.IValidator                 //The list of validators which will be run on the request and the response to find malicious activity
	rules         []rules.Rule                      //The list of rules which will try to find anomalies in the requests and the responses
	apiWsConn     *websocket.APIWebSocketConnection //The WS connection to the API
}

// Creates a new AgentHandlerStructure
func NewAgentHandler(logger logging.ILogger, apiBaseURL string, configuration config.Configuration, checkers []code.IValidator, rules []rules.Rule, apiWsConn *websocket.APIWebSocketConnection) *AgentHandler {
	return &AgentHandler{logger: logger, apiBaseURL: apiBaseURL, configuration: configuration, checkers: checkers, rules: rules, apiWsConn: apiWsConn}
}

// Forwards the request to the target server
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

// Forwards the response back to the client
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

// Combines the request and response findings into a single slice
func (agentHandler *AgentHandler) combineFindings(requestFindings []data.FindingData, responseFindings []data.FindingData) []data.Finding {
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

	return allFindings
}

func (agentHandler *AgentHandler) combineRuleFindings(requestRuleFindings []*data.RuleFindingData, responseRuleFindings []*data.RuleFindingData) []data.RuleFinding {
	//Add all the findings from all the validators to a list which will be sent to the API
	allFindings := make([]data.RuleFinding, 0)
	//Add all request findings
	for index, finding := range requestRuleFindings {
		if index < len(responseRuleFindings) {
			allFindings = append(allFindings, data.RuleFinding{Request: finding, Response: responseRuleFindings[index]})
		} else {
			allFindings = append(allFindings, data.RuleFinding{Request: finding, Response: nil})
		}
	}

	//Add the response findings
	for index, finding := range responseRuleFindings {
		//If the index is less than the length of the all findings list then complete the index structure with the response findings
		if index < len(allFindings) {
			allFindings[index].Response = finding
		} else {
			//Otherwise add a new structure to the list of all findings which will have the Request empty
			allFindings = append(allFindings, data.RuleFinding{Request: nil, Response: finding})
		}
	}

	return allFindings
}

// Converts the request and the response to raw string then base64 encodes both of them
func (agentHandler *AgentHandler) convertRequestAndResponseToB64(req *http.Request, resp *http.Response) (string, string, error) {
	//Dump the HTTP request to raw string
	rawRequest, _ := utils.DumpHTTPRequest(req)
	//Dump the response as raw string
	rawResponse, err := utils.DumpHTTPResponse(resp)
	//Check if an error occured when dumping the response as raw string
	if err != nil {
		agentHandler.logger.Error(err.Error())
		return "", "", err
	}
	//Convert raw request to base64
	b64RawRequest := b64.StdEncoding.EncodeToString(rawRequest)
	//Convert raw response to base64
	b64RawResponse := b64.StdEncoding.EncodeToString(rawResponse)
	//Return the base64 string of the request and the response
	return b64RawRequest, b64RawResponse, nil
}

// Sends the raw request encoded base64 to the LLM API
// Returns LLMResponseData which will contain a series of headers and a body generated by the LLM
func (agentHandler *AgentHandler) sendB64RequestToLLMAPI(req *http.Request) *data.LLMResponse {
	//Convert raw HTTP request to base64
	rawRequest, _ := utils.DumpHTTPRequest(req)
	b64RawRequest := b64.StdEncoding.EncodeToString(rawRequest)

	//Create the request to the LLM API
	requestURL := fmt.Sprintf("%s/generic?raw_request=%s", agentHandler.configuration.LLMAPIURL, b64RawRequest)
	llm_req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		agentHandler.logger.Error("Failed to create request to LLM API", err.Error())
		return nil
	}

	//Send the request to the LLM API
	llm_response, err := http.DefaultClient.Do(llm_req)
	if err != nil {
		agentHandler.logger.Error("Failed to send request to LLM API", err.Error())
		return nil
	}

	//Parse the response from the LLM
	var llm_response_data data.LLMResponse
	err = llm_response_data.FromJSON(llm_response.Body)
	if err != nil {
		agentHandler.logger.Error("Failed to decode JSON from LLM API response", err.Error())
		return nil
	}

	return &llm_response_data
}

// Handles the requests received by the agent
func (agentHandler *AgentHandler) HandleRequest(rw http.ResponseWriter, r *http.Request) {
	//Log the endpoint where the request was made
	agentHandler.logger.Info("Received request on", r.URL.Path)

	//Create the validator runner
	validatorRunner := code.NewValidatorRunner(agentHandler.checkers, agentHandler.logger)
	//Create the rule runner
	ruleRunner := rules.NewRuleRunner(agentHandler.logger, agentHandler.rules, agentHandler.apiWsConn, agentHandler.configuration)

	//Run all the validators on the request
	requestFindings, _ := validatorRunner.RunValidatorsOnRequest(r)
	//Run all the rules on the request
	requestRuleFindings, _ := ruleRunner.RunRulesOnRequest(r)

	//Log request findings
	agentHandler.logger.Debug("Request findings", requestFindings)
	//Log the request rule findings
	agentHandler.logger.Debug("Request rule findings", requestRuleFindings)

	//If the mode of operation is adaptive then send the raw request encoded base64 to LLM
	if agentHandler.configuration.OperationMode == "adaptive" {
		//Check if the request should be sent to the LLM API
		//If it shouldn't be sent then serve a static page

		//Send the request
		agentHandler.logger.Debug("Sending request to LLM API...")
		llm_response_data := agentHandler.sendB64RequestToLLMAPI(r)

		//Debug log the response from the LLM API
		agentHandler.logger.Debug(*llm_response_data)

		//Add the headers to the response
		for header_name, header_value := range llm_response_data.Headers {
			rw.Header().Set(header_name, header_value)
		}

		//Send the response back to the client
		rw.WriteHeader(http.StatusOK)

		agentHandler.logger.Debug("Sending body")
		//Send the body
		rw.Write([]byte(llm_response_data.Body))
		return
	}

	//Forward the request to the destination web server
	response, err := agentHandler.forwardRequest(r)
	if err != nil {
		agentHandler.logger.Error(err.Error())
		return
	}

	//Run the validators on the response
	responseFindings, _ := validatorRunner.RunValidatorsOnResponse(response)
	//Run the rules on the response
	responseRuleFindings, _ := ruleRunner.RunRulesOnResponse(response)

	//Log response findings
	agentHandler.logger.Debug("Response findings", responseFindings)
	//Log the rules response findings
	agentHandler.logger.Debug("Response rule findings", responseRuleFindings)

	//Combine the findings into a single structure
	allFindings := agentHandler.combineFindings(requestFindings, responseFindings)
	//Combine the rule findings into a single structure
	allRuleFindings := agentHandler.combineRuleFindings(requestRuleFindings, responseRuleFindings)

	//Convert the request and response to base64 string
	b64RawRequest, b64RawResponse, _ := agentHandler.convertRequestAndResponseToB64(r, response)

	//Create the log structure that should be sent to the API
	logData := data.LogData{AgentId: agentHandler.configuration.UUID, RemoteIP: r.RemoteAddr, Timestamp: time.Now().Unix(), Request: b64RawRequest, Response: b64RawResponse, Findings: allFindings, RuleFindings: allRuleFindings}
	if true {
		agentHandler.logger.Debug("Log data", logData)
	}

	if agentHandler.apiWsConn != nil {
		//Send log information to the API
		apiHandler := api.NewAPIHandler(agentHandler.logger, agentHandler.configuration)
		_, err = apiHandler.SendLog(agentHandler.apiBaseURL, logData)
		//Check if an error occured when sending log to the API
		if err != nil {
			agentHandler.logger.Error(err.Error())
			//return
		}
	}

	//If the mode is testing then send the log data as response
	if strings.EqualFold(agentHandler.configuration.OperationMode, "testing") {
		rw.WriteHeader(http.StatusOK)
		logData.ToJSON(rw)
		return
	}

	//Send the response from the web server back to the client
	agentHandler.forwardResponse(rw, response)
}
