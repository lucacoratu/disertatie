package data

import (
	"encoding/json"
	"io"

	"github.com/lucacoratu/disertatie/api/data"
)

type LogsGetResponse struct {
	Logs []data.LogDataShort `json:"logs"` // The list of logs
}

func (logs *LogsGetResponse) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(logs)
}

func (logs *LogsGetResponse) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(logs)
}

type LogGetResponse struct {
	Log data.LogDataDatabase `json:"log"` //The log to be returned
}

func (logs *LogGetResponse) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(logs)
}

func (logs *LogGetResponse) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(logs)
}
