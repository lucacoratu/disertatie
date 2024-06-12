/*
    ID           string   `json:"id"`           // The ID of the machine
	OS           string   `json:"os"`           //The operating system of the machine
	Hostname     string   `json:"hostname"`     //Hostname of the machine
	IPAddresses  []string `json:"ip_addresses"` //A list of ip addreses of the machine on all network interfaces
	NumberAgents int64    `json:"numberAgents"` //The number of agents the machine has on itself
}
*/
type Machine = {
    id: string,
    os: string,
    hostname: string,
    ipAddresses: string[],
    numberAgents: number
}

type MachineResponse = {
    machines: Machine[]
}

interface MachineProps {
    machine: Machine
}

interface OSProps {
    os: string
}

type MachineRegisterResponse = {
    uuid: string
}

// type MachinesStatisticsResponse struct {
// 	TotalMachines   int64 `json:"totalMachines"`
// 	TotalInterfaces int64 `json:"totalInterfaces"`
// }

type MachinesStatisticsResponse = {
    totalMachines: number,
    totalInterfaces: number
}