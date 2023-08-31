package adventure

import (
	"bytes"
	"errors"
	"fmt"
	gptBase "robot/base/gpt_base"
	"robot/base/log"
	robotfile "robot/base/robot_file"
	"strings"
	"time"
)

const (
	adventureKey = "讲述故事"
	expiration int64 = 60 * 1000

	promotBegin = "我想让你扮演讲故事的角色。您将根据我提供的信息，说出引人入胜、富有想象力和吸引观众的有趣冒险故事。接下来我将提供参与冒险的角色名（我用[]标出,你的输出不用标）和对应的角色介绍（我用{}标出，你的输出不用标）。"
	promotBackProvide = "冒险的背景是："
	promotBackNotProvide = "冒险的背景你自由发挥"
	promotEnd = "请直接输出故事。"
)

func NewAdventureHandler() *AdventureHandler {
	return &AdventureHandler{
		Key: adventureKey,
		heroIntroCache: make(map[string]*cacheInfoStruct),
	}
}

type AdventureHandler struct {
	Key      string
	heroList []string
	storyKey string
	heroIntroCache map[string]*cacheInfoStruct
}

type cacheInfoStruct struct {
	introduction 	string
	lastTime 		int64
}

func (handler *AdventureHandler) ParseParm(parm *string) error {
	//清空之前的数据
	handler.init()

	datas := strings.Split(*parm, "\n")
	if len(datas) < 2 {
		return errors.New("用户输入至少为2行")
	}

	heros := datas[1]
	heroParse := strings.Split(heros, " ")
	if len(heroParse) < 1 {
		return errors.New("至少输入一个英雄名")
	}

	for key, value := range heroParse {
		heroParse[key] = strings.TrimSpace(value)
	}

	handler.heroList = heroParse
	if len(datas) >= 3 {
		handler.storyKey = datas[2]
	}

	return nil
}

func (handler *AdventureHandler) HandlerRequest() (*string, error) {
	heroInfoMap := make(map[string]*string)

	//检查缓存中是否有，如果有则添加，没有则从文件中读取并添加到缓存中
	for _, name := range handler.heroList {
		cacheInfo, ok := handler.getInfoFromCache(name)
		if ok {
			heroInfoMap[name] = cacheInfo
			continue
		}

		fileInfo, err := handler.getInfoFromFile(name)
		if err != nil {
			log.Error.Println("读取英雄信息错误")
			return nil, errors.New("没有该英雄的信息：" + name)
		}

		heroInfoMap[name] = fileInfo
		handler.heroIntroCache[name] = &cacheInfoStruct{
			introduction: *fileInfo,
			lastTime: time.Now().UnixMilli(),
		}
	}

	promot := handler.getPromot(heroInfoMap)
	client := gptBase.Client
	result, err := client.SendMessage([]*gptBase.GptMsg{
		&gptBase.GptMsg{
			Role: gptBase.UserRole,
			Content: *promot,
		},
	})
	if err != nil {
		log.Error.Println("和大语言模型通信错误：" + err.Error())
		return nil, err
	}

	return result, nil
}

func (handler *AdventureHandler) getPromot(infoMap map[string]*string) *string {
	buf := bytes.NewBufferString(promotBegin)
	for key, value := range infoMap {
		buf.WriteString(fmt.Sprintf("[%s] {%s}。", key, *value))
	}

	if len(handler.storyKey) == 0 {
		buf.WriteString(promotBackNotProvide)
	} else {
		buf.WriteString(promotBackProvide + handler.storyKey + "。")
	}

	buf.WriteString(promotEnd)
	
	result := buf.String()
	return &result
}

func (handler *AdventureHandler) getInfoFromCache(name string) (*string, bool) {
	cacheInfo, ok := handler.heroIntroCache[name]
	if !ok {
		return nil, false
	}

	if time.Now().UnixMilli() - cacheInfo.lastTime > expiration {
		return nil, false
	}

	return &cacheInfo.introduction, true
}

func (handler *AdventureHandler) getInfoFromFile(name string) (*string, error) {
	url := robotfile.HeroInfoPath + name

	content, err := robotfile.ReadFile(url)
	if err != nil {
		log.Error.Println("读取文件错误：" + url)
		return nil, err
	}

	return content, nil
}

func (handler *AdventureHandler) init() {
	handler.storyKey = ""
	handler.heroList = []string{}
}