package database

import "github.com/lucacoratu/disertatie/api/data"

type IConnection interface {
	Init() error
	CheckMachineExists(os string, hostname string) (string, error)
	InsertMachine(os string, hostname string, ip_addresses []string) (string, error)
	InsertAgent(protocol string, ip_address string, port string, webserver_protocol string, webserver_ip string, webserver_port string, machine_id string) (string, error)
	GetAgents() ([]data.Agent, error)
	GetMachines() ([]data.MachineDatabase, error)
	GetMachine(id string) (data.MachineInformation, error)
	GetAgentLogs(uuid string) ([]data.LogData, error)
	GetAgentLogsShort(agent_id string) ([]data.LogDataShort, error)
	InsertLog(logData data.LogData) (bool, error)
	GetLog(uuid string) (data.LogDataDatabase, error)
	GetLogFindings(log_uuid string) ([]data.FindingDatabase, error)
}
