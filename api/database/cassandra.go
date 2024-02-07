package database

import (
	"errors"
	"strings"

	b64 "encoding/base64"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"github.com/lucacoratu/disertatie/api/config"
	"github.com/lucacoratu/disertatie/api/data"
	"github.com/lucacoratu/disertatie/api/logging"
)

// Structure that will hold the necessary fields and functions for interacting with Cassandra
type CassandraConnection struct {
	logger        logging.ILogger
	configuration config.Configuration
	session       *gocql.Session
}

// Creates a new instance of the CassandraConnection structure
func NewCassandraConnection(logger logging.ILogger, configuration config.Configuration) *CassandraConnection {
	return &CassandraConnection{logger: logger, configuration: configuration}
}

// Creates the tables needed in the cassandra database (keyspace)
func (cassandra *CassandraConnection) createTables() error {
	// _ = cassandra.session.Query("DROP TABLE api.machines").Exec()
	// _ = cassandra.session.Query("DROP TABLE api.agents").Exec()
	// _ = cassandra.session.Query("DROP TABLE api.logs").Exec()
	// _ = cassandra.session.Query("DROP TABLE api.findings").Exec()
	// _ = cassandra.session.Query("DROP TABLE api.rulefindings").Exec()
	// _ = cassandra.session.Query("DROP TABLE api.exploitcodes").Exec()

	//Create the machines table
	err := cassandra.session.Query("CREATE TABLE IF NOT EXISTS " + cassandra.configuration.CassandraKeyspace + ".machines (id TEXT PRIMARY KEY, os TEXT, hostname TEXT, ip_addresses TEXT)").Exec()
	//Check if an error occured when creating the machines table
	if err != nil {
		return errors.New("cannot create machines table, " + err.Error())
	}

	//Create the agents table
	err = cassandra.session.Query("CREATE TABLE IF NOT EXISTS " + cassandra.configuration.CassandraKeyspace + ".agents (id text PRIMARY KEY, name TEXT, protocol TEXT, ip_address TEXT, port INT, webserver_protocol TEXT, webserver_ip TEXT, webserver_port INT, machine_id TEXT);").Exec()
	//Check if an error occured when creating the agents table
	if err != nil {
		return errors.New("cannot create agents table, " + err.Error())
	}

	//Create the table that will hold the logs of an agent
	err = cassandra.session.Query("CREATE TABLE IF NOT EXISTS " + cassandra.configuration.CassandraKeyspace + ".logs (id TEXT, agent_id TEXT, request_preview TEXT, response_preview TEXT, remote_ip TEXT, timestamp INT, raw_request TEXT, raw_response TEXT, PRIMARY KEY (id, agent_id, timestamp))").Exec()
	//Check if an error occured when creating the table logs
	if err != nil {
		return errors.New("cannot create logs table, " + err.Error())
	}

	//Create the findings table which will hold all the findings of a log
	err = cassandra.session.Query("CREATE TABLE IF NOT EXISTS " + cassandra.configuration.CassandraKeyspace + ".findings (id TEXT, log_id TEXT, line INT, line_index INT, length INT, matched_string TEXT, classification INT, severity INT, validator_name TEXT, finding_type INT, PRIMARY KEY (id, log_id))").Exec()
	//Check if an error occured when creating the findings table
	if err != nil {
		return errors.New("cannot create logs table, " + err.Error())
	}

	//Create the rule findings table which will hold all the rule based findings of a log
	err = cassandra.session.Query("CREATE TABLE IF NOT EXISTS " + cassandra.configuration.CassandraKeyspace + ".rulefindings (id TEXT, log_id TEXT, rule_id TEXT, rule_name TEXT, rule_description TEXT, line INT, line_index INT, length INT, matched_string TEXT, matched_hash TEXT, matched_hash_alg TEXT, classification TEXT, severity INT, finding_type INT, PRIMARY KEY (id, log_id))").Exec()
	//Check if an error occured when creating the rules findings table
	if err != nil {
		return errors.New("cannot create rules findings table, " + err.Error())
	}

	//Create the exploitcodes table which will hold the exploit code of a log
	err = cassandra.session.Query("CREATE TABLE IF NOT EXISTS " + cassandra.configuration.CassandraKeyspace + ".exploitcodes (id TEXT, log_id TEXT, exploit_code TEXT, PRIMARY KEY (id, log_id))").Exec()
	//Check if an error occured when creating the rules findings table
	if err != nil {
		return errors.New("cannot create rules findings table, " + err.Error())
	}

	//Create the index for the agent id in the logs table
	err = cassandra.session.Query("CREATE INDEX IF NOT EXISTS logs_agent_index ON " + cassandra.configuration.CassandraKeyspace + ".logs(agent_id)").Exec()
	if err != nil {
		return errors.New("cannot create agent_id index in logs table, " + err.Error())
	}
	return nil
}

