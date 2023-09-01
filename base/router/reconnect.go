package router

import robotfile "robot/base/robot_file"

func PreReconnect() {
	Wg.Wait()

	robotfile.DeleteFile(robotfile.LastestMsgPath)
	robotfile.DeleteFile(robotfile.LastestSession)
}

func SendReconnectSignal() {
	Connections <- true
	HeartBeatExit <- true
}