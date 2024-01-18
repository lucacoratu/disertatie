/*
// This structure holds the log data that will be sent to the client (short version)
type LogDataShort struct {
	Id              string            `json:"id"`               //The UUID of the log from the database
	AgentId         string            `json:"agentId"`          //The UUID of the agent that collected the log data
	RemoteIP        string            `json:"remoteIp"`         //The IP address of the sender of the request
	Timestamp       int64             `json:"timestamp"`        //Timestamp when the request was received
	RequestPreview  string            `json:"request_preview"`  //The preview of the request
	ResponsePreview string            `json:"response_preview"` //The preview of the response
	Findings        []FindingDatabase `json:"findings"`         //The findings extracted from the database
}
 */

 // Findings found by the agent, one for request, one for response
/*
 type Finding struct {
	Request  FindingData `json:"request"`  //The finding for the request
	Response FindingData `json:"response"` //The finding for the response
}
*/

/*
// Structure that holds the information about the finding
type FindingData struct {
	Line           int64  `json:"line"`           //The line from the request where the finding is located
	LineIndex      int64  `json:"lineIndex"`      //The offset from the start of the line
	Length         int64  `json:"length"`         //The length of the finding string
	Classification int64  `json:"classification"` //The classification of the finding based on the constants above
	Severity       int64  `json:"severity"`       //The severity of the finding
	ValidatorName  string `json:"validatorName"`  //The name of the validator who made the discovery
}
 */

/*
// Severity types
const (
	LOW      int64 = 0
	MEDIUM   int64 = 1
	HIGH     int64 = 2
	CRITICAL int64 = 3
)
 */

enum severityTypes {
    LOW = 0,
    MEDIUM = 1,
    HIGH = 2,
    CRITICAL = 3,
}

/*
// Types of finding classifications
const (
	//Unknown classification
	UNKNOWN int64 = -1

	//Request classifications
	LFI_ATTACK int64 = 0
    SCRIPT_USER_AGENT int64 = 1

	//Response classifications
	UNAUTHORIZED_ACCESS int64 = 100
	FILE_OUT            int64 = 101
	FLAG_OUT            int64 = 102
)
*/

enum classification {
    UNKNOWN = -1,
    
    //Request classifications
    LFI_ATTACK = 0,
    SCRIPT_USER_AGENT = 1,

    //Response classification
	UNAUTHORIZED_ACCESS = 100,
	FILE_OUT = 101,
	FLAG_OUT = 102,
}

type FindingData = {
    id: string,
    log_id: string,
    line: number,
    lineIndex: number,
    length: number,
    classification: number,
    severity: number,
    validatorName: string,
}

type Finding = {
    request: FindingData,
    response: FindingData,
}

type LogShort = {
    id: string,
    agentId: string,
    remoteIp: string,
    timestamp: number,
    request_preview: string,
    response_preview: string,
    findings: Finding[],
}

type LogShortResponse = {
    logs: LogShort[]
}

type LogShortProps = {
    logs: LogShort[]
}

/*
// This structure holds the log data that is in the database
type LogDataDatabase struct {
	Id              string    `json:"id"`               //The UUID of the log from the database
	AgentId         string    `json:"agentId"`          //The UUID of the agent that collected the log data
	RemoteIP        string    `json:"remoteIp"`         //The IP address of the sender of the request
	Timestamp       int64     `json:"timestamp"`        //Timestamp when the request was received
	RequestPreview  string    `json:"request_preview"`  //The preview of the request
	ResponsePreview string    `json:"response_preview"` //The preview of the response
	Request         string    `json:"request"`          //The request base64 encoded
	Response        string    `json:"response"`         // The response base64 encoded
	Findings        []Finding `json:"findings"`         //A list of findings
}
*/

type LogFull = {
	id: string,
	agentId: string,
	remoteIp: string,
	timestamp: number,
	request_preview: string,
	response_preview: string,
	request: string,
	response: string,
	findings: Finding[], 
}

type LogFullResponse = {
	log: LogFull
}