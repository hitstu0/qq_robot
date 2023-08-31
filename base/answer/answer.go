package answer

import (
	"encoding/json"
	"fmt"
	"net/http"
	"robot/base"
	"robot/base/log"
	robothttp "robot/base/robot_http"
)

const (
	answerUrlPre    string = "channels"
	answerUrlTail   string = "messages"
	contentType     string = "application/json"
)

var (
	client *http.Client
)

func init() {
	client = &http.Client{}
}

func PushMsg(channelId *string, content string, auth *string) error {
	data := getPushBody(&content)

	url := fmt.Sprintf("%s%s/%s/%s", base.GetUrl(), answerUrlPre, *channelId, answerUrlTail)
	log.Info.Println(url)
	_, err := robothttp.DoPOST(client, &url, data, auth)
	if err != nil {
		log.Error.Println("返回响应错误: " + err.Error())
		return err
	}

	return nil
}

func AnswerMsg(channelId *string, content string, replyMsgId *string, auth *string) error {
	data := getAnswerBody(&content, replyMsgId)

	url := fmt.Sprintf("%s%s/%s/%s", base.GetUrl(), answerUrlPre, *channelId, answerUrlTail)
	log.Info.Println(url)
	_, err := robothttp.DoPOST(client, &url, data, auth)
	if err != nil {
		log.Error.Println("返回响应错误: " + err.Error())
		return err
	}

	return nil
}

func getAnswerBody(content *string, replyMsgId *string) *string {
	answerData := answerData{
		Content: *content,
		MsgId:   *replyMsgId,
	}
	answerJson, _ := json.Marshal(answerData)
	answerString := string(answerJson)
	return &answerString
}

func getPushBody(content *string) *string {
	answerData := answerData{
		Content: *content,
	}
	answerJson, _ := json.Marshal(answerData)
	answerString := string(answerJson)
	return &answerString
}