package robotfile

import (
	"errors"
	"io/ioutil"
	"os"
	"robot/base/log"
)

func IsFileExist(url string) bool {
	_, err := os.Stat(url)
	if err != nil {
		if os.IsNotExist(err) {
			return false;
		}
	}

	return true
}
func ReadFile(url string) (*string, error) {
	_, err := os.Stat(url)
	if err != nil {
		if os.IsNotExist(err) {
			return nil, errors.New("不存在信息" + url)
		}
		return nil, err
	}

	file, err := os.Open(url)
	if err != nil {
		log.Error.Println("打开文件错误：" + err.Error())
		return nil, err
	}
	defer file.Close()

	content, err := ioutil.ReadAll(file)
	if err != nil {
		log.Error.Println("读取错误：" + err.Error())
		return nil, err
	}
	contentString := string(content)

	return &contentString, nil
}

func WriteFile(url string, content string) error {
	// 创建新文件并检查错误
	file, err := os.Create(url)
	if err != nil {
		log.Error.Println("创建文件错误")
		return err
	}
	defer file.Close()

	_, err = file.WriteString(content)
	if err != nil {
		log.Error.Println("写入文件错误")
		return err
	}

	return nil
}

func GetAllFileName(url string) (*string, error) {
	// 读取目录内容
	files, err := ioutil.ReadDir(url)
	if err != nil {
		return nil, err
	}

	result := ""
	for _, file := range files {
		result = result + file.Name() + "\n"
	}

	return &result, nil
}

func MoveFile(url string, newUrl string) error {
	// 判断源文件是否存在
	if _, err := os.Stat(url); os.IsNotExist(err) {
		return errors.New("文件不存在")
	}

	// 将文件移动到目标文件夹
	err := os.Rename(url, newUrl)
	if err != nil {
		return err
	}

	return nil
}

func DeleteFile(url string) error {
	// 检查文件是否存在
	_, err := os.Stat(url)
	if os.IsNotExist(err) {
		return nil
	}

	// 删除文件
	err = os.Remove(url)
	if err != nil {
		return err
	}

	return nil
}
