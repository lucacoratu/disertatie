package data

import (
	"encoding/json"
	"io"
)

type RegisterProxyResponse struct {
	Uuid string `json:"uuid"`
}

func (reg *RegisterProxyResponse) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(reg)
}

func (reg *RegisterProxyResponse) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(reg)
}
