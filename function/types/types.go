package types

import (
	"encoding/json"
	"errors"
	"io"
)

type TaskResult struct {
	ID   string      `json:"id"`
	Data interface{} `json:"data"`
}

type PollRequest struct {
	UUID   string  `json:"uuid"`
	CPU    float64 `json:"cpu"`
	Memory float64 `json:"memory"`
}

type TaskRequest struct {
	ID   string                 `json:"id"`
	Type string                 `json:"type"`
	Args map[string]interface{} `json:"args"`
}

func (tr *TaskRequest) ParseBody(body io.Reader) error {
	if err := json.NewDecoder(body).Decode(tr); err != nil {
		return err
	}
	return nil
}

type RegisterResponse struct {
	UUID string
}

func (rr *RegisterResponse) ParseBody(body io.Reader) error {
	var uuid string
	if err := json.NewDecoder(body).Decode(&uuid); err != nil {
		return err
	}
	if len(uuid) == 0 {
		return errors.New("invalid meeting")
	}
	rr.UUID = uuid
	return nil
}