// Function that will initialize the connection to cassandra server
func (cassandra *CassandraConnection) Init() error {
	//Create the configuration for the cassandra cluster
	cluster := gocql.NewCluster(cassandra.configuration.CassandraNodes...)

	//Set the database name (keyspace)
	cluster.Keyspace = "api"

	//Try to connect to the cluster and keyspace
	session, err := cluster.CreateSession()
	if err != nil {
		//Return a new error that will specify that the client couldn't connect to cassandra cluster
		return errors.New("cannot connect to cassandra cluster, " + strings.Join(cassandra.configuration.CassandraNodes, ","))
	}
	//Save the session object in the CassandraConnection struct
	cassandra.session = session

	//Create the needed tables
	err = cassandra.createTables()
	if err != nil {
		return err
	}
	return nil
}

// Check if the machine already exists
func (cassandra *CassandraConnection) CheckMachineExists(os string, hostname string) (string, error) {
	query := cassandra.session.Query("SELECT id FROM "+cassandra.configuration.CassandraKeyspace+".machines WHERE os = ? AND hostname = ?", os, hostname)
	var id string
	iter := query.Iter()
	exists := iter.Scan(&id)
	if !exists {
		return "", nil
	}
	return id, nil
}

// Insert a new machine in the database
func (cassandra *CassandraConnection) InsertMachine(os string, hostname string, ip_addresses []string) (string, error) {
	//Generate a new UUID
	id := uuid.New().String()
	if id == "" {
		return "", errors.New("could not create a new uuid for the proxy")
	}
	//Insert the new agent in the keyspace in the agents table
	err := cassandra.session.Query("INSERT INTO "+cassandra.configuration.CassandraKeyspace+".machines (id, os, hostname, ip_addresses) VALUES (?, ?, ?, ?)", id, os, hostname, strings.Join(ip_addresses[:], " ")).Exec()
	if err != nil {
		return "", errors.New("could not insert a new instance in the database (keyspace), " + err.Error())
	}
	//Get the uuid
	return id, nil
}

// Insert a new agent in the list of agents
func (cassandra *CassandraConnection) InsertAgent(protocol string, ip_address string, port string, webserver_protocol string, webserver_ip string, webserver_port string, machine_id string) (string, error) {
	//Generate a new UUID
	id := uuid.New().String()
	if id == "" {
		return "", errors.New("could not create a new uuid for the proxy")
	}
	//Insert the new agent in the keyspace in the agents table
	err := cassandra.session.Query("INSERT INTO "+cassandra.configuration.CassandraKeyspace+".agents (id, protocol, ip_address, port, webserver_protocol, webserver_ip, webserver_port, machine_id) VALUES (?, ?, ?, ?, ?, ?, ?, ?)", id, protocol, ip_address, port, webserver_protocol, webserver_ip, webserver_port, machine_id).Exec()
	if err != nil {
		return "", errors.New("could not insert a new instance in the database (keyspace), " + err.Error())
	}
	//Get the uuid
	return id, nil
}

