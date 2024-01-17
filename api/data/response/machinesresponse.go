package data

import (
	"encoding/json"
	"io"

	"github.com/lucacoratu/disertatie/api/data"
)

type MachinesGetResponse struct {
	Machines []data.MachineDatabase `json:"machines"` // The list of agents
}

func (ag *MachinesGetResponse) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(ag)
}

func (ag *MachinesGetResponse) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(ag)
}
