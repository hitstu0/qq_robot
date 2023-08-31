package manager

import (
	"robot/base/log"
	robotfile "robot/base/robot_file"
)

const (
	getRIntroKey        = "获取审核列表"

	resultPre 			= "审核列表如下:\n"
)

func NewGetReviewInfoHandler() *GetReviewInfoHandler {
	return &GetReviewInfoHandler{
		Key: getRIntroKey,
	}
}

type GetReviewInfoHandler struct {
	Key string
}

func (handler *GetReviewInfoHandler) ParseParm(parm *string) error {
	return nil
}

func (handler *GetReviewInfoHandler) HandlerRequest() (*string, error) {
    result, err := robotfile.GetAllFileName(robotfile.ReviewPath)
	if err != nil {
		log.Error.Println("读取文件结构错误")
        return nil, err
	}

	answer := resultPre + *result
	return &answer, nil
}