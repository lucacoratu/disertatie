package data

import (
	"encoding/json"
	"io"
)

// Structure that holds information about the metrics of methods found by the agent
type MethodsMetrics struct {
	Id     int64  `json:"id"`     //The ID of the method metrics
	Method string `json:"method"` //The method name
	Count  int64  `json:"count"`  //The number of occurences of specific method
}

// Convert json data to MethodsMetrics structure
func (mm *MethodsMetrics) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(mm)
}

// Convert MethodsMetrics structure to json string
func (mm *MethodsMetrics) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(mm)
}

// Structure that holds information about the metrics of requests received by agent based on timestamp
type DayMetrics struct {
	Id    int64  `json:"id"`    //The ID of the metric
	Date  string `json:"date"`  //The date when the counts have been computed
	Count int64  `json:"count"` //The number of requests till the specified date
}

// Convert json data to MethodsMetrics structure
func (dm *DayMetrics) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(dm)
}

// Convert MethodsMetrics structure to json string
func (dm *DayMetrics) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(dm)
}

// Structure that holds information about the metrics of response status codes received by agent
type StatusCodesMetrics struct {
	Id         int64  `json:"id"`         //The ID of the metric
	StatusCode string `json:"statusCode"` //The status code of the metric
	Count      int64  `json:"count"`      //The number of occurences of the status code
}

// Convert json data to StatusCodesMetrics structure
func (scm *StatusCodesMetrics) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(scm)
}

// Convert StatusCodesMetrics structure to json string
func (scm *StatusCodesMetrics) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(scm)
}

// Structure that holds information about the metrics of remote IP addreses received by agent
type IPMetrics struct {
	Id        int64  `json:"id"`        //The ID of the metric
	IPAddress string `json:"ipAddress"` //The IP address of the metric
	Count     int64  `json:"count"`     //The number of occurences of the IP address
}

// Convert json data to StatusCodesMetrics structure
func (ipm *IPMetrics) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(ipm)
}

// Convert StatusCodesMetrics structure to json string
func (ipm *IPMetrics) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(ipm)
}

// Structure that holds metrics about the findings (the classification string and the number of appearences)
type FindingsMetrics struct {
	Classification string `json:"classification"` //The classification string
	Count          int64  `json:"count"`          //The count of occurences
}

// Convert json data to StatusCodesMetrics structure
func (fm *FindingsMetrics) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(fm)
}

// Convert StatusCodesMetrics structure to json string
func (fm *FindingsMetrics) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(fm)
}

// Structure that holds the metrics about counts of findings
type FindingsCountMetrics struct {
	FindingsCount     int64 `json:"findingsCount"`
	RuleFindingsCount int64 `json:"ruleFindingsCount"`
}

func (fcm *FindingsCountMetrics) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(fcm)
}

func (fcm *FindingsCountMetrics) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(fcm)
}

// Structure that holds the metrics about number of logs collected from each agent
type AgentsMetrics struct {
	AgentId   string `json:"agentId"`
	AgentName string `json:"agentName"`
	Count     int64  `json:"count"`
}

func (am *AgentsMetrics) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(am)
}

func (am *AgentsMetrics) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(am)
}

type ClassificationMetrics struct {
	ClassifiedCount   int64 `json:"classifiedCount"`
	UnclassifiedCount int64 `json:"unclassifiedCount"`
}

func (cm *ClassificationMetrics) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(cm)
}

func (cm *ClassificationMetrics) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(cm)
}
