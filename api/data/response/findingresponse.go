package data

import (
	"encoding/json"
	"io"

	"github.com/lucacoratu/disertatie/api/data"
)

type FindingClassificationStringResponse struct {
	FindingsString []data.FindingClassificationString `json:"findingsString"`
}

func (fs *FindingClassificationStringResponse) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(fs)
}

func (fs *FindingClassificationStringResponse) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(fs)
}
