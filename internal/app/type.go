package app

import (
	"encoding/json"
	"time"
)

type Message struct {
	Token     string
	Success   bool
	StartTime time.Time
	Error     error
}

func (m *Message) MarshalJSON() ([]byte, error) {
	var errMsg = ""
	if m.Error != nil {
		errMsg = m.Error.Error()
	}

	msg := struct {
		Token    string `json:"token"`
		Success  bool   `json:"success"`
		Duration string `json:"duration"`
		Error    string `json:"error"`
	}{
		m.Token,
		m.Success,
		time.Since(m.StartTime).String(),
		errMsg,
	}

	return json.Marshal(msg)
}