// Insert the findings from an agent
func (cassandra *CassandraConnection) InsertFindings(log_id string, findings []data.Finding) (bool, error) {
	//Go through all the findings from the agent
	for _, finding := range findings {
		//Check if the request finding is not an empty object
		if finding.Request != (data.FindingData{}) {
			//Generate a new UUID
			id_request := uuid.New().String()
			//cassandra.logger.Debug(finding.Request)
			//Check if the UUID has been sucessfully generated
			if id_request == "" {
				cassandra.logger.Warning("could not save the request finding - uuid generation failed")
			} else {
				//Insert the request finding
				err := cassandra.session.Query("INSERT INTO "+cassandra.configuration.CassandraKeyspace+".findings (id, log_id, line, line_index, length, matched_string, classification, severity, validator_name, finding_type) VALUES (?,?,?,?,?,?,?,?,?,?)", id_request, log_id, finding.Request.Line, finding.Request.LineIndex, finding.Request.Length, finding.Request.MatchedString, finding.Request.Classification, finding.Request.Severity, finding.Request.ValidatorName, 0).Exec()
				//Check if an error occured when inserting the request finding
				if err != nil {
					cassandra.logger.Error("Could not insert the request finding in the database", err.Error())
				}
			}
		}

		//Check if the finding is not empty object
		if finding.Response != (data.FindingData{}) {
			//Generate a new UUID
			id_response := uuid.New().String()
			//Check if the UUID has been sucessfully generated
			if id_response == "" {
				cassandra.logger.Warning("could not save the response finding - uuid generation failed")
			} else {
				//Insert the response finding
				err := cassandra.session.Query("INSERT INTO "+cassandra.configuration.CassandraKeyspace+".findings (id, log_id, line, line_index, length, matched_string, classification, severity, validator_name, finding_type) VALUES (?,?,?,?,?,?,?,?,?,?)", id_response, log_id, finding.Response.Line, finding.Response.LineIndex, finding.Response.Length, finding.Response.MatchedString, finding.Response.Classification, finding.Response.Severity, finding.Response.ValidatorName, 1).Exec()
				//Check if an error occured when inserting the request finding
				if err != nil {
					cassandra.logger.Error("Could not insert the response finding in the database", err.Error())
				}
			}
		}
	}
	return true, nil
}

// Insert the rule findings from an agent
func (cassandra *CassandraConnection) InsertRuleFindings(log_id string, ruleFindings []data.RuleFinding) (bool, error) {
	for _, ruleFinding := range ruleFindings {
		//If the rule finding has request field defined
		if ruleFinding.Request != nil {
			//Generate a new UUID
			id_request := uuid.New().String()
			//Check if the UUID has been sucessfully generated
			if id_request == "" {
				cassandra.logger.Warning("could not save the request finding - uuid generation failed")
			} else {
				//Insert the request finding
				err := cassandra.session.Query("INSERT INTO "+cassandra.configuration.CassandraKeyspace+".rulefindings (id, log_id, line, line_index, length, matched_string, matched_hash, matched_hash_alg, classification, severity, rule_id, rule_name, rule_description, finding_type) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)", id_request, log_id, ruleFinding.Request.Line, ruleFinding.Request.LineIndex, ruleFinding.Request.Length, ruleFinding.Request.MatchedString, ruleFinding.Request.MatchedBodyHash, ruleFinding.Request.MatchedBodyHashAlg, ruleFinding.Request.Classification, ruleFinding.Request.Severity, ruleFinding.Request.RuleId, ruleFinding.Request.RuleName, ruleFinding.Request.RuleDescription, 0).Exec()
				if err != nil {
					cassandra.logger.Error("Could not insert the request rule findings in the database", err.Error())
				}
			}
		}
		//If the rule finding has the response field defined
		if ruleFinding.Response != nil {
			//Generate a new UUID
			id_response := uuid.New().String()
			//Check if the UUID has been sucessfully generated
			if id_response == "" {
				cassandra.logger.Warning("could not save the request finding - uuid generation failed")
			} else {
				//Insert the request finding
				err := cassandra.session.Query("INSERT INTO "+cassandra.configuration.CassandraKeyspace+".rulefindings (id, log_id, line, line_index, length, matched_string, matched_hash, matched_hash_alg, classification, severity, rule_id, rule_name, rule_description, finding_type) VALUES (?,?,?,?,?,?,?,?,?,?,?,?,?,?)", id_response, log_id, ruleFinding.Response.Line, ruleFinding.Response.LineIndex, ruleFinding.Response.Length, ruleFinding.Response.MatchedString, ruleFinding.Response.MatchedBodyHash, ruleFinding.Response.MatchedBodyHashAlg, ruleFinding.Response.Classification, ruleFinding.Response.Severity, ruleFinding.Response.RuleId, ruleFinding.Response.RuleName, ruleFinding.Response.RuleDescription, 0).Exec()
				if err != nil {
					cassandra.logger.Error("Could not insert the request rule findings in the database", err.Error())
				}
			}
		}
	}

	return true, nil
}

