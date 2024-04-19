package websocket

import (
	"encoding/json"
	"strings"

	"github.com/lucacoratu/disertatie/api/config"
	"github.com/lucacoratu/disertatie/api/data"
	"github.com/lucacoratu/disertatie/api/database"
	"github.com/lucacoratu/disertatie/api/logging"
)

/*
 * This structure will handle concurrent connections using channels
 * Each channel will have a particular functionality
 */
type Pool struct {
	RegisterDashboard   chan *DashboardClient     //Channel which will handle new dashboard client connections
	UnregisterDashboard chan *DashboardClient     //Channel which will handle a dashboard client disconnecting
	DashboardClients    map[*DashboardClient]bool //A map of dashboard client connections and associated state of the connection (true for online)
	DashboardBroadcast  chan DashboardMessage     //Channel which will be used to handle a message from the dashboard client
	RegisterAgent       chan *AgentClient         //Channel which will handle new agent connections
	UnregisterAgent     chan *AgentClient         //Channgel which will handle agent client disconnecting
	AgentClients        map[*AgentClient]bool     //A map of dashboard client connections and associated state of the connection (true for online)
	AgentBroadcast      chan AgentMessage         //Channel which will be used to handle a message from the agent
	logger              logging.ILogger           //The logger
	dbConnection        database.IConnection      //The database connection
	configuration       config.Configuration      //The configuration
}

/*
 * This function will create a new pool that can then be used when starting the chat service
 */
func NewPool(l logging.ILogger, dbConn database.IConnection, conf config.Configuration) *Pool {
	return &Pool{
		RegisterDashboard:   make(chan *DashboardClient),
		UnregisterDashboard: make(chan *DashboardClient),
		DashboardClients:    make(map[*DashboardClient]bool),
		DashboardBroadcast:  make(chan DashboardMessage),
		RegisterAgent:       make(chan *AgentClient),
		UnregisterAgent:     make(chan *AgentClient),
		AgentClients:        make(map[*AgentClient]bool),
		AgentBroadcast:      make(chan AgentMessage),
		logger:              l,
		dbConnection:        dbConn,
		configuration:       conf,
	}
}

/*
 * This function will handle the register of a client (the client connected to the websocket)
 * It should be registered in the pool
 */
func (pool *Pool) DashboardClientRegistered(c *DashboardClient) {
	pool.logger.Info("Dashboard client connected to websocket")
}

func (pool *Pool) DashboardClientUnregistered(c *DashboardClient) {
	pool.logger.Info("Dashboard client disconnected from websocket")
}

/*
 * This function will handle when a message is recevied from a client
 * There should be more types of messages that can be received from the client
 */
func (pool *Pool) DashboardMessageReceived(message DashboardMessage) {
	//Log that a message has been received on the websocket
	pool.logger.Info("Dashboard message received on the websocket")
}

func (pool *Pool) AgentRegistered(c *AgentClient) {
	c.Status = "online"
	pool.logger.Info("Agent connected to websocket, id:", c.Id)
	//Send a notification to all the dashboard client to announce that the agent connected to the API
	message := Notification{AgentId: c.Id, Message: "Agent connected"}
	wsMessage := WebSocketMessage{Type: WsAgentConnectedNotification, Data: message}
	for client := range pool.AgentClients {
		err := client.Conn.WriteJSON(wsMessage)
		//Check if an error occured when sending the notification to the dashboard client
		if err != nil {
			pool.logger.Error("Error occured when sending agent connect notification to dashboard client, id:", client.Id)
		}
	}
}

func (pool *Pool) AgentUnregistered(c *AgentClient) {
	c.Status = "offline"
	pool.logger.Info("Agent disconnected from websocket, id: ", c.Id)
	//Send a disconnect message to all the dashboard clients
	message := Notification{AgentId: c.Id, Message: "Agent disconnected"}
	wsMessage := WebSocketMessage{Type: WsAgentDisconnectedNotification, Data: message}
	for client := range pool.AgentClients {
		err := client.Conn.WriteJSON(wsMessage)
		//Check if an error occured when sending the notification to the dashboard client
		if err != nil {
			pool.logger.Error("Error occured when sending agent disconnect notification to dashboard client, id:", client.Id)
		}
	}
}

/*
 * This function will handle when a message is recevied from a client
 * There should be more types of messages that can be received from the client
 */
