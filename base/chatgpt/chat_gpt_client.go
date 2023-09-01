package chatgpt

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"net/url"
	"robot/base/gpt_base"
	"robot/base/log"
	robothttp "robot/base/robot_http"
)

const(
	//oglejimneto@mail.com----kXwqKbxteo----O3wTLO528----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
	authKey = "Bearer xxxxxxxxxxxxxxx"
	model = "gpt-3.5-turbo"
	gptUrl = "https://api.openai.com/v1/chat/completions"
	temperature = 0.7

)

type ChatGptClient struct{
	AuthKey string
	Model string
	GptUrl string
	Temperature float64
	HttpClient *http.Client
}

//采用默认设置的ChatGPT客户端
func NewDefaultChatGptClient() (*ChatGptClient, error) {
	gptClient := &ChatGptClient{
		AuthKey: authKey,
		Model: model,
		GptUrl: gptUrl,
	}

	proxyUrl, err := url.Parse("http://127.0.0.1:18081")
	if(err != nil) {
		log.Error.Println("获取代理URL错误: " + err.Error())
		return nil, fmt.Errorf("proxy url error: %s", err.Error()) 
	}

	gptClient.HttpClient = &http.Client{
		Transport: &http.Transport{
			// 设置代理
			Proxy: http.ProxyURL(proxyUrl),
		},
	}
	
	return gptClient, nil
}

//使用GPT客户端向ChatGPT3.5API接口发起请求，并接收响应
func (chatClient * ChatGptClient) SendMessage(gptMsgs[] *gptBase.GptMsg) (*string, error) {
	if len(gptMsgs) == 0 {
		log.Error.Println("消息长度不能为0")
		return nil, errors.New("gptMsgs len should bigger than 0")
	}

	log.Info.Println("开始发送")
	//创建请求体
	gptRequest := &chatGptRequest{
		Model: chatClient.Model,
		Temperature: chatClient.Temperature,
		Messages: gptMsgs,
	}
	reqByte, _ := json.Marshal(gptRequest)
	reqString := string(reqByte)

	//创建http请求并接收响应
	resp, err := robothttp.DoPOST(chatClient.HttpClient, &chatClient.GptUrl, &reqString, &chatClient.AuthKey)
	if err != nil {
		log.Error.Println("获取GPT响应错误：" + err.Error())
		return nil, err
	}
	log.Info.Println("接收响应：" + *resp)

	var gptResp chatGptResponse
	err = json.Unmarshal([]byte(*resp), &gptResp)
	if err != nil {
		log.Error.Println("反序列化GPT响应错误：" + err.Error())
		return nil, err
	}

	if len(gptResp.Choices) < 1 {
		log.Error.Println("GPT通信错误")
		return nil, errors.New("GPT通信错误")
	}
	
	return &gptResp.Choices[0].Message.Content, nil
}
