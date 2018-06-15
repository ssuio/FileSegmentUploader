package main

import (
	"os"
	"fmt"
	"src/utils"
	"src/api"
)


const LOG_PATH = "resource/test.txt"

var conf utils.Configuration
var getLogInfoObj GetLogInfoRequest

func main() {
	logs := api.GetLogInfo(utils.ToJsonStr(getLogInfoObj))
	fmt.Println(logs)
}

func init() {
	conf = utils.LoadConfiguration()
	getLogInfoObj = GetLogInfoRequest{
		ServiceId: conf.ServiceId,
	}
}

type GetLogInfoRequest struct {
	ServiceId string `json:"Serviceid"`
}

func clientDo() {
	file, _ := os.Open(LOG_PATH)
	defer file.Close()
	//md5Str := getMd5FromFile(file)
}