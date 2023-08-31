package chatgpt

import "robot/base/gpt_base" 

//GPT API使用的请求结构体
type chatGptRequest struct {
	Model string 				`json:"model"`
	Messages []*gptBase.GptMsg 	`json:"messages"`
	Temperature float64 		`json:"temperature"`
}

//GPT API使用的响应结构体
type chatGptResponse struct {
	Id string 			`json:"id"`
	Object string 		`json:"object"`
	Created int 		`json:"created"`
	Model string 		`json:"model"`
	Usage usage 		`json:"usage"`
	Choices []choice 	`json:"choices"`
}
type usage struct {
	PromptTokens int 		`json:"prompt_tokens"`
	CompletionTokens int 	`json:"completion_tokens"`
	TotalTokens int 		`json:"total_tokens"`
}

type choice struct {
	Message gptBase.GptMsg 	`json:"message"`
	FinishReason string 	`json:"finish_reason"`
	Index int 				`json:"index"`
}