package adventure

import (
	"robot/base/log"
	robotfile "robot/base/robot_file"
)

const (
	getListKey  	 = "英雄列表"

	resultPre = "支持的英雄如下:\n"
)

func NewGetHeroListHandler() *GetHeroListHandler {
	return &GetHeroListHandler{
		Key: getListKey,
	}
}

type GetHeroListHandler struct {
	Key string
}

func (handler *GetHeroListHandler) ParseParm(parm *string) error {
	return nil
}

func (handler *GetHeroListHandler) HandlerRequest() (*string, error) {
    result, err := robotfile.GetAllFileName(robotfile.HeroInfoPath)
	if err != nil {
		log.Error.Println("读取文件结构错误")
        return nil, err
	}

	answer := resultPre + *result
	return &answer, nil
}