// Insert a log from an agent
func (cassandra *CassandraConnection) InsertLog(logData data.LogData) (bool, error) {
	//Generate a new UUID
	id := uuid.New().String()
	if id == "" {
		return false, errors.New("could not create a new uuid for the log")
	}

	//Convert request from base64 to string
	rawRequest, err := b64.StdEncoding.DecodeString(logData.Request)
	//Check if an error occured when decoding the request from base64
	if err != nil {
		return false, errors.New("could not decode the request from base64, " + err.Error())
	}

	//Create the request preview (the first line of the request)
	request_preview := strings.Split(string(rawRequest), "\n")[0]

	//Convert response from base64 to string
	rawResponse, err := b64.StdEncoding.DecodeString(logData.Response)
	//Check if an error occured when decoding the response from base64
	if err != nil {
		return false, errors.New("could not decode the response from base64, " + err.Error())
	}

	//Create the response preview (the first line of the response)
	response_preview := strings.Split(string(rawResponse), "\n")[0]

	//Insert log data into the database
	err = cassandra.session.Query("INSERT INTO "+cassandra.configuration.CassandraKeyspace+".logs (id, agent_id, request_preview, response_preview, remote_ip, timestamp, raw_request, raw_response) VALUES (?,?,?,?,?,?,?,?)", id, logData.AgentId, request_preview, response_preview, logData.RemoteIP, logData.Timestamp, rawRequest, rawResponse).Exec()
	if err != nil {
		return false, errors.New("could not insert the log in the database, " + err.Error())
	}

	//cassandra.logger.Debug(logData.Findings)

	//Insert all the findings from the agent in the findings table
	cassandra.InsertFindings(id, logData.Findings)

	//Insert all rule findings from the agent in the findings table
	cassandra.InsertRuleFindings(id, logData.RuleFindings)

	return true, nil
}

// Get all the agents from the database
func (cassandra *CassandraConnection) GetAgents() ([]data.Agent, error) {
	query := cassandra.session.Query("SELECT id, name, protocol, ip_address, port, webserver_protocol, webserver_ip, webserver_port, machine_id FROM " + cassandra.configuration.CassandraKeyspace + ".agents")
	agents := make([]data.Agent, 0)
	agent := data.Agent{}
	iter := query.Iter()
	for iter.Scan(&agent.ID, &agent.Name, &agent.ListeningProtocol, &agent.ListeningAddress, &agent.ListeningPort, &agent.ForwardServerProtocol, &agent.ForwardServerAddress, &agent.ForwardServerPort, &agent.MachineId) {
		agents = append(agents, agent)
	}
	return agents, nil
}

// Get a single agent from the database
func (cassandra *CassandraConnection) GetAgent(id string) (data.Agent, error) {
	query := cassandra.session.Query("SELECT id, name, protocol, ip_address, port, webserver_protocol, webserver_ip, webserver_port, machine_id FROM "+cassandra.configuration.CassandraKeyspace+".agents WHERE id = ?", id)
	agent := data.Agent{}
	res := query.Iter().Scan(&agent.ID, &agent.Name, &agent.ListeningProtocol, &agent.ListeningAddress, &agent.ListeningPort, &agent.ForwardServerProtocol, &agent.ForwardServerAddress, &agent.ForwardServerPort, &agent.MachineId)
	if !res {
		return agent, errors.New("agent does not exist")
	}
	return agent, nil
}

// Get number of agents deployed on machine based on id
func (cassandra *CassandraConnection) GetNumberAgentsDeployed(machineId string) (int64, error) {
	//Prepare the query which will get the number of agents deployed on the machine
	query := cassandra.session.Query("SELECT COUNT(id) FROM "+cassandra.configuration.CassandraKeyspace+".agents WHERE machine_id = ? ALLOW FILTERING", machineId)
	var numberAgents int64
	err := query.Scan(&numberAgents)
	if err != nil {
		return 0, errors.New("could not get the number of agents deployed on machine " + machineId + ", " + err.Error())
	}
	return numberAgents, nil
}

// Get all the machines from the database
func (cassandra *CassandraConnection) GetMachines() ([]data.MachineDatabase, error) {
	query := cassandra.session.Query("SELECT id, os, hostname, ip_addresses FROM " + cassandra.configuration.CassandraKeyspace + ".machines")
	machines := make([]data.MachineDatabase, 0)
	machine := data.MachineDatabase{}
	iter := query.Iter()
	var ip_addresses string
	for iter.Scan(&machine.ID, &machine.OS, &machine.Hostname, &ip_addresses) {
		//cassandra.logger.Debug(machine)
		machine.IPAddresses = strings.Split(ip_addresses, " ")
		//Get the number of agents it has deployed
		numberAgents, err := cassandra.GetNumberAgentsDeployed(machine.ID)
		if err != nil {
			cassandra.logger.Warning(err.Error())
		}
		machine.NumberAgents = numberAgents
		machines = append(machines, machine)
	}

	return machines, nil
}

