package odin

import (
	"bytes"
	"encoding/json"
)

type request struct {
	Domain string        `json:"domain"`
	Limit  int           `json:"limit"`
	Start  []interface{} `json:"start"`
}

func (r *request) ToJSON() (*bytes.Buffer, error) {
	var buf bytes.Buffer
	err := json.NewEncoder(&buf).Encode(r)
	return &buf, err
}
