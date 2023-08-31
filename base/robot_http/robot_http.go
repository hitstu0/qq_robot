package robothttp

import (
	"errors"
	"io/ioutil"
	"net/http"
	"robot/base/log"
	"strings"
)

const (
	post = "POST"
	get  = "GET"
	contentType = "Content-Type"
	authorization = "Authorization"
)

const (
	applicationJson = "application/json"
)

func DoPOST(client *http.Client, url *string, content *string, auth *string) (*string, error) {
	req, _ := http.NewRequest(post, *url, strings.NewReader(*content))
	req.Header.Set(contentType, applicationJson)
	req.Header.Set(authorization, *auth)

	resp, err := client.Do(req)
	if err != nil {
		log.Error.Println("获取响应错误： " + err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	respByte, _ := ioutil.ReadAll(resp.Body)
	respString := string(respByte)

	if resp.StatusCode < 200 || resp.StatusCode > 300 {
		log.Error.Printf("响应码不为2XX: %d\n", resp.StatusCode)
		return nil, errors.New(respString)
	}

	return &respString, nil
}

func DoGET(client *http.Client, url string, auth string) (*string, error) {
	req, _ := http.NewRequest(get, url, nil)
	req.Header.Set(authorization, auth)

	resp, err := client.Do(req)
	if err != nil {
		log.Error.Println("获取响应错误： " + err.Error())
		return nil, err
	}
	defer resp.Body.Close()

	respByte, _ := ioutil.ReadAll(resp.Body)
	respString := string(respByte)
	log.Info.Println("接收响应：" + respString)

	if resp.StatusCode < 200 || resp.StatusCode > 300 {
		log.Error.Printf("响应码不为2XX: %d\n", resp.StatusCode)
		return nil, errors.New(respString)
	}

	return &respString, nil
}