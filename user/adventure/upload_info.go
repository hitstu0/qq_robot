package adventure

import (
	"errors"
	robotfile "robot/base/robot_file"
	"strings"
)

const (
	uploadKey = "上传英雄介绍"

	successInfo = "保存成功，等待管理员审核"
)

func NewUploadInfoHandler() *UploadInfoHandler {
	return &UploadInfoHandler{
		Key: uploadKey,
	}
}

type UploadInfoHandler struct {
	Key string

	name         string
	introduction string
}

func (handler *UploadInfoHandler) ParseParm(parm *string) error {
	handler.init()

	datas := strings.Split(*parm, "\n")
	if len(datas) < 3 {
		return errors.New("用户输入至少为3行")
	}

	name := datas[1]
	introduction := datas[2]
	if len(name) < 1 || len(introduction) < 1 {
		return errors.New("名字和介绍不能为空")
	}

	handler.name = name
	handler.introduction = introduction
	return nil
}

func (handler *UploadInfoHandler) HandlerRequest() (*string, error) {
	url := robotfile.ReviewPath + handler.name
	err := robotfile.WriteFile(url, handler.introduction)
	if err != nil {
		return nil, err
	}

	result := successInfo
	return &result, nil
}

func (handler *UploadInfoHandler) init() {
	handler.name = ""
	handler.introduction = ""
}
