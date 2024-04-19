package data

import (
	"encoding/json"
	"io"
)

type LoginResponse struct {
	Token string `json:"token"`
}

func (lr *LoginResponse) ToJSON(w io.Writer) error {
	e := json.NewEncoder(w)
	return e.Encode(lr)
}
