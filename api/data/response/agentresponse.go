package data

import (
	"encoding/json"
	"io"

	"github.com/lucacoratu/disertatie/api/data"
)

type AgentsGetResponse struct {
	Agents []data.Agent `json:"agents"` // The list of agents
}

func (ag *AgentsGetResponse) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(ag)
}

func (ag *AgentsGetResponse) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(ag)
}

type AgentsCountResponse struct {
	Count int64 `json:"count"` //The number of agents registered on the API
}

func (acr *AgentsCountResponse) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(acr)
}

func (acr *AgentsCountResponse) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(acr)
}

type ConnectedAgentsResponse struct {
	Agents []string `json:"agents"`
}

func (car *ConnectedAgentsResponse) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(car)
}

func (car *ConnectedAgentsResponse) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(car)
}
