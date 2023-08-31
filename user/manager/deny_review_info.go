package manager

import (
	"errors"
	"robot/base/log"
	robotfile "robot/base/robot_file"
	"strings"
)

const (
	denyIntroKey        = 		"拒绝审核记录"
	denyResultInfo 		= 		"操作成功"
)

func NewDenyReviewInfoHandler() *DenyReviewInfoHandler {
	return &DenyReviewInfoHandler{
		Key: denyIntroKey,
	}
}

type DenyReviewInfoHandler struct {
	Key string
	name string
}

func (handler *DenyReviewInfoHandler) ParseParm(parm *string) error {
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

func (handler *DenyReviewInfoHandler) HandlerRequest() (*string, error) {
    err := robotfile.DeleteFile(robotfile.ReviewPath + handler.name)
	if err != nil {
		log.Error.Println("审核拒绝错误：" + err.Error())
		return nil, errors.New("不存在该审核记录")
	}

	result := denyResultInfo
	return &result, nil
}	