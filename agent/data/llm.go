package data

import (
	"encoding/json"
	"io"
)

type LLMResponse struct {
	Headers map[string]string `json:"headers"`
	Body    string            `json:"body"`
}

func (lr *LLMResponse) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(lr)
}