func (pool *Pool) AgentMessageReceived(message AgentMessage) {
	//Log that a message has been received on the websocket
	pool.logger.Info("Agent message received on the websocket", message.Body)
	//Parse the message body to a websocket message
	wsMessage := WebSocketMessage{}
	err := wsMessage.FromJSON(strings.NewReader(message.Body))
	//Check if an error occured when parsing the WebSocketMessage from JSON
	if err != nil {
		//Send an error message back to the client
		errMessage := WebSocketMessage{Type: WsError, Data: data.APIError{Code: data.PARSE_ERROR, Message: "Cannot parse the websocket message from JSON"}}
		message.C.Conn.WriteJSON(errMessage)
		return
	}

	//Select the action based on the message type
	switch wsMessage.Type {
	case WsError:
		pool.logger.Debug("Error message received")
	case WsNotification:
		//Handle the notification
		err = pool.HandleNotification(wsMessage)
		//Check if an error occured when handling the notification
		if err != nil {
			//Send an error message back to the client
			errMessage := WebSocketMessage{Type: WsError, Data: data.APIError{Code: data.WS_ERROR, Message: err.Error()}}
			message.C.Conn.WriteJSON(errMessage)
			return
		}
		//pool.logger.Debug("Notification received")
	case WsAgentStatusRequest:
		response, err := pool.HandleAgentStatusRequest(wsMessage)
		//Check if an error occured
		if err != nil {
			//Send an error message back to the client
			errMessage := WebSocketMessage{Type: WsError, Data: data.APIError{Code: data.WS_ERROR, Message: err.Error()}}
			message.C.Conn.WriteJSON(errMessage)
			return
		}
		//Send the response back to the client
		message.C.Conn.WriteJSON(response)
	case WsRuleDetectionAlert:
		err = pool.HandleRuleDetectionAlert(wsMessage)
		if err != nil {
			//Send an error message back to the client
			errMessage := WebSocketMessage{Type: WsError, Data: data.APIError{Code: data.WS_ERROR, Message: err.Error()}}
			message.C.Conn.WriteJSON(errMessage)
			return
		}
	}
}

/*
 * This function will start the pool which will handle client connections, client disconnections and broadcast messages
 */
func (pool *Pool) Start() {
	//Loop infinetly
	for {
		//Check what kind of event occured (connect, disconnect, broadcast message)
		select {
		case client := <-pool.RegisterDashboard:
			//A new client connected on the chat service websocket uri
			//Add the client connection to the pool of current connections
			pool.DashboardClients[client] = true
			pool.logger.Debug("Size of dashboard connection pool", len(pool.DashboardClients))
			pool.DashboardClientRegistered(client)

		case client := <-pool.UnregisterDashboard:
			//A client disconnected from the chat service
			pool.DashboardClientUnregistered(client)
			delete(pool.DashboardClients, client)
			pool.logger.Debug("Size of dashboard connection pool: ", len(pool.DashboardClients))

		case message := <-pool.DashboardBroadcast:
			//Message received on the websocket
			pool.DashboardMessageReceived(message)

		case client := <-pool.RegisterAgent:
			//Agent connected to the websocket
			pool.AgentClients[client] = true
			pool.logger.Debug("Size of agents connection pool", len(pool.AgentClients))
			pool.AgentRegistered(client)

		case client := <-pool.UnregisterAgent:
			//Agent client disconnected from the websocket
			pool.AgentUnregistered(client)
			delete(pool.AgentClients, client)
			pool.logger.Debug("Size of agents connection pool: ", len(pool.AgentClients))

		case message := <-pool.AgentBroadcast:
			//Message received from the agent on the websocket
			pool.AgentMessageReceived(message)
		}
	}
}

func (pool *Pool) HandleNotification(msg WebSocketMessage) error {
	// data, _ := json.Marshal(msg.Data)
	// notif := Notification{}
	// _ = json.Unmarshal(data, &notif)
	// pool.logger.Debug(notif)
	//Send the notification to all the dashboard clients
	for client := range pool.DashboardClients {
		client.Conn.WriteJSON(msg)
	}
	return nil
}

func (pool *Pool) HandleAgentStatusRequest(msg WebSocketMessage) (AgentStatusResponse, error) {
	//Convert the message data field to the corresponding structure
	data, _ := json.Marshal(msg.Data)
	agentStatusRequest := AgentStatusRequest{}
	err := json.Unmarshal(data, &agentStatusRequest)
	if err != nil {
		return AgentStatusResponse{}, err
	}
	response := AgentStatusResponse{}
	//Get the status of the agent from the pool of connections
	var found bool = false
	for agent := range pool.AgentClients {
		if agent.Id == agentStatusRequest.AgentId {
			response.AgentId = agent.Id
			response.Status = agent.Status
			found = true
			break
		}
	}
	//Check if he agent has been found
	if !found {
		//If the agent was not found then it is offline
		response.AgentId = agentStatusRequest.AgentId
		response.Status = "offline"
	}
	//Return the response
	return response, nil
}

func (pool *Pool) HandleRuleDetectionAlert(msg WebSocketMessage) error {
	//Send the alert to all the dashboard clients
	for client := range pool.DashboardClients {
		client.Conn.WriteJSON(msg)
	}
	return nil
}
