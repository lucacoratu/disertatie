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

type MachinesRegisterResponse struct {
	Uuid string `json:"uuid"` //The UUID of the created machine
}

func (mrr *MachinesRegisterResponse) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(mrr)
}

type MachinesStatisticsResponse struct {
	TotalMachines   int64 `json:"totalMachines"`
	TotalInterfaces int64 `json:"totalInterfaces"`
}

func (msr *MachinesStatisticsResponse) FromJSON(r io.Reader) error {
	e := json.NewDecoder(r)
	return e.Decode(msr)
}

func (msr *MachinesStatisticsResponse) ToJSON(w io.Writer) error {
	d := json.NewEncoder(w)
	return d.Encode(msr)
}
