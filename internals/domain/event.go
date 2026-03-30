package domain

import "encoding/json"

type DocEvent struct {
	DocID  string `json:"doc_id"`
	Status string `json:"status"`
	UserID int64  `json:"user_id"`
}

func (event DocEvent) MarshalBinary() ([]byte, error) {
	return json.Marshal(event)
}
