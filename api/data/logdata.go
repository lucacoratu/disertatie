package data

import (
	"encoding/json"
	"io"
)

// This structure holds the log data that is sent to the api
type LogData struct {
	Id           string        `json:"id"`           //The UUID of the log from the database
	AgentId      string        `json:"agentId"`      //The UUID of the agent that collected the log data
	RemoteIP     string        `json:"remoteIp"`     //The IP address of the sender of the request
	Timestamp    int64         `json:"timestamp"`    //Timestamp when the request was received
	Websocket    bool          `json:"websocket"`    //If the message is a websocket
	Request      string        `json:"request"`      //The request base64 encoded
	Response     string        `json:"response"`     // The response base64 encoded
	Findings     []Finding     `json:"findings"`     //A list of findings
	RuleFindings []RuleFinding `json:"ruleFindings"` //The list of rule findings
}

// Convert json data to LogData structure
func (ld *LogData) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(ld)
}

// Convert LogData structure to json string
func (ld *LogData) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(ld)
}

// This structure holds the log data that is in the database
type LogDataDatabase struct {
	Id              string                `json:"id"`               //The UUID of the log from the database
	AgentId         string                `json:"agentId"`          //The UUID of the agent that collected the log data
	RemoteIP        string                `json:"remoteIp"`         //The IP address of the sender of the request
	Timestamp       int64                 `json:"timestamp"`        //Timestamp when the request was received
	RequestPreview  string                `json:"request_preview"`  //The preview of the request
	ResponsePreview string                `json:"response_preview"` //The preview of the response
	Request         string                `json:"request"`          //The request base64 encoded
	Response        string                `json:"response"`         // The response base64 encoded
	Findings        []FindingDatabase     `json:"findings"`         //A list of findings
	RuleFindings    []RuleFindingDatabase `json:"ruleFindings"`     //The list of rule findings
}

// This structure holds the log data that will be sent to the client (short version)
type LogDataShort struct {
	Id              string                `json:"id"`               //The UUID of the log from the database
	AgentId         string                `json:"agentId"`          //The UUID of the agent that collected the log data
	RemoteIP        string                `json:"remoteIp"`         //The IP address of the sender of the request
	Timestamp       int64                 `json:"timestamp"`        //Timestamp when the request was received
	RequestPreview  string                `json:"request_preview"`  //The preview of the request
	ResponsePreview string                `json:"response_preview"` //The preview of the response
	Findings        []FindingDatabase     `json:"findings"`         //The findings extracted from the database
	RuleFindings    []RuleFindingDatabase `json:"ruleFindings"`     //The list of rule findings
}

// Convert json data to LogData structure
func (lds *LogDataShort) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(lds)
}

// Convert LogData structure to json string
func (lds *LogDataShort) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(lds)
}

// This structure holds the log data that is sent to the api
type LogDataElastic struct {
	Id              string        `json:"id"`               //The UUID of the log from the database
	AgentId         string        `json:"agentId"`          //The UUID of the agent that collected the log data
	AgentName       string        `json:"agentName"`        //The name of the agent that collected the log data
	RemoteIP        string        `json:"remoteIp"`         //The IP address of the sender of the request
	Timestamp       int64         `json:"timestamp"`        //Timestamp when the request was received
	RequestPreview  string        `json:"request_preview"`  //The preview of the request
	ResponsePreview string        `json:"response_preview"` //The preview of the response
	Findings        []Finding     `json:"findings"`         //A list of findings
	Websocket       bool          `json:"websocket"`        //If the message is an websocket message
	RuleFindings    []RuleFinding `json:"ruleFindings"`     //The list of rule findings
}

// Convert json data to LogData structure
func (ld *LogDataElastic) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(ld)
}

// Convert LogData structure to json string
func (ld *LogDataElastic) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(ld)
}
