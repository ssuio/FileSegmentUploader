package api

import (
	"src/model"
	"encoding/json"
	"src/utils"
)

func GetLogInfo(reqStr string) ([]model.LogInfo) {
	retBytes := utils.RequestPost("https://"+utils.Conf.IFT_URL+utils.Conf.GetLogInfoUrl, reqStr)
	var retData map[string][]model.LogInfo
	json.Unmarshal(retBytes, &retData)
	for key, value := range retData {
		if key == "Items" {
			return value
		}
	}
	return nil
}

func uploadAttachmentResume() {

}
