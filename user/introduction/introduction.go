package introduction

import (
	"robot/base/robot_file"
	"robot/base/log"
	"time"
)

const (
	introKey string = "使用说明"
	expiration int64  = 10 * 1000
)

func NewIntroHandler() *IntroHandler {
	return &IntroHandler{
		Key: introKey,
	}
}

type IntroHandler struct {
	Key string
	introduction string
	lastTime int64
}

func (handler *IntroHandler) ParseParm(parm *string) error {
	return nil
}

func (handler *IntroHandler) HandlerRequest() (*string, error) {
	if(time.Now().UnixMilli() - handler.lastTime < expiration) {
		return &handler.introduction, nil
	} else {
		content, err := robotfile.ReadFile(robotfile.RobotInfoPath)
		if err != nil {
			log.Error.Println("读取文件错误：" + robotfile.RobotInfoPath)
			return nil, err
		}


		handler.introduction = *content
		handler.lastTime = time.Now().UnixMilli()
		
		return content, nil
	}
	
}