// Get a specific machine based on the id
func (cassandra *CassandraConnection) GetMachine(id string) (data.MachineInformation, error) {
	//Prepare the query which will get the machine based on the id provided
	query := cassandra.session.Query("SELECT os, hostname, ip_addresses FROM "+cassandra.configuration.CassandraKeyspace+".machines WHERE id = ?", id)
	machine := data.MachineInformation{}
	iter := query.Iter()
	var ip_addreses string
	res := iter.Scan(&machine.OS, &machine.Hostname, &ip_addreses)
	if !res {
		return machine, errors.New("could not find the machine with id: " + id)
	}
	//Convert the string of ip_addresses to a list
	machine.IPAddresses = strings.Split(ip_addreses, " ")
	return machine, nil
}

// Delete a specific machine based on the id
func (cassandra *CassandraConnection) DeleteMachine(id string) error {
	//Prepare the query to delete the machine
	err := cassandra.session.Query("DELETE FROM "+cassandra.configuration.CassandraKeyspace+".machines WHERE id = ?", id).Exec()
	return err
}

// Get number of request findings of a log
func (cassandra *CassandraConnection) GetLogFindingsCount(log_id string, finding_type int64) (int64, error) {
	//Prepare the query to select the number of findings for the specified type
	query := cassandra.session.Query("SELECT COUNT(id) FROM "+cassandra.configuration.CassandraKeyspace+".findings WHERE log_id = ? AND finding_type = ? ALLOW FILTERING", log_id, finding_type)
	var count int64
	result := query.Iter().Scan(&count)
	if !result {
		return -1, errors.New("could not get the number of findings for the specified type")
	}
	return count, nil
}

// Get all the findings of the request
func (cassandra *CassandraConnection) GetLogFindings(log_id string) ([]data.FindingDatabase, error) {
	//TO DO...Check which has more findings the request or the response and build the findings array based on the one that has more entries

	//Get the count for request
	request_findings_count, err := cassandra.GetLogFindingsCount(log_id, 0)
	//Check if an error occured when getting the count for request
	if err != nil {
		return nil, err
	}

	//Get the count for response
	response_findings_count, err := cassandra.GetLogFindingsCount(log_id, 1)
	//Check if an error occured when getting the count for request
	if err != nil {
		return nil, err
	}

	//Check which one is larger and initialize the return findings structure array size to that value
	var necessary_structures int64 = 0
	if request_findings_count > response_findings_count {
		necessary_structures = request_findings_count
	} else {
		necessary_structures = response_findings_count
	}

	//Prepare the query to select all the findings on the request of a specific log
	query := cassandra.session.Query("SELECT id, log_id, line, line_index, length, matched_string, classification, severity, validator_name FROM "+cassandra.configuration.CassandraKeyspace+".findings WHERE log_id = ? AND finding_type = 0 ALLOW FILTERING", log_id)
	findings := make([]data.FindingDatabase, necessary_structures)
	findingRequest := data.FindingDataDatabase{}
	iter := query.Iter()
	var index int64 = 0
	for iter.Scan(&findingRequest.Id, &findingRequest.LogId, &findingRequest.Line, &findingRequest.LineIndex, &findingRequest.Length, &findingRequest.MatchedString, &findingRequest.Classification, &findingRequest.Severity, &findingRequest.ValidatorName) {
		//cassandra.logger.Debug(findingRequest)
		findings[index].Request = findingRequest
		index += 1
	}

	//Prepare the query to select all the findings on the response of a specific log
	query = cassandra.session.Query("SELECT id, log_id, line, line_index, length, matched_string, classification, severity, validator_name FROM "+cassandra.configuration.CassandraKeyspace+".findings WHERE log_id = ? AND finding_type = 1 ALLOW FILTERING", log_id)
	findingResponse := data.FindingDataDatabase{}
	iter = query.Iter()
	index = 0
	for iter.Scan(&findingResponse.Id, &findingResponse.LogId, &findingResponse.Line, &findingResponse.LineIndex, &findingResponse.Length, &findingResponse.MatchedString, &findingResponse.Classification, &findingResponse.Severity, &findingResponse.ValidatorName) {
		findings[index].Response = findingResponse
		index += 1
	}

	//Return the data to the client
	return findings, nil
}

