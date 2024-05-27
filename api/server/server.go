package server

import (
	"context"
	"flag"
	"net/http"
	"os"
	"os/signal"
	"time"

	"github.com/gorilla/mux"
	"github.com/lucacoratu/disertatie/api/config"
	"github.com/lucacoratu/disertatie/api/database"
	"github.com/lucacoratu/disertatie/api/handlers"
	"github.com/lucacoratu/disertatie/api/jwt"
	"github.com/lucacoratu/disertatie/api/logging"
	"github.com/lucacoratu/disertatie/api/websocket"
)

type APIServer struct {
	srv               *http.Server
	logger            logging.ILogger
	configuration     config.Configuration
	dbConnection      database.IConnection
	elasticConnection database.IElasticConnection
	configFile        string
}

func (api *APIServer) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api.logger.Info(r.Method, "-", r.URL.Path, r.RemoteAddr)

		w.Header().Set("Access-Control-Allow-Origin", "*")
		// compare the return-value to the authMW
		next.ServeHTTP(w, r)
	})
}

func (api *APIServer) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//Check the JWT in to cookies
		sessionCookie, err := r.Cookie("session")
		if err != nil {
			api.logger.Error("Error occured when getting the session cookie in the middleware", err.Error())
			//Return unauthenticated
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		token := sessionCookie.Value
		_, err = jwt.ValidateJWT(token)
		if err != nil {
			api.logger.Error("Error occured when validating JWT token in the middleware", err.Error())
			//Return unauthenticated
			w.WriteHeader(http.StatusUnauthorized)
			return
		}

		//Authenticated
		next.ServeHTTP(w, r)
	})
}

