/*
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
*/

type Agent = {
    id:                    string,
    name:                  string,
    listeningProtocol:     string,
    listeningAddress:      string,
    listeningPort:         string,
    forwardServerProtocol: string,
	forwardServerAddress:  string,   
	forwardServerPort:     string,
	machineId:             string,
	machineOs:             string,
	machineHostname:       string,
	machineIPAddreses:     string[],
} 

type AgentResponse = {
    agents: Agent[]
}

interface AgentProps {
    agent: Agent
}

type UpdateAgent = {
	name:                  string,
    listeningProtocol:     string,
    listeningAddress:      string,
    listeningPort:         number,
    forwardServerProtocol: string,
	forwardServerAddress:  string,   
	forwardServerPort:     number,
	machineId: 			   string,
}