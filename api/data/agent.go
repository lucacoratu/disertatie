package data

import (
	"encoding/json"
	"io"
)

type Agent struct {
	ID                    string   `json:"id"`                    //The UUID of the agent
	Name                  string   `json:"name"`                  //The name of the agent
	ListeningProtocol     string   `json:"listeningProtocol"`     //The listening protocol of the agent
	ListeningAddress      string   `json:"listeningAddress"`      //The listening address of the agent
	ListeningPort         int64    `json:"listeningPort"`         //The port the agent is listening on
	ForwardServerProtocol string   `json:"forwardServerProtocol"` //The protocol of the webserver that the agent sends requests to
	ForwardServerAddress  string   `json:"forwardServerAddress"`  //The address of the webserver that the agent sends requests to
	ForwardServerPort     string   `json:"forwardServerPort"`     //The port of the webserver that the agent sends requests to
	MachineId             string   `json:"machineId"`             //The id of the machine the agent is deployed on
	MachineOS             string   `json:"machineOs"`             //The OS of the machine the agent is deployed on
	MachineHostname       string   `json:"machineHostname"`       //The hostname of the machine the agent is deployed on
	MachineIPAddreses     []string `json:"machineIpAddreses"`     //A list with all the ip addreses of the machine the agent is deployed on
	LogsCollected         int64    `json:"logsCollected"`         //The number of logs collected by the agent
}

// Convert from json into the agent structure
func (ag *Agent) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(ag)
}

// Convert to json the agent structure
func (ag *Agent) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(ag)
}

type UpdateAgent struct {
	Name                  string `json:"name"`                  //The name of the agent
	ListeningProtocol     string `json:"listeningProtocol"`     //The listening protocol of the agent
	ListeningAddress      string `json:"listeningAddress"`      //The listening address of the agent
	ListeningPort         int64  `json:"listeningPort"`         //The port the agent is listening on
	ForwardServerProtocol string `json:"forwardServerProtocol"` //The protocol of the webserver that the agent sends requests to
	ForwardServerAddress  string `json:"forwardServerAddress"`  //The address of the webserver that the agent sends requests to
	ForwardServerPort     int64  `json:"forwardServerPort"`     //The port of the webserver that the agent sends requests to
	MachineId             string `json:"machineId"`             //The id of the machine the agent is deployed on
}

func (ua *UpdateAgent) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(ua)
}

func (ua *UpdateAgent) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(ua)
}
