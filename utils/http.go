package utils

import (
	"strings"
	"net/url"
	"encoding/base64"
	"log"
	"crypto/hmac"
	"crypto/sha256"
	"net/http"
	"bytes"
	"io/ioutil"
	"encoding/json"
	"src/model"
)

var conf Configuration

func init(){
	conf = LoadConfiguration()
}



func EncodeUrl(urlStr string) (string) {
	return strings.ToLower(url.QueryEscape(urlStr))
}

func ExecutePost(url string, data string)([]model.LogInfo) {
	req, err := http.NewRequest("POST", url, bytes.NewBufferString(data))
	headerStr := generateHeader(url, "POST", bytes.NewBufferString(data).Bytes())
	req.Header.Set("Authorization", headerStr)
	req.Header.Set("content-type", "application/json")
	if err != nil {
		panic(err)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	if resp.StatusCode >= 200 && resp.StatusCode <= 299 {
		body, _ := ioutil.ReadAll(resp.Body)
		var retData map[string][]model.LogInfo
		json.Unmarshal([]byte(body), &retData)
		for key, value:= range retData{
			if key == "Items"{
				return value
			}
		}
	}
	return nil
}

func generateHeader(url string, method string, content []byte) (string) {
	timestamp := string(GetTimestamp())
	nonce := GetUUID()
	url = EncodeUrl(url)

	contentMd5 := GetMd5FromBytes(content)
	contentRes := base64.StdEncoding.EncodeToString(contentMd5)

	sigRaw := conf.APPID + method + url + timestamp + nonce + contentRes

	secretKeyByteArr, err := base64.StdEncoding.DecodeString(conf.API_KEY)
	if err != nil {
		log.Fatal("decode API key failed.", err)
	}
	mac := hmac.New(sha256.New, secretKeyByteArr)
	mac.Write([]byte(sigRaw))
	sigRes := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	authStr := conf.APPID + ":" + sigRes + ":" + nonce + ":" + timestamp
	return "icss " + authStr
}