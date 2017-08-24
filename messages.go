package jsonstream

import (
	"encoding/json"
	"github.com/aphistic/gomol"
	"time"
)

type Message interface{}

type outputMsg struct {
	Type string      `json:"type"`
	Data interface{} `json:"data"`
}
type partialMsg struct {
	Type string          `json:"type"`
	Data json.RawMessage `json:"data"`
}

type HeaderMsg struct {
	Timestamp time.Time              `json:"time"`
	Headers   map[string]interface{} `json:"headers"`
}

type LogMsg struct {
	Timestamp time.Time              `json:"time"`
	Msg       string                 `json:"msg"`
	Level     gomol.LogLevel         `json:"level"`
	Attrs     map[string]interface{} `json:"attrs"`
}
