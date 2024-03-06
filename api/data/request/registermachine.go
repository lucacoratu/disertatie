package data

import (
	"encoding/json"
	"io"
)

// Holds the data from the WebUI for registering a new machine
type RegisterMachineRequest struct {
	OS         string `json:"os"`         //The operating system of the machine
	Hostname   string `json:"hostname"`   //Hostname of the machine
	IPAddress  string `json:"ipAddress"`  //A list of ip addreses of the machine on all network interfaces
	Username   string `json:"username"`   //The username used for connection (SSH)
	Password   string `json:"password"`   //The password used for connection (SSH)
	PrivateKey string `json:"privateKey"` //The private key used for connection (SSH)
}

func (rmr *RegisterMachineRequest) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(rmr)
}

func (rmr *RegisterMachineRequest) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(rmr)
}
