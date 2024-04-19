package data

import (
	"encoding/json"
	"io"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

func (lr *LoginRequest) FromJSON(r io.Reader) error {
	d := json.NewDecoder(r)
	return d.Decode(lr)
}
