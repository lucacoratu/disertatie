package data

import (
	"encoding/json"
	"io"
)

const (
	DATABASE_ERROR int64 = 1
	PARSE_ERROR    int64 = 2
	REQUEST_ERROR  int64 = 3
)

type APIError struct {
	Code    int64  `json:"code"`    //The code of the error
	Message string `json:"message"` //The message of the error
}

// Convert from json into the APIError structure
func (err *APIError) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(err)
}

// Convert to json the APIError structure
func (err *APIError) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(err)
}