// Initialize the api http server based on the configuration file
func (api *APIServer) Init() error {
	//Initialize the logger
	api.logger = logging.NewDefaultDebugLogger()
	api.logger.Debug("Logger initialized")

	//Define command line arguments of the agent
	flag.StringVar(&api.configFile, "config", "", "The path to the configuration file")
	//Parse command line arguments
	flag.Parse()

	//Load the configuration from file
	err := api.configuration.LoadConfigurationFromFile(api.configFile)
	if err != nil {
		api.logger.Fatal("Error occured when loading the config from file,", err.Error())
		return err
	}
	api.logger.Debug("Loaded configuration from file")

	//Initialize the database connection
	api.dbConnection = database.NewCassandraConnection(api.logger, api.configuration)
	err = api.dbConnection.Init()
	if err != nil {
		api.logger.Error("Error occured when initializing database connection", err.Error())
		return err
	}
	api.logger.Debug("Connection to the database has been initialized")

	//Initialize the connection to elasticsearch
	api.elasticConnection = database.NewElasticConnection(api.logger, api.configuration)
	err = api.elasticConnection.Init()
	if err != nil {
		api.logger.Error("Error occured when initializing the elasticsearch connection", err.Error())
		return err
	}
	api.logger.Debug("Connection to elasticsearch has been initialized")

	//Create the admin account if it doesn't already exist
	res, err := api.dbConnection.CheckUserExists(api.configuration.AdminUsername)
	if err != nil {
		api.logger.Error("Error occured when checking if the admin user exists in the database", err.Error())
		return err
	}

	if !res {
		//Insert the user in the database
		_, err := api.dbConnection.InsertUser(api.configuration.AdminUsername, api.configuration.AdminPassword)
		if err != nil {
			api.logger.Error("Error occured when inserting the admin user in the database", err.Error())
			return err
		}
		api.logger.Info("Created admin user")
	}

	//Create the pool
	pool := websocket.NewPool(api.logger, api.dbConnection, api.configuration)
	//Start the pool in a goroutine
	go pool.Start()

	//Create the router
	r := mux.NewRouter()
	//Use the logging middleware
	r.Use(api.LoggingMiddleware)

	//Create the handlers
	healthCheckHandler := handlers.NewHealthCheckHandler(api.logger, api.configuration)
	authHandler := handlers.NewAuthHandler(api.logger, api.configuration, api.dbConnection, api.elasticConnection)
	//registerHandler := handlers.NewRegisterHandler(api.logger, api.configuration, api.dbConnection)
	agentsHandler := handlers.NewAgentsHandler(api.logger, api.configuration, api.dbConnection, api.elasticConnection)
	logsHandler := handlers.NewLogsHandler(api.logger, api.configuration, api.dbConnection, api.elasticConnection)
	machinesHandler := handlers.NewMachinesHandler(api.logger, api.configuration, api.dbConnection)
	wsHandler := handlers.NewWebsocketHandler(api.logger, api.configuration, api.dbConnection)

	//Create the standalone login route
	r.HandleFunc("/api/v1/auth/login", authHandler.Login)

	//Add the routes
	//Create the subrouter for the API path
	apiGetSubrouter := r.PathPrefix("/api/v1/").Methods("GET").Subrouter()
	apiPostSubrouter := r.PathPrefix("/api/v1/").Methods("POST").Subrouter()
	apiDeleteSubrouter := r.PathPrefix("/api/v1/").Methods("DELETE").Subrouter()
	apiPutSubrouter := r.PathPrefix("/api/v1/").Methods("PUT").Subrouter()

	//Add the auth middleware
	apiGetSubrouter.Use(api.AuthMiddleware)
	apiPutSubrouter.Use(api.AuthMiddleware)
	apiDeleteSubrouter.Use(api.AuthMiddleware)

	//Create the route that will send all the registered agents
	apiGetSubrouter.HandleFunc("/agents", agentsHandler.GetAgents)
	//Create the route that will send the number of registered agents
	apiGetSubrouter.HandleFunc("/agents/count", agentsHandler.GetAgentsCount)
	//Create the route that will send a single agent details
	apiGetSubrouter.HandleFunc("/agents/{uuid:[0-9a-f]+-[0-9a-f]+-[0-9a-f]+-[0-9a-f]+-[0-9a-f]+}", agentsHandler.GetAgent)
	//Create the route that will send the logs of an agent
	apiGetSubrouter.HandleFunc("/agents/{uuid:[0-9a-f]+-[0-9a-f]+-[0-9a-f]+-[0-9a-f]+-[0-9a-f]+}/logs", logsHandler.GetLogsShort)
	//Create the route that will send the logs of an agent from elasticsearch
	apiGetSubrouter.HandleFunc("/agents/{uuid:[0-9a-f]+-[0-9a-f]+-[0-9a-f]+-[0-9a-f]+-[0-9a-f]+}/logs-elastic", logsHandler.GetLogsShortElastic)
	//Create the route for exporting logs of an agent
	apiGetSubrouter.HandleFunc("/agents/{uuid:[0-9a-f]+-[0-9a-f]+-[0-9a-f]+-[0-9a-f]+-[0-9a-f]+}/export-logs", logsHandler.GetLogsShortElastic)
	//Create the route that will send the logs methods metrics
	apiGetSubrouter.HandleFunc("/agents/{uuid:[0-9a-f]+-[0-9a-f]+-[0-9a-f]+-[0-9a-f]+-[0-9a-f]+}/logs-methods-metrics", logsHandler.GetLogsMethodMetrics)
	//Create the route that will send the logs each day metrics
	apiGetSubrouter.HandleFunc("/agents/{uuid:[0-9a-f]+-[0-9a-f]+-[0-9a-f]+-[0-9a-f]+-[0-9a-f]+}/logs-each-day-metrics", logsHandler.GetLogsCountPerDay)
	//Create the route that will send the logs status codes metrics
	apiGetSubrouter.HandleFunc("/agents/{uuid:[0-9a-f]+-[0-9a-f]+-[0-9a-f]+-[0-9a-f]+-[0-9a-f]+}/logs-statuscode-metrics", logsHandler.GetResponseStatusCodesMetrics)
	//Create the route that will send the logs for IP addreses
	apiGetSubrouter.HandleFunc("/agents/{uuid:[0-9a-f]+-[0-9a-f]+-[0-9a-f]+-[0-9a-f]+-[0-9a-f]+}/logs-ipaddresses-metrics", logsHandler.GetIPAddressesMetrics)
	//Create the route that will send a specific log
	apiGetSubrouter.HandleFunc("/logs/{loguuid:[0-9a-f]+-[0-9a-f]+-[0-9a-f]+-[0-9a-f]+-[0-9a-f]+}", logsHandler.GetLog)
	//Create the route that will send the exploit code of a log
	apiGetSubrouter.HandleFunc("/logs/{loguuid:[0-9a-f]+-[0-9a-f]+-[0-9a-f]+-[0-9a-f]+-[0-9a-f]+}/exploit", logsHandler.GetLogExploitPythonCode)
	//Create the route that will send recent logs (10) to the client
	apiGetSubrouter.HandleFunc("/logs/recent", logsHandler.GetRecentLogsElastic)
	//Create the route that will send recent classified logs (10) to the client
	apiGetSubrouter.HandleFunc("/logs/recent-classified", logsHandler.GetRecentClassifiedLogsElastic)
	//Create the route that will send total count of logs
	apiGetSubrouter.HandleFunc("/logs/count", logsHandler.GetTotalLogsCount)
	//Create the route that will send all the logs
	apiGetSubrouter.HandleFunc("/logs", logsHandler.GetAllLogs)
	//Create the route that will send agent logs metrics
	apiGetSubrouter.HandleFunc("/logs/agent-metrics", logsHandler.GetAgentsMetrics)
	//Create the route that will send logs classification metrics
	apiGetSubrouter.HandleFunc("/logs/classification-metrics", logsHandler.GetClassificationMetrics)
	//Create the route that will send logs IP addresses metrics
	apiGetSubrouter.HandleFunc("/logs/ip-address-metrics", logsHandler.GetAllIPAddressesMetrics)
	//Create the route that will send all the classified logs
	apiGetSubrouter.HandleFunc("/logs/classified", logsHandler.GetAllClassifiedLogs)

	//Create the route that will send the findings count metrics
	apiGetSubrouter.HandleFunc("/findings/count-metrics", logsHandler.GetFindingsCount)
	//Create the route that will send the string format of all findings
	apiGetSubrouter.HandleFunc("/findings/string", logsHandler.GetFindingsClassificationString)
	//Create the route that will send the rule findings metrics
	apiGetSubrouter.HandleFunc("/findings/rule/metrics", logsHandler.GetLogsRuleFindingsMetrics)
	//Create the route that will send the rule id metrics
	apiGetSubrouter.HandleFunc("/findings/rule/id-metrics", logsHandler.GetLogsRuleIdMetrics)
	//Create the route that will send all the registered machines
	apiGetSubrouter.HandleFunc("/machines", machinesHandler.GetMachines)
	//Create the route that will send the machines statistics
	apiGetSubrouter.HandleFunc("/machines/metrics", machinesHandler.GetMachinesStatistics)

	//Create the route which will handle websocket dashboard connections
	r.HandleFunc("/ws", func(rw http.ResponseWriter, r *http.Request) {
		wsHandler.ServeDashboardWs(pool, rw, r)
	})

	//Create the route which will handle websocket agent connections
	r.HandleFunc("/agents/{uuid:[0-9a-f]+-[0-9a-f]+-[0-9a-f]+-[0-9a-f]+-[0-9a-f]+}/ws", func(rw http.ResponseWriter, r *http.Request) {
		wsHandler.ServeAgentWs(pool, rw, r)
	})

	//Create route to register agent
	apiPostSubrouter.HandleFunc("/registeragent", agentsHandler.RegisterAgent)
	//Create route to receive logs from agents
	apiPostSubrouter.HandleFunc("/addlog", agentsHandler.AddLog)
	//Create the route to register a new machine
	apiPostSubrouter.HandleFunc("/machines", machinesHandler.RegisterMachine)

	//Create the route to delete a machine
	apiDeleteSubrouter.HandleFunc("/machines/{machineuuid:[0-9a-f]+-[0-9a-f]+-[0-9a-f]+-[0-9a-f]+-[0-9a-f]+}", machinesHandler.DeleteMachine)

	//Create the route to update an agent
	apiPutSubrouter.HandleFunc("/agents/{uuid:[0-9a-f]+-[0-9a-f]+-[0-9a-f]+-[0-9a-f]+-[0-9a-f]+}", agentsHandler.ModifyAgent)

	//Create the healthcheck route
	r.HandleFunc("/api/v1/healthcheck", healthCheckHandler.HealthCheck)

	api.srv = &http.Server{
		Addr: api.configuration.ListeningAddress + ":" + api.configuration.ListeningPort,
		// Good practice to set timeouts to avoid Slowloris attacks.
		WriteTimeout: time.Second * 15,
		ReadTimeout:  time.Second * 15,
		IdleTimeout:  time.Second * 60,
		Handler:      r, // Pass our instance of gorilla/mux in.
	}

	return nil
}

// Start the api server
func (api *APIServer) Run() {
	var wait time.Duration = 5
	// Run our server in a goroutine so that it doesn't block.
	go func() {
		if err := api.srv.ListenAndServe(); err != nil {
			api.logger.Error(err.Error())
		}
	}()

	api.logger.Info("Started server on port", api.configuration.ListeningPort)

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
	api.srv.Shutdown(ctx)
	// Optionally, you could run srv.Shutdown in a goroutine and block on
	// <-ctx.Done() if your application should wait for other services
	// to finalize based on context cancellation.
	api.logger.Info("shutting down")
	os.Exit(0)
}