// Get number of findings of a certain type of a log
func (cassandra *CassandraConnection) GetLogRuleFindingsCount(log_id string, finding_type int64) (int64, error) {
	//Prepare the query to select the number of findings for the specified type
	query := cassandra.session.Query("SELECT COUNT(id) FROM "+cassandra.configuration.CassandraKeyspace+".rulefindings WHERE log_id = ? AND finding_type = ? ALLOW FILTERING", log_id, finding_type)
	var count int64
	result := query.Iter().Scan(&count)
	if !result {
		return -1, errors.New("could not get the number of findings for the specified type")
	}
	return count, nil
}

// Get all the rule findings of the log
func (cassandra *CassandraConnection) GetLogRuleFindings(log_id string) ([]data.RuleFindingDatabase, error) {
	//Get the count for request
	request_findings_count, err := cassandra.GetLogRuleFindingsCount(log_id, 0)
	//Check if an error occured when getting the count for request
	if err != nil {
		return nil, err
	}

	//Get the count for response
	response_findings_count, err := cassandra.GetLogRuleFindingsCount(log_id, 1)
	//Check if an error occured when getting the count for request
	if err != nil {
		return nil, err
	}

	//Check which one is larger and initialize the return findings structure array size to that value
	var necessary_structures int64 = 0
	if request_findings_count > response_findings_count {
		necessary_structures = request_findings_count
	} else {
		necessary_structures = response_findings_count
	}

	if necessary_structures == 0 {
		return make([]data.RuleFindingDatabase, 0), nil
	}

	//cassandra.logger.Debug("Necessary structures rule findings", necessary_structures)

	//Prepare the query to select all the findings on the request of a specific log
	query := cassandra.session.Query("SELECT id, log_id, line, line_index, length, matched_string, classification, severity, rule_id, rule_name, rule_description, matched_hash, matched_hash_alg FROM "+cassandra.configuration.CassandraKeyspace+".rulefindings WHERE log_id = ? AND finding_type = 0 ALLOW FILTERING", log_id)
	findings := make([]data.RuleFindingDatabase, necessary_structures)
	findingRequest := data.RuleFindingDataDatabase{}
	iter := query.Iter()
	var index int64 = 0
	for iter.Scan(&findingRequest.Id, &findingRequest.LogId, &findingRequest.Line, &findingRequest.LineIndex, &findingRequest.Length, &findingRequest.MatchedString, &findingRequest.Classification, &findingRequest.Severity, &findingRequest.RuleId, &findingRequest.RuleName, &findingRequest.RuleDescription, &findingRequest.MatchedBodyHash, &findingRequest.MatchedBodyHashAlg) {
		//cassandra.logger.Debug(findingRequest)
		findings[index].Request = &findingRequest
		index += 1
	}

	//Prepare the query to select all the findings on the response of a specific log
	query = cassandra.session.Query("SELECT id, log_id, line, line_index, length, matched_string, classification, severity, rule_id, rule_name, rule_description, matched_hash, matched_hash_alg FROM "+cassandra.configuration.CassandraKeyspace+".rulefindings WHERE log_id = ? AND finding_type = 1 ALLOW FILTERING", log_id)
	findingResponse := data.RuleFindingDataDatabase{}
	iter = query.Iter()
	index = 0
	for iter.Scan(&findingResponse.Id, &findingResponse.LogId, &findingResponse.Line, &findingResponse.LineIndex, &findingResponse.Length, &findingResponse.MatchedString, &findingResponse.Classification, &findingResponse.Severity, &findingResponse.RuleId, &findingResponse.RuleName, &findingResponse.RuleDescription, &findingResponse.MatchedBodyHash, &findingResponse.MatchedBodyHashAlg) {
		findings[index].Response = &findingResponse
		index += 1
	}

	//cassandra.logger.Debug(findings)

	//Return the data to the client
	return findings, nil
}

