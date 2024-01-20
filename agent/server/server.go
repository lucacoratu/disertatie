package server

import (
	"bytes"
	"context"
	b64 "encoding/base64"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/lucacoratu/disertatie/agent/api"
	"github.com/lucacoratu/disertatie/agent/config"
	"github.com/lucacoratu/disertatie/agent/data"
	"github.com/lucacoratu/disertatie/agent/logging"
	"github.com/lucacoratu/disertatie/agent/preliminary"
	rules "github.com/lucacoratu/disertatie/agent/preliminary/rules"
	"github.com/lucacoratu/disertatie/agent/runtimedata"
	"github.com/lucacoratu/disertatie/agent/utils"
)

type AgentServer struct {
	srv              *http.Server
	logger           logging.ILogger
	collectorBaseURL string
	configuration    config.Configuration
	checkers         []preliminary.IValidator
	runtimeData      runtimedata.RuntimeData
	rules            []rules.YamlRule
}

func (agent *AgentServer) ForwardRequest(req *http.Request) (*http.Response, error) {
	// we need to buffer the body if we want to read it here and send it
	// in the request.
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, errors.New("could not send the request to the target web server, " + err.Error())
	}

	// you can reassign the body if you need to parse it as multipart
	req.Body = io.NopCloser(bytes.NewReader(body))

	// create a new url from the raw RequestURI sent by the client
	url := fmt.Sprintf("%s://%s%s", agent.configuration.ForwardServerProtocol, agent.configuration.ForwardServerAddress+":"+agent.configuration.ForwardServerPort, req.RequestURI)

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

