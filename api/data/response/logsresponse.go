package data

import (
	"encoding/json"
	"io"

	"github.com/lucacoratu/disertatie/api/data"
)

type LogsGetResponse struct {
	Logs     []data.LogDataShort `json:"logs"`     // The list of logs
	NextPage string              `json:"nextPage"` //The next page to navigate to
}

func (logs *LogsGetResponse) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(logs)
}

func (logs *LogsGetResponse) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(logs)
}

type LogsGetResponseElastic struct {
	Logs []data.LogDataElastic `json:"logs"` // The list of logs
}

func (logs *LogsGetResponseElastic) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(logs)
}

func (logs *LogsGetResponseElastic) FromJSON(r io.Reader) error {
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

type TotalLogCountResponse struct {
	Count int64 `json:"count"`
}

func (tlcr *TotalLogCountResponse) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(tlcr)
}

func (tlcr *TotalLogCountResponse) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(tlcr)
}