// Get all the logs of an agent in a short format
func (cassandra *CassandraConnection) GetAgentLogsShort(agent_id string) ([]data.LogDataShort, error) {
	//Prepare the query to select all the logs that are generated by the specified agent
	query := cassandra.session.Query("SELECT id, agent_id, request_preview, response_preview, remote_ip, timestamp FROM "+cassandra.configuration.CassandraKeyspace+".logs WHERE agent_id = ?", agent_id)
	logs := make([]data.LogDataShort, 0)
	log := data.LogDataShort{}
	iter := query.Iter()
	for iter.Scan(&log.Id, &log.AgentId, &log.RequestPreview, &log.ResponsePreview, &log.RemoteIP, &log.Timestamp) {
		//Get the all the findings for the log
		findings, err := cassandra.GetLogFindings(log.Id)
		//cassandra.logger.Debug(findings)
		//Check if an error occured when getting the findings for the log
		if err != nil {
			cassandra.logger.Warning("Could not get the findings for the log", log.Id)
			continue
		}
		//Get all the rule findings for the log
		ruleFindings, err := cassandra.GetLogRuleFindings(log.Id)
		//Check if an error occured when getting the rule findings for the log
		if err != nil {
			cassandra.logger.Warning("Could not get the rule findings for the log", log.Id)
			continue
		}
		//Add the findings to the logs findings array
		if len(findings) > 0 {
			log.Findings = findings
		}
		if len(ruleFindings) > 0 {
			log.RuleFindings = ruleFindings
		}
		//Append the selected log to the list of logs
		logs = append(logs, log)
	}

	return logs, nil
}

// Get all logs of an agent
func (cassandra *CassandraConnection) GetAgentLogs(uuid string) ([]data.LogData, error) {
	query := cassandra.session.Query("SELECT id, raw_request, raw_response, remote_ip, timestamp FROM "+cassandra.configuration.CassandraKeyspace+".logs WHERE agent_id = ?", uuid)
	logs := make([]data.LogData, 0)
	log := data.LogData{}
	iter := query.Iter()
	for iter.Scan(&log.Id, &log.Request, &log.Response, &log.RemoteIP, &log.Timestamp) {
		log.AgentId = uuid
		log.Request = b64.StdEncoding.EncodeToString([]byte(log.Request))
		log.Response = b64.StdEncoding.EncodeToString([]byte(log.Response))
		logs = append(logs, log)
	}
	return logs, nil
}

// Get a specific log
func (cassandra *CassandraConnection) GetLog(uuid string) (data.LogDataDatabase, error) {
	query := cassandra.session.Query("SELECT id, raw_request, raw_response, remote_ip, timestamp, request_preview, response_preview FROM "+cassandra.configuration.CassandraKeyspace+".logs WHERE id = ?", uuid)
	log := data.LogDataDatabase{}
	iter := query.Iter()
	for iter.Scan(&log.Id, &log.Request, &log.Response, &log.RemoteIP, &log.Timestamp, &log.RequestPreview, &log.ResponsePreview) {
		log.AgentId = uuid
		log.Request = b64.StdEncoding.EncodeToString([]byte(log.Request))
		log.Response = b64.StdEncoding.EncodeToString([]byte(log.Response))
	}
	return log, nil
}

// Get log request
func (cassandra *CassandraConnection) GetLogRequest(uuid string) (string, error) {
	query := cassandra.session.Query("SELECT raw_request FROM "+cassandra.configuration.CassandraKeyspace+".logs WHERE id = ?", uuid)
	var rawRequest string = ""
	result := query.Iter().Scan(&rawRequest)
	if !result {
		return "", errors.New("could not get the raw request of the log")
	}
	// rawReq, err := b64.StdEncoding.DecodeString(rawRequest)
	// //Check if an error occured when decoding the raw request
	// if err != nil {
	// 	return "", errors.New("could not decode raw request from base64")
	// }
	return rawRequest, nil
}

// Check if exploit code exists for a log
func (cassandra *CassandraConnection) CheckExploitCodeExists(log_uuid string) (bool, error) {
	//Prepare the query which will get the exploit codes number for a log
	query := cassandra.session.Query("SELECT COUNT(id) FROM "+cassandra.configuration.CassandraKeyspace+".exploitcodes WHERE log_id = ? ALLOW FILTERING", log_uuid)
	var countLogExploits int64 = 0
	//Get the result
	result := query.Iter().Scan(&countLogExploits)
	//Check if an error occured
	if !result {
		return false, errors.New("cannot get the number of exploit codes for log " + log_uuid)
	}
	//Return the result
	return countLogExploits > 0, nil
}
