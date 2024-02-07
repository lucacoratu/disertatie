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
	"github.com/lucacoratu/disertatie/api/logging"
)

type APIServer struct {
	srv           *http.Server
	logger        logging.ILogger
	configuration config.Configuration
	dbConnection  database.IConnection
	configFile    string
}

func (api *APIServer) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		api.logger.Info(r.Method, "-", r.URL.Path, r.RemoteAddr)

		w.Header().Set("Access-Control-Allow-Origin", "*")
		// compare the return-value to the authMW
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

	//Create the router
	r := mux.NewRouter()
	//Use the logging middleware
	r.Use(api.LoggingMiddleware)

	//Create the handlers
	healthCheckHandler := handlers.NewHealthCheckHandler(api.logger, api.configuration)
	//registerHandler := handlers.NewRegisterHandler(api.logger, api.configuration, api.dbConnection)
	agentsHandler := handlers.NewAgentsHandler(api.logger, api.configuration, api.dbConnection)
	logsHandler := handlers.NewLogsHandler(api.logger, api.configuration, api.dbConnection)
	machinesHandler := handlers.NewMachinesHandler(api.logger, api.configuration, api.dbConnection)

	//Add the routes
	//Create the subrouter for the API path
	apiGetSubrouter := r.PathPrefix("/api/v1/").Methods("GET").Subrouter()
	apiPostSubrouter := r.PathPrefix("/api/v1/").Methods("POST").Subrouter()
	apiDeleteSubrouter := r.PathPrefix("/api/v1/").Methods("DELETE").Subrouter()

	//Create the route that will send all the registered agents
	apiGetSubrouter.HandleFunc("/agents", agentsHandler.GetAgents)
	//Create the route that will send a single agent details
	apiGetSubrouter.HandleFunc("/agents/{uuid:[0-9a-f]+-[0-9a-f]+-[0-9a-f]+-[0-9a-f]+-[0-9a-f]+}", agentsHandler.GetAgent)
	//Create the route that will send the logs of an agent
	apiGetSubrouter.HandleFunc("/agents/{uuid:[0-9a-f]+-[0-9a-f]+-[0-9a-f]+-[0-9a-f]+-[0-9a-f]+}/logs", logsHandler.GetLogsShort)
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
	//Create the route that will send the string format of all findings
	apiGetSubrouter.HandleFunc("/findings/string", logsHandler.GetFindingsClassificationString)
	//Create the route that will send all the registered machines
	apiGetSubrouter.HandleFunc("/machines", machinesHandler.GetMachines)

	//Create route to register agent
	apiPostSubrouter.HandleFunc("/registeragent", agentsHandler.RegisterAgent)
	//Create route to receive logs from agents
	apiPostSubrouter.HandleFunc("/addlog", agentsHandler.AddLog)
	//Create the route to register a new machine
	apiPostSubrouter.HandleFunc("/machines", machinesHandler.RegisterMachine)

	//Create the route to delete a machine
	apiDeleteSubrouter.HandleFunc("/machines/{machineuuid:[0-9a-f]+-[0-9a-f]+-[0-9a-f]+-[0-9a-f]+-[0-9a-f]+}", machinesHandler.DeleteMachine)

	//Create the healthcheck route
	apiGetSubrouter.HandleFunc("/healthcheck", healthCheckHandler.HealthCheck)

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
