package router

import "time"

//OpCode
type OpCode int32
const(
	Dispatch 		OpCode = 0
	Heartbeat 		OpCode = 1
	Identify 		OpCode = 2
	Resume   		OpCode = 6
	Reconnect 		OpCode = 7
	InvalidSession 	OpCode = 9
	Hello 			OpCode = 10
	HeartbeatACK 	OpCode = 11
	HTTPCallbackACK OpCode = 12
)

//Key
const(
	MsgCreateKey   string = "AT_MESSAGE_CREATE"
)

//Intents
const(
	PUBLIC_GUILD_MESSAGES int = 1 << 30
)

var Connections chan bool
var HeartBeatExit chan bool

//websocket所有报文的统一结构
type EventPayLoad struct {
	Op OpCode 			`json:"op"` 		//opcode,用于区分操作类型
	D interface{} 		`json:"d"`			//不同事件的报文保存的内容不同
	S int 				`json:"s"`			//下行消息的唯一标识
	T string 			`json:"t"`			//事件类型
}

//获取wsURL
type WsUrlInfo struct {
	Url string `json:"url"`
}

//心跳内容
type HeartBeat struct {
	HeartBeatInterval int `json:"heartbeat_interval"`
}

//鉴权连接的内容
type AuthInfo struct {
	Token string 			`json:"token"`
	Intents int 			`json:"intents"`
	Shard []int 			`json:"shard"`
	Properties Properties 	`json:"properties"`
}

type Properties struct {
	Os string 			`json:"$os"`
	Browser string 		`json:"$browser"`
	Device string 		`json:"$device"`
}

type ReadyPayLoad struct {
	Op OpCode 			`json:"op"` 	
	D ReadyInfo			`json:"d"`		
	S int 				`json:"s"`		
	T string 			`json:"t"`		
}

//准备的内容
type ReadyInfo struct {
	Version int 		`json:"version"`
	SessionID string 	`json:"session_id"`
	User User 			`json:"user"`
	Shard []int 		`json:"shard"`
}

type User struct {
	ID string 			`json:"id"`
	Username string 	`json:"username"`
	Bot bool 			`json:"bot"`
}

//消息内容
type MsgEventPayLoad struct {
	Op OpCode 		`json:"op"` 		
	D Message 		`json:"d"`		
	S int 			`json:"s"`		
	T string 		`json:"t"`			
}

type Message struct {
	ID string 						`json:"id"`
	Author Author 					`json:"author"`
	ChannelID string 				`json:"channel_id"`
	Content string 					`json:"content"`
	GuildID string 					`json:"guild_id"`
	Timestamp time.Time 			`json:"timestamp"`
	MsgReference MessageReference 	`json:"message_reference"`
}

type Author struct {
	Avatar string 		`json:"avatar"`
	Bot bool 			`json:"bot"`
	ID string 			`json:"id"`
	Username string 	`json:"username"`
}

type MessageReference struct {
	MessageId string				`json:"message_id"`
	IgnoreGetMessageError bool	`json:"ignore_get_message_error"`
}

type Member struct {
	JoinedAt time.Time 	`json:"joined_at"`
	Roles []string 		`json:"roles"`
}

type ResumeInfo struct {
	Token string 		`json:"token"`
	SessionID string 	`json:"session_id"`
	Seq int 			`json:"seq"`
}
