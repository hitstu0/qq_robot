package router

//事件处理的接口，统一抽象事件处理
type EventHandler interface {
	GetEventItent() int
	GetEventKey() string
	GetOpCode() OpCode
	HandleEvent(content *[]byte) 
}