func (agent *AgentServer) ForwardResponse(rw http.ResponseWriter, response *http.Response) {
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

func (agent *AgentServer) handleRequest(rw http.ResponseWriter, r *http.Request) {
	agent.logger.Debug("Received request on", r.URL.Path)
	//Create the raw HTTP request string
	rawRequest, _ := utils.DumpHTTPRequest(r)
	agent.logger.Debug(string(rawRequest))

	//Create the list of findings
	requestFindings := make([]data.FindingData, 0)

	//Run all the validators to check if the request seems valid
	for _, valid := range agent.checkers {
		//Call the validate method of the validator for the request
		findingsRequest, err := valid.ValidateRequest(r)
		//Check if there was an error when looking for malicious input in the request
		if err != nil {
			agent.logger.Error("Error occured when trying to find malicious input in the request", err.Error())
		}
		//Check if there was a finding returned by the validator
		if findingsRequest != nil {
			//Add the findings to the list of findings for the request
			requestFindings = append(requestFindings, findingsRequest...)
			//Log the findings
			for _, finding := range findingsRequest {
				agent.logger.Debug(finding)
			}
		}

		// //If the request is not valid then
		// if !validRequest {
		// 	//Send a forbidden page
		// 	rw.WriteHeader(http.StatusForbidden)
		// 	file, err := os.Open(proxy.configuration.ForbiddenPagePath)
		// 	if err != nil {
		// 		proxy.logger.Error("could not send forbidden page", err.Error())
		// 		return
		// 	}
		// 	_, err = io.Copy(rw, file)
		// 	if err != nil {
		// 		proxy.logger.Error("could not copy the contents of forbidden page to response writer", err.Error())
		// 		return
		// 	}
		// 	return
		// }
	}

	//Log request findings
	agent.logger.Debug("Request findings", requestFindings)

	//Forward the request to the destination web server
	response, err := agent.ForwardRequest(r)
	if err != nil {
		agent.logger.Error(err.Error())
		return
	}

	//Create the structure which will hold response findings
	responseFindings := make([]data.FindingData, 0)

	//Run all the validators to check if the response seems valid
	for _, valid := range agent.checkers {
		//Call the validate method of the validator for the response
		findingsResponse, err := valid.ValidateResponse(response)
		//Check if an error occured when looking for malicious input in the response
		if err != nil {
			agent.logger.Error("Error occured when trying to find malicious input in the response", err.Error())
		}
		//Check if there was any finding
		if findingsResponse != nil {
			//Add the finding to the list of response findings
			responseFindings = append(responseFindings, findingsResponse...)
			//Log the findings
			for _, finding := range findingsResponse {
				agent.logger.Debug(finding)
			}
		}
	}

	//Log response findings
	agent.logger.Debug("Response findings", responseFindings)

	//Dump the response as raw string
	rawResponse, err := utils.DumpHTTPResponse(response)
	//Check if an error occured when dumping the response as raw string
	if err != nil {
		agent.logger.Error(err.Error())
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
	logData := data.LogData{AgentId: agent.runtimeData.Uuid, RemoteIP: r.RemoteAddr, Timestamp: time.Now().Unix(), Request: b64RawRequest, Response: b64RawResponse, Findings: allFindings}
	agent.logger.Debug(logData)
	//Send log information to the API
	apiHandler := api.NewAPIHandler(agent.logger, agent.configuration)
	_, err = apiHandler.SendLog(agent.collectorBaseURL, logData)
	//Check if an error occured when sending log to the API
	if err != nil {
		agent.logger.Error(err.Error())
		return
	}

	//Send the response from the web server back to the client
	agent.ForwardResponse(rw, response)
}

// Initialize the proxy http server based on the configuration file
func (agent *AgentServer) Init() error {
	//Initialize the logger
	agent.logger = logging.NewDefaultDebugLogger()
	agent.logger.Info("Logger initialized")

	//Load the configuration from file
	err := agent.configuration.LoadConfigurationFromFile(".\\agent.conf")
	if err != nil {
		agent.logger.Fatal("Error occured when loading the config from file,", err.Error())
		return err
	}
	agent.logger.Info("Loaded configuration from file")

	//Check if the rules directory was specified in the configuration file
	if agent.configuration.RulesDirectory != "" {
		//Load the rules from the rules directory
		allRules, err := rules.LoadRulesFromDirectory(agent.configuration.RulesDirectory, agent.logger)
		if err != nil {
			agent.logger.Error("Could not load rules from", agent.configuration.RulesDirectory, err.Error())
		}
		agent.logger.Info("Loaded", len(allRules), "rules from", agent.configuration.RulesDirectory)
		//Add the list of rules to the server structure
		agent.rules = allRules
	} else {
		agent.logger.Warning("No rules were loaded because the rules directory was not specified")
		//Assign nil to the rules slice of the server structure
		agent.rules = nil
	}

	//Assemble the collector base URL
	agent.collectorBaseURL = agent.configuration.APIProtocol + "://" + agent.configuration.APIIpAddress + ":" + agent.configuration.APIPort + "/api/v1"

	//Check connection to the collector
	if !utils.CheckCollectorConnection(agent.collectorBaseURL) {
		agent.logger.Error("Cannot connect to the API")
		return errors.New("could not connect to the API")
	}

	apiHandler := api.NewAPIHandler(agent.logger, agent.configuration)

	//Check if the proxy instance isn't already registered in the collector
	exists := utils.CheckFileExists(".\\uuid.txt")
	if !exists {
		//Collect information of the operating system
		machineInfo, err := utils.GetMachineInfo()
		if err != nil {
			agent.logger.Error(err.Error())
			return err
		}

		//Log the machine info extracted
		agent.logger.Debug(machineInfo)

		//Populate information about this agent
		agentInfo := data.AgentInformation{Protocol: agent.configuration.ListeningProtocol, IPAddress: agent.configuration.ListeningAddress, Port: agent.configuration.ListeningPort, WebServerProtocol: agent.configuration.ForwardServerProtocol, WebServerIP: agent.configuration.ForwardServerAddress, WebServerPort: agent.configuration.ForwardServerPort, MachineInfo: machineInfo}

		//Send the information to the collector
		uuid, err := apiHandler.RegisterAgent(agent.collectorBaseURL, agentInfo)
		if err != nil {
			agent.logger.Error("could not register this proxy on the collector", err.Error())
			return err
		}

		agent.logger.Debug("UUID received", uuid)
		//Save the uuid into uuid.txt
		uuidFile, err := os.OpenFile(".\\uuid.txt", os.O_APPEND|os.O_WRONLY|os.O_CREATE, 0600)
		//Check if an error occured when opening the uuid.txt file
		if err != nil {
			agent.logger.Error("could not open uuid file", err.Error())
			return err
		}
		_, err = uuidFile.WriteString(uuid)
		//Check if an error occured when writing the uuid into the file
		if err != nil {
			agent.logger.Error("could not write uuid in the uuid file", err.Error())
			return err
		}

		//Save the uuid in the runtime data
		agent.runtimeData.Uuid = uuid
	} else {
		//Read the uuid from file
		uuid, err := utils.ReadAllDataFromFile(".\\uuid.txt")
		if err != nil {
			agent.logger.Error("could not read uuid from uuid file", err.Error())
			return err
		}
		//Save the uuid in the runtime data
		agent.runtimeData.Uuid = uuid
		agent.logger.Debug("Loaded UUID from file", agent.runtimeData.Uuid)
	}

	//Add the validators to the list of validators
	agent.checkers = append(agent.checkers, preliminary.NewUserAgentValidator(agent.logger, agent.configuration))

	//Create the router
	r := mux.NewRouter()

	//Create a single route that will catch every request on every method
	r.PathPrefix("/").HandlerFunc(agent.handleRequest)

	agent.srv = &http.Server{
		Addr: agent.configuration.ListeningAddress + ":" + agent.configuration.ListeningPort,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	return nil
}

// Start the proxy server
func (agent *AgentServer) Run() {
	var wait time.Duration = 5
	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := agent.srv.ListenAndServe(); err != nil {
			agent.logger.Error(err.Error())
		}
	}()

	agent.logger.Info("Started server on port", agent.configuration.ListeningPort)

	c := make(chan os.Signal, 1)
	// We'll accept graceful shutdowns when quit via SIGINT (Ctrl+C)
	// SIGKILL, SIGQUIT or SIGTERM (Ctrl+/) will not be caught.
	signal.Notify(c, os.Interrupt)

	// Block until we receive our signal.
	<-c

	// Create a deadline to wait for.
	ctx, cancel := context.WithTimeout(context.Background(), wait)
	defer cancel()
	// Doesn't block if no connections, but will otherwise wait
	// until the timeout deadline.
	agent.srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	agent.logger.Info("shutting down")
	os.Exit(0)
}
