package router

import (
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"robot/base"
	"robot/base/log"
	robotdata "robot/base/robot_data"
	robotfile "robot/base/robot_file"
	robothttp "robot/base/robot_http"
	"strconv"
	"time"

	"github.com/gorilla/websocket"
)

const (
	connetPath   = "gateway"
	heartBeatKey = "heartbeat_interval"
)

type Router struct {
	Handlers map[string]EventHandler
	RetryInteval int

	connection *websocket.Conn
	heartBeat float64
	lastestS int
}

func NewDefaultRouter(handlers *[]EventHandler) *Router {
	handlerMap := make(map[string]EventHandler)
	for _, value := range *handlers {
		handler := value
		handlerMap[GetHandlerMapKey(handler.GetOpCode(), handler.GetEventKey())] = handler
	}

	return &Router{
        RetryInteval: 1000,
        Handlers: handlerMap,
    }
}

func (router *Router) Connect() error {
	log.Info.Println("开始建立连接")

	//获取网关Url
	wsUrlInfo, err := router.getGateWayUrl()
	if err != nil {
		log.Error.Println("获取网关Url错误：" + err.Error())
		return err
	}

	//连接到网关
	err = router.connectToGateWay(wsUrlInfo)
	if err != nil {
		log.Error.Println("连接网关Url错误：" + err.Error())
		return err
	}

	//判断是否需要故障恢复
	exist := robotfile.IsFileExist(robotfile.LastestMsgPath)
	if !exist {
		//不需要则进行鉴权连接
		err = router.authConnect()
		if err != nil {
			log.Error.Println("连接鉴权错误: " + err.Error())
			return err
		}

	} else {
		//需要则发送重新连接消息
		err = router.sendResumeInfo()
		if err != nil {
			log.Error.Println("重连错误: " + err.Error())
			return err
		}
	}

	//开启心跳
	go router.startHeartBeat()

	log.Info.Println("连接建立成功")
	return nil
}

func (router *Router) ListenAndRouteEvent() {
	log.Info.Println("开始监听事件")

	conn := router.connection
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Error.Println("远程连接已关闭 " + err.Error())

				Connections <- true
				HeartBeatExit <- true

				break
			} else {
				log.Error.Println("获取事件错误： " + err.Error())
				continue
			}
		}
		msgString := string(msg)
		var eventLoad EventPayLoad

		err = json.Unmarshal(msg, &eventLoad)
		if err != nil {
			log.Error.Println("反序列化event事件错误 " + err.Error())
			continue
		}

		log.Info.Printf("获取到事件：%+v\n", msgString)

		//更新获取到的最新事件
		router.lastestS = eventLoad.S
		robotfile.WriteFile(robotfile.LastestMsgPath, strconv.Itoa(eventLoad.S))

		eventHandler := router.Handlers[GetHandlerMapKey(eventLoad.Op, eventLoad.T)]
		if eventHandler == nil {
			continue
		}
		
		go eventHandler.HandleEvent(&msg)
	}
}

func (router *Router) getGateWayUrl() (*WsUrlInfo, error) {
	url := base.GetUrl() + connetPath
	auth := robotdata.GetRobotData().AuthKey

	resp, err := robothttp.DoGET(&http.Client{}, url, auth)
	if err != nil {
		log.Info.Printf("获取链接错误：" + err.Error())
		return nil, err
	}

	wsUrlInfo := &WsUrlInfo{}
	err = json.Unmarshal([]byte(*resp), wsUrlInfo)
	if(err != nil) {
		log.Error.Println("解析WsURL错误 " + err.Error())
		return nil, err
	}

	return wsUrlInfo, nil
}

func (router *Router) connectToGateWay(wsInfo *WsUrlInfo) error {
	conn, _, err := websocket.DefaultDialer.Dial(wsInfo.Url, nil)
    if err != nil {
		log.Error.Println("建立webSocket连接错误: " + err.Error())
		return err
    }
	router.connection = conn

	_, msg, err := conn.ReadMessage()
	if err != nil {
		log.Error.Println("获取ws连接响应错误： " + err.Error())
		return err
	}

	var heartBeatLoad EventPayLoad
	err = json.Unmarshal(msg, &heartBeatLoad)
	if err != nil {
		log.Error.Println("反序列化心跳信息错误：" + err.Error())
		return err
	}

	heartBeatInfo, ok := (heartBeatLoad.D).(map[string] interface{})
	if !ok {
		errMsg := fmt.Sprintf("心跳信息类型转化错误: %T", heartBeatLoad.D)
		log.Error.Println(errMsg)
		return errors.New(errMsg)
	}

	router.heartBeat = heartBeatInfo[heartBeatKey].(float64)

	return nil
}

