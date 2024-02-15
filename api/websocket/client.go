package websocket

import (
	"log"

	"github.com/gorilla/websocket"
)

/*
 * This structure will define a client that connected to the chat service.
 * Each client will have a unique id, a websocket connection that will be used to send and receive messages
 * The client structure will also have a pointer to the pool structure which will be used for conccurency
 */
type DashboardClient struct {
	Id     int64
	Status string
	Conn   *websocket.Conn
	Pool   *Pool
}

type AgentClient struct {
	Id     string
	Status string
	Conn   *websocket.Conn
	Pool   *Pool
}

/*
 * This structure will define a message that can be sent/received on the websocket
 * The type variable will be used to determine if the websocket message is text or binary as it will have different values based on that
 * The body will be the payload that the other users it is delivered to should receive
 */
type DashboardMessage struct {
	C    *DashboardClient
	Type int    `json:"type"`
	Body string `json:"body"`
}

type AgentMessage struct {
	C    *AgentClient
	Type int    `json:"type"`
	Body string `json:"body"`
}

// // Constants
// const (
// 	// Time allowed to write a message to the peer.
// 	writeWait = 10 * time.Second

// 	// Time allowed to read the next pong message from the peer.
// 	pongWait = 60 * time.Second

// 	// Send pings to peer with this period. Must be less than pongWait.
// 	pingPeriod = (pongWait * 9) / 10

// 	// Maximum message size allowed from peer.
// 	maxMessageSize = 512
// )

/*
 * This function will wait for a message to be sent by the client and based on the message type different functions from the pool will be called
 */
func (c *DashboardClient) Read() {
	//Unregister a client when it disconnects from the server (this function will be called after the infinite loop)
	defer func() {
		c.Pool.UnregisterDashboard <- c
		c.Conn.Close()
	}()

	// c.Conn.SetReadLimit(maxMessageSize)
	// c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	// c.Conn.SetPongHandler(func(string) error {
	// 	c.Pool.logger.Debug("Received pong message from dashboard")
	// 	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	// 	return nil
	// })

	//Check if a message is received from the client
	for {
		//read the message from the client (here is where the message type is specified)
		messageType, p, err := c.Conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		//Create the message structure based on the message received from the client
		message := DashboardMessage{Type: messageType, Body: string(p), C: c}
		//Broadcast the message to all the other users (TO DO ... create private conversations and groups)
		c.Pool.DashboardBroadcast <- message
	}
}

// func (c *DashboardClient) Write() {
// 	ticker := time.NewTicker(pingPeriod)
// 	defer func() {
// 		ticker.Stop()
// 		c.Conn.Close()
// 	}()

// 	for ; ; <-ticker.C {
// 		c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
// 		c.Pool.logger.Debug("Sent pong message to dashboard")
// 		if err := c.Conn.WriteMessage(websocket.PingMessage, []byte("Test")); err != nil {
// 			c.Pool.logger.Debug(err.Error())
// 			return
// 		}
// 	}
// }

/*
 * This function will wait for a message to be sent by the client and based on the message type different functions from the pool will be called
 */
func (c *AgentClient) Read() {
	//Unregister a client when it disconnects from the server (this function will be called after the infinite loop)
	defer func() {
		c.Pool.UnregisterAgent <- c
		c.Conn.Close()
	}()
	// c.Conn.SetReadLimit(maxMessageSize)
	// c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	// c.Conn.SetPongHandler(func(string) error {
	// 	c.Pool.logger.Debug("Received pong message from agent")
	// 	c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	// 	return nil
	// })
	//Check if a message is received from the server
	for {
		//Read the message from the server
		messageType, p, err := c.Conn.ReadMessage()
		if err != nil {
			c.Pool.logger.Error("Error occured when reading message from agent", c.Id, "on the websocket,", err.Error())
			return
		}
		//Create the message structure based on the message received from the client
		message := AgentMessage{Type: messageType, Body: string(p), C: c}
		//Send the message to the WS message handler
		c.Pool.AgentBroadcast <- message
	}
}

// func (c *AgentClient) Write() {
// 	ticker := time.NewTicker(pingPeriod)
// 	defer func() {
// 		ticker.Stop()
// 		c.Conn.Close()
// 	}()

// 	for ; ; <-ticker.C {
// 		c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
// 		c.Pool.logger.Debug("Sent pong message to agent")
// 		if err := c.Conn.WriteMessage(websocket.PingMessage, []byte("Test")); err != nil {
// 			c.Pool.logger.Debug(err.Error())
// 			return
// 		}
// 	}
// }
