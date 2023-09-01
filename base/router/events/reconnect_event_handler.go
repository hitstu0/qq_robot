package events

import (
	"robot/base/router"
)
type ReconnectEventHandler struct {
	EventItent int
	EventKey string
	OpCode router.OpCode
}

func (handler *ReconnectEventHandler) GetEventItent() int {
	return handler.EventItent
}

func (handler *ReconnectEventHandler) GetEventKey() string {
	return handler.EventKey
}

func (handler *ReconnectEventHandler) GetOpCode() router.OpCode {
	return handler.OpCode
}

func (handler *ReconnectEventHandler) HandleEvent(content *[]byte)  {
	router.Connections <- true
	router.HeartBeatExit <- true
}