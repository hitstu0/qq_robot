package log

import (
	"io"
	"log"
	"os"
)

var (
	Info    *log.Logger // 重要的信息
	Error   *log.Logger // 错误信息
)

func init() {
	errFile, err := os.OpenFile("log/errors.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open error log file: ", err)
	}

	infoFile, err := os.OpenFile("log/info.txt", os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0666)
	if err != nil {
		log.Fatalln("Failed to open info log file: ", err)
	}

	Info = log.New(io.MultiWriter(infoFile, os.Stdout), "Info: ", log.Ldate|log.Ltime|log.Lshortfile)
	Error = log.New(io.MultiWriter(errFile, os.Stderr), "Error: ", log.Ldate|log.Ltime|log.Lshortfile)
}
