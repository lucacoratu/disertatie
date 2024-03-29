package database

import "github.com/lucacoratu/disertatie/api/data"

type IConnection interface {
	Init() error
	CheckMachineExists(os string, hostname string) (string, error)
	InsertMachine(os string, hostname string, ip_addresses []string) (string, error)
	InsertAgent(protocol string, ip_address string, port string, webserver_protocol string, webserver_ip string, webserver_port string, machine_id string) (string, error)
	GetAgents() ([]data.Agent, error)
	GetAgent(id string) (data.Agent, error)
	ModifyAgent(id string, agent data.UpdateAgent) error
	GetMachines() ([]data.MachineDatabase, error)
	GetMachine(id string) (data.MachineInformation, error)
	DeleteMachine(id string) error
	GetAgentLogs(uuid string) ([]data.LogData, error)
	GetAgentLogsShort(agent_id string) ([]data.LogDataShort, error)
	InsertLog(logData data.LogData) (bool, error)
	GetLog(uuid string) (data.LogDataDatabase, error)
	GetLogRequest(uuid string) (string, error)
	GetLogFindings(log_uuid string) ([]data.FindingDatabase, error)
	GetLogRuleFindings(log_uuid string) ([]data.RuleFindingDatabase, error)
	CheckExploitCodeExists(log_uuid string) (bool, error)
}
