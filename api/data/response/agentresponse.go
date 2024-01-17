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
