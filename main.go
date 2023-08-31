package main

import (
	"robot/base/chatgpt"
	"robot/base/gpt_base"
	"robot/base/log"
	robotdata "robot/base/robot_data"
	"robot/base/router"
	"robot/base/router/events"
)
//修复故障恢复
func main() {
    initRobot()
    initGPT()

    handlers := initHandlers()
    router := router.NewDefaultRouter(handlers)

    err := router.Connect()
    if err != nil{
        log.Info.Println(err.Error())
        panic(err)
    }

    router.ListenAndRouteEvent()
}

//初始化程序机器人配置
func initRobot() {
    robotdata.AddNewRobot(robotdata.DefaultRobotKey, 102066650, "pM32zSLslhweVD5r6WEmJgnXk0pNFSOV")
}

//初始化选择的大语言模型
func initGPT() {
    var err error
    gptBase.Client , err = chatgpt.NewDefaultChatGptClient()
    if err != nil {
        log.Error.Printf("初始化语言模型客户端错误")
        panic(err)
    }
}

//配置程序监听的事件
func initHandlers() *[]router.EventHandler{
    handlers := []router.EventHandler {
        &events.MsgEventHandler {
            OpCode: router.Dispatch,
            EventItent: router.PUBLIC_GUILD_MESSAGES,
            EventKey: router.MsgCreateKey,
        },

        &events.HeartBeatEventHandler {
            OpCode: router.HeartbeatACK,
            EventItent: router.PUBLIC_GUILD_MESSAGES,
        },
    }

    return &handlers
}