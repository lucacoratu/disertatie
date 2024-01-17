package data

import (
	"encoding/json"
	"io"
)

// Types of finding classifications
const (
	//Unknown classification
	UNKNOWN int64 = -1

	//Request classifications
	LFI_ATTACK        int64 = 0
	SCRIPT_USER_AGENT int64 = 1

	//Response classifications
	UNAUTHORIZED_ACCESS int64 = 100
	FILE_OUT            int64 = 101
	FLAG_OUT            int64 = 102
)

var ClassificationsMap = map[int64]string{
	LFI_ATTACK:        "LFI",
	SCRIPT_USER_AGENT: "Script UA",
}

var ClassificationDescriptionMap = map[int64]string{
	LFI_ATTACK:        "Local File Inclusion Attack",
	SCRIPT_USER_AGENT: "User Agent used by scripts/tools to automatically enumerate websites",
}

type FindingClassificationString struct {
	IntegerFormat int64  `json:"intFormat"`    //The integer format
	StringFormat  string `json:"stringFormat"` //The string format
	Description   string `json:"description"`  //The description of the classification
}

// Severity types
const (
	LOW      int64 = 0
	MEDIUM   int64 = 1
	HIGH     int64 = 2
	CRITICAL int64 = 3
)

// Structure that holds the information about the finding
type FindingData struct {
	Line           int64  `json:"line"`           //The line from the request where the finding is located
	LineIndex      int64  `json:"lineIndex"`      //The offset from the start of the line
	Length         int64  `json:"length"`         //The length of the finding string
	Classification int64  `json:"classification"` //The classification of the finding based on the constants above
	Severity       int64  `json:"severity"`       //The severity of the finding
	ValidatorName  string `json:"validatorName"`  //The name of the validator who made the discovery
}

// Structure that holds the information stored in the database about findings
type FindingDataDatabase struct {
	Id             string `json:"id"`             //The ID of the finding
	LogId          string `json:"log_id"`         //The log ID
	Line           int64  `json:"line"`           //The line from the request where the finding is located
	LineIndex      int64  `json:"lineIndex"`      //The offset from the start of the line
	Length         int64  `json:"length"`         //The length of the finding string
	Classification int64  `json:"classification"` //The classification of the finding based on the constants above
	Severity       int64  `json:"severity"`       //The severity of the finding
	ValidatorName  string `json:"validatorName"`  //The name of the validator who made the discovery
}

// Findings found by the agent, one for request, one for response
type Finding struct {
	Request  FindingData `json:"request"`  //The finding for the request
	Response FindingData `json:"response"` //The finding for the response
}

// Findings extracted from the database
type FindingDatabase struct {
	Request  FindingDataDatabase `json:"request"`  //The findings for the request
	Response FindingDataDatabase `json:"response"` //The findings for the response
}

// Convert finding to JSON
func (f *Finding) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(f)
}

// Convert finding from JSON
func (f *Finding) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(f)
}

// Convert finding to JSON
func (fd *FindingDatabase) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(fd)
}

// Convert finding from JSON
func (fd *FindingDatabase) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(fd)
}
