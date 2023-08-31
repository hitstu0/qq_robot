package events

import (
	"encoding/json"
	"errors"
	"robot/base/answer"
	"robot/base/log"
	"robot/base/robot_data"
	"robot/base/router"
	"robot/user"
	"robot/user/adventure"
	"robot/user/introduction"
	"robot/user/manager"
	"strings"
)

var (
	userHandlers map[string]user.UserFuncHandler
)

//初始化这个事件支持各个功能，包括获取使用说明，生成冒险故事等
func init() {
	userHandlers = make(map[string]user.UserFuncHandler)

	introHandler := introduction.NewIntroHandler()
	adventureHandler := adventure.NewAdventureHandler()
	uploadInfoHandler := adventure.NewUploadInfoHandler()
	getHeroListHandler := adventure.NewGetHeroListHandler()
	getReviewHandler := manager.NewGetReviewInfoHandler()
	getReviewDetailHandler := manager.NewGetReviewDetailHandler()
	passReviewHandler := manager.NewPassReviewInfoHandler()
	denyReviewHandler := manager.NewDenyReviewInfoHandler()


	userHandlers[introHandler.Key] = introHandler
	userHandlers[adventureHandler.Key] = adventureHandler
	userHandlers[uploadInfoHandler.Key] = uploadInfoHandler
	userHandlers[getHeroListHandler.Key] = getHeroListHandler
	userHandlers[getReviewHandler.Key] = getReviewHandler
	userHandlers[passReviewHandler.Key] = passReviewHandler
	userHandlers[denyReviewHandler.Key] = denyReviewHandler
	userHandlers[getReviewDetailHandler.Key] = getReviewDetailHandler
}

type MsgEventHandler struct {
	EventItent int
	EventKey string
	AuthKey      string
	OpCode router.OpCode
}

func (handler *MsgEventHandler) GetEventItent() int {
	return handler.EventItent
}

func (handler *MsgEventHandler) GetEventKey() string {
	return handler.EventKey
}

func (handler *MsgEventHandler) GetOpCode() router.OpCode {
	return handler.OpCode
}

func (handler *MsgEventHandler) HandleEvent(eventLoad *[]byte)  {
	//解析消息对象
	msg, err := handler.getMsg(eventLoad)
	if err != nil {
		answer.AnswerMsg(&msg.ChannelID, err.Error(), 
			&msg.MsgReference.MessageId, &robotdata.GetRobotData().AuthKey)
		log.Error.Println("解析消息错误 " + err.Error())
		return 
	}
	
	//从用户输入解析出Key
	key, err := handler.parseContent(msg.Content)
	if err != nil {
		answer.AnswerMsg(&msg.ChannelID, err.Error(), 
			&msg.MsgReference.MessageId, &robotdata.GetRobotData().AuthKey)
		log.Error.Println("解析Key错误 " + err.Error())
		return
	}

	//查找处理事件的handler
	curHandler, err := handler.getHandler(key)
	if err != nil {
		answer.AnswerMsg(&msg.ChannelID, err.Error(), 
			&msg.MsgReference.MessageId, &robotdata.GetRobotData().AuthKey)
		log.Error.Println("获取处理器错误: " + err.Error())
		return 
	}

	//handler解析参数，处理并返回结果
	err = curHandler.ParseParm(&msg.Content)
	if err != nil {
		answer.AnswerMsg(&msg.ChannelID, err.Error(), 
			&msg.MsgReference.MessageId, &robotdata.GetRobotData().AuthKey)
		log.Error.Println("解析参数错误: " + err.Error())
		return 
	}

	result, err := curHandler.HandlerRequest()
	if err != nil {
		answer.AnswerMsg(&msg.ChannelID, err.Error(), 
			&msg.MsgReference.MessageId, &robotdata.GetRobotData().AuthKey)
		log.Error.Println("处理结果错误: " + err.Error())
		return 
	}

	//封装返回的结果
	err = answer.AnswerMsg(&msg.ChannelID, *result, 
		&msg.MsgReference.MessageId, &robotdata.GetRobotData().AuthKey)
	if err != nil {
		errs := "请求GTP超时，请重试"
		answer.PushMsg(&msg.ChannelID, errs, &robotdata.GetRobotData().AuthKey)
	}
}

func (handler *MsgEventHandler) getMsg(content *[]byte) (router.Message, error){
	var eventLoad router.MsgEventPayLoad
	err := json.Unmarshal(*content, &eventLoad)
	if err != nil {
		log.Error.Println("反序列化event事件错误 " + err.Error())
		return router.Message{}, err
	}

	return eventLoad.D, nil
}

func (handler *MsgEventHandler) parseContent(text string) (string, error) {
	var firstLine string
	if strings.Contains(text, "\n") {
		firstLine = strings.Split(text, "\n")[0]
	} else {
		firstLine = text
	}

	datas := strings.SplitN(firstLine, " ", 2)
	if len(datas) < 2 {
		return "", errors.New("请求格式错误")
	}
	
	return strings.TrimSpace(datas[1]), nil
}

func (handler *MsgEventHandler) getHandler(key string) (user.UserFuncHandler, error){
	cur, ok := userHandlers[key]
	if !ok {
		return nil, errors.New("暂不支持该功能：" + key + "。请查看使用说明")
	}

	return cur, nil
}