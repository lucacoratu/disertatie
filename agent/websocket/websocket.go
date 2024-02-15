package websocket

import (
	"github.com/gorilla/websocket"
	"github.com/lucacoratu/disertatie/agent/config"
	"github.com/lucacoratu/disertatie/agent/data"
	"github.com/lucacoratu/disertatie/agent/logging"
)

type message struct {
	Type int    `json:"type"`
	Body string `json:"body"`
}

type Notification struct {
	AgentId string `json:"agentId"`
	Message string `json:"message"`
}

type APIWebSocketConnection struct {
	logger        logging.ILogger      //The logger
	apiWsURL      string               //The ws url of the API
	configuration config.Configuration //The configuration structure of the agent
	State         bool                 //The state of the websocket connection (true for active, false for inactive)
	connection    *websocket.Conn      //The connection structure
}

func NewAPIWebSocketConnection(logger logging.ILogger, apiWsURL string, configuration config.Configuration) *APIWebSocketConnection {
	return &APIWebSocketConnection{logger: logger, apiWsURL: apiWsURL, configuration: configuration}
}

// Connects to the API websocket URL for the agent
func (awsc *APIWebSocketConnection) Connect() (bool, error) {
	//Connect to the websocket URL from the API
	c, _, err := websocket.DefaultDialer.Dial(awsc.apiWsURL, nil)

	//Check if an error occured
	if err != nil {
		return false, err
	}

	awsc.connection = c
	awsc.State = true
	return true, nil
}

// // Handle the connection closed
// func (awsc *APIWebSocketConnection) connectionClosed(code int, text string) error {
// 	return nil
// }

// Handle the message received
func (awsc *APIWebSocketConnection) handleReceivedMessage(message message) {
	awsc.logger.Debug("Message received", message)
}

func (awsc *APIWebSocketConnection) Start() {
	//Close the connection at the end of the function
	defer awsc.connection.Close()

	//Start listening for incomming messages
	for {
		mt, msg, err := awsc.connection.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				awsc.logger.Error(err.Error())
				awsc.Connect()
				continue
			}
			//If it is other type of error exit
			awsc.logger.Debug("Websocket connection to the API closed")
			return
		}

		//Call the handle message function
		awsc.handleReceivedMessage(message{Type: mt, Body: string(msg)})
	}
}

// Function to send a notification to the API
func (awsc *APIWebSocketConnection) SendNotification(message string) error {
	notif := Notification{AgentId: awsc.configuration.UUID, Message: message}
	return awsc.connection.WriteJSON(data.WebSocketMessage{Type: data.Notification, Data: notif})
}
