package manager

import (
	"errors"
	"robot/base/log"
	robotfile "robot/base/robot_file"
	"strings"
)

const (
	passIntroKey        = "通过审核记录"
	resultInfo 			= "操作成功"
)

func NewPassReviewInfoHandler() *PassReviewInfoHandler {
	return &PassReviewInfoHandler{
		Key: passIntroKey,
	}
}

type PassReviewInfoHandler struct {
	Key string

	name string
}

func (handler *PassReviewInfoHandler) ParseParm(parm *string) error {
	datas := strings.Split(*parm, "\n")
	if len(datas) < 2 {
		return errors.New("用户输入至少为2行")
	}

	record := datas[1]
	if len(record) < 1 {
		return errors.New("至少指定一个文件")
	}

	handler.name = record
	return nil
}

func (handler *PassReviewInfoHandler) HandlerRequest() (*string, error) {
    err := robotfile.MoveFile(robotfile.ReviewPath + handler.name, robotfile.HeroInfoPath + handler.name)
	if err != nil {
		log.Error.Println("审核通过错误：" + err.Error())
		return nil, errors.New("不存在该审核记录")
	}

	result := resultInfo
	return &result, nil
}	