func (router *Router) authConnect() error {
	conn := router.connection

	intents := 0
	for _, value := range router.Handlers {
		intents = intents | value.GetEventItent()
	}

	identifyInfo := &EventPayLoad{
		Op: Identify,
		D: AuthInfo{
			Token: robotdata.GetRobotData().AuthKey,
			Intents: intents,
			Shard: []int{0, 1},
		},
	}

	identifyInfoJson, err := json.Marshal(identifyInfo)
	if err != nil {
		log.Error.Println("序列化认证信息错误: " + err.Error())
		return err
	}

	err = conn.WriteMessage(websocket.TextMessage, identifyInfoJson)
	if err != nil {
		log.Error.Println("发送认证信息错误: " + err.Error())
		return err
	}

	_, msg, err := conn.ReadMessage()
	if err != nil {
		log.Error.Println("获取ready响应错误： " + err.Error())
		return err
	}

	var readyLoad ReadyPayLoad
	err = json.Unmarshal(msg, &readyLoad)
	if err != nil {
		log.Error.Println("反序列化ready响应错误： " + err.Error())
		return err
	}
	router.lastestS = readyLoad.S

	err = robotfile.WriteFile(robotfile.LastestSession, readyLoad.D.SessionID)
	if err != nil {
		log.Error.Println("写入session错误 " + err.Error())
		return err
	}

	return nil
}

func (router *Router) sendResumeInfo() error{
	conn := router.connection
	
	seqString, err := robotfile.ReadFile(robotfile.LastestMsgPath)
	if err != nil {
		log.Error.Println("读取文件错误 " + err.Error())
		return err
	}

	seq, err := strconv.Atoi(*seqString)
	if err != nil {
		log.Error.Println("转化失败 " + err.Error())
		return err
	}

	session, err := robotfile.ReadFile(robotfile.LastestSession)
	if err != nil {
		log.Error.Println("读取文件错误 " + err.Error())
		return err
	}

	resumeLoad := &EventPayLoad{
		Op: Resume,
		D:  ResumeInfo{
			Seq: seq,
			SessionID: *session,
			Token: robotdata.GetRobotData().AuthKey,
		},
	}

	resumeJson, err := json.Marshal(resumeLoad)
	log.Info.Println(string(resumeJson))
	if err != nil {
		log.Error.Println("序列错误： " + err.Error())
		return err
	}

	err = conn.WriteMessage(websocket.TextMessage, resumeJson)
	if err != nil {
		log.Error.Println("发送错误: " + err.Error())
		return err
	}

	return nil
}

func (router *Router) startHeartBeat() {
	for {
		select {
		case <- HeartBeatExit :{
			log.Info.Printf("心跳退出")
			break
		}

		default :
			conn := router.connection
			
			heartBeatLoad := &EventPayLoad{
				Op: Heartbeat,
				D:  router.lastestS,
			}

			heartBeatLoadJson, err := json.Marshal(heartBeatLoad)
			if err != nil {
				log.Error.Println("序列化心跳包错误： " + err.Error())
				router.sleepToRetry()
				continue
			}

			err = conn.WriteMessage(websocket.TextMessage, heartBeatLoadJson)
			if err != nil {
				log.Error.Println("发送认证信息错误: " + err.Error())
				router.sleepToRetry()
				continue
			}

			time.Sleep(time.Duration(router.heartBeat) * time.Millisecond)
		}
	}
}

func (router *Router) sleepToRetry () {
	time.Sleep(time.Duration(router.RetryInteval) * time.Millisecond)
}

func GetHandlerMapKey(code OpCode, key string) string {
	return fmt.Sprintf("%d.%s", code, key)
}