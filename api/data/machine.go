package data

import (
	"encoding/json"
	"io"
)

const (
	UNKOWN  int64 = 0
	WINDOWS int64 = 1
	LINUX   int64 = 2
)

// Holds data collected by the proxy about the machine
type MachineInformation struct {
	OS          string   `json:"os"`          //The operating system of the machine
	Hostname    string   `json:"hostname"`    //Hostname of the machine
	IPAddresses []string `json:"ipAddresses"` //A list of ip addreses of the machine on all network interfaces
}

func (mach *MachineInformation) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(mach)
}

func (mach *MachineInformation) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(mach)
}

type MachineDatabase struct {
	ID           string   `json:"id"`           // The ID of the machine
	OS           string   `json:"os"`           //The operating system of the machine
	Hostname     string   `json:"hostname"`     //Hostname of the machine
	IPAddresses  []string `json:"ipAddresses"`  //A list of ip addreses of the machine on all network interfaces
	NumberAgents int64    `json:"numberAgents"` //The number of agents the machine has on itself
}

func (mach *MachineDatabase) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(mach)
}

func (mach *MachineDatabase) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(mach)
}
