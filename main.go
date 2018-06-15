package main

import (
	"os"
	"fmt"
	"encoding/json"
	"src/utils"
)


const LOG_PATH = "resource/test.txt"

var conf utils.Configuration
var getLogInfoObj GetLogInfoRequest

func main() {
	logs := utils.ExecutePost("https://"+conf.IFT_URL+conf.GetLogInfoUrl, toJsonStr(getLogInfoObj))
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

func toJsonStr(o GetLogInfoRequest) (string) {
	s, _ := json.Marshal(o)
	return string(s)
}

func clientDo() {
	file, _ := os.Open(LOG_PATH)
	defer file.Close()
	//md5Str := getMd5FromFile(file)
}