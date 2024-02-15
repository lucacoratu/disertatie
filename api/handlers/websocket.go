package handlers

import (
	"net/http"

	"github.com/gorilla/mux"
	"github.com/lucacoratu/disertatie/api/config"
	"github.com/lucacoratu/disertatie/api/database"
	"github.com/lucacoratu/disertatie/api/logging"
	"github.com/lucacoratu/disertatie/api/websocket"
)

type WebsocketHandler struct {
	logger        logging.ILogger
	configuration config.Configuration
	dbConnection  database.IConnection
}

func NewWebsocketHandler(logger logging.ILogger, configuration config.Configuration, dbConnection database.IConnection) *WebsocketHandler {
	return &WebsocketHandler{logger: logger, configuration: configuration, dbConnection: dbConnection}
}

/*
 * This function will handle when a client connects to the websocket endpoint
 */
func (wsh *WebsocketHandler) ServeDashboardWs(pool *websocket.Pool, rw http.ResponseWriter, r *http.Request) {
	//Upgrade the connection to a Websocket connection
	ws, err := websocket.Upgrade(rw, r)
	//Check if an error occured
	if err != nil {
		//Log the error
		wsh.logger.Error(err.Error())
		return
	}

	//Create the client structure which will be saved in the pool
	client := &websocket.DashboardClient{
		Conn:   ws,
		Pool:   pool,
		Status: "Offline",
		Id:     0,
	}

	//Call the client register function
	pool.RegisterDashboard <- client
	//Start reading data from the connection
	//go client.Write()
	go client.Read()
}

/*
 * This function will handle when a client connects to the websocket endpoint
 */
func (wsh *WebsocketHandler) ServeAgentWs(pool *websocket.Pool, rw http.ResponseWriter, r *http.Request) {
	//Get the agent UUID from the mux variables
	vars := mux.Vars(r)
	agent_uuid := vars["uuid"]

	//Upgrade the connection to a Websocket connection
	ws, err := websocket.Upgrade(rw, r)
	//Check if an error occured
	if err != nil {
		//Log the error
		wsh.logger.Error(err.Error())
		return
	}

	//Create the client structure which will be saved in the pool
	client := &websocket.AgentClient{
		Conn:   ws,
		Pool:   pool,
		Status: "Offline",
		Id:     agent_uuid,
	}

	//Call the client register function
	pool.RegisterAgent <- client
	//Start reading data from the connection
	//go client.Write()
	go client.Read()
}
