package data

import (
	"encoding/json"
	"io"
)

type SuccessMessage struct {
	Message string `json:"message"`
}

func (sm *SuccessMessage) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(sm)
}

func (sm *SuccessMessage) ToJSON(w io.Writer) error {
	d := json.NewEncoder(w)
	return d.Encode(sm)
}
