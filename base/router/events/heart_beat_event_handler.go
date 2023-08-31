package events

import (
	"robot/base/log"
	"robot/base/router"
)
type HeartBeatEventHandler struct {
	EventItent int
	EventKey string
	OpCode router.OpCode
}

func (handler *HeartBeatEventHandler) GetEventItent() int {
	return handler.EventItent
}

func (handler *HeartBeatEventHandler) GetEventKey() string {
	return handler.EventKey
}

func (handler *HeartBeatEventHandler) GetOpCode() router.OpCode {
	return handler.OpCode
}

func (handler *HeartBeatEventHandler) HandleEvent(content *[]byte)  {
	log.Info.Println("心跳成功")
}