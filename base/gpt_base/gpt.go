package gptBase

const(
	UserRole = "user"
)

var (
	Client GptClient
)

//消息封装结构体
type GptMsg struct{
	Role string `json:"role"`
	Content string `json:"content"`
}



//统一的大语言模型接口
type GptClient interface {
	SendMessage(gptMsgs[] *GptMsg) (*string, error) 
}