package manager

import (
	"errors"
	"robot/base/log"
	"robot/base/robot_file"
	"strings"
)

const (
	reviewDetailKey string = "获取审核记录详情"
)

func NewGetReviewDetailHandler() *GetReviewDetailHandler {
	return &GetReviewDetailHandler{
		Key: reviewDetailKey,
	}
}

type GetReviewDetailHandler struct {
	Key  string
	name string
}

func (handler *GetReviewDetailHandler) ParseParm(parm *string) error {
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

func (handler *GetReviewDetailHandler) HandlerRequest() (*string, error) {
	content, err := robotfile.ReadFile(robotfile.ReviewPath + handler.name)
	if err != nil {
		log.Error.Println("读取文件错误：" + robotfile.ReviewPath)
		return nil, errors.New("不存在该审核记录")
	}

	return content, nil
}