package robotdata

import "fmt"

var (
	DefaultRobotKey = "default"

	robotMap = make(map[string]*RobotData)

	curRobot string
)

func init() {
	curRobot = DefaultRobotKey
}

type RobotData struct {
	AppId    int
	AppToken string
	AuthKey  string
}

func AddNewRobot(key string, appId int, appToken string) {
	data := &RobotData{
		AppId:    appId,
		AppToken: appToken,
		AuthKey:  fmt.Sprintf("Bot %d.%s", appId, appToken),
	}

	robotMap[key] = data
}

func GetRobotData() *RobotData {
	return robotMap[curRobot]
}