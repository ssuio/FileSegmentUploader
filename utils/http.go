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
)

var Conf Configuration

func init() {
	Conf = LoadConfiguration()
}

func EncodeUrl(urlStr string) (string) {
	return strings.ToLower(url.QueryEscape(urlStr))
}

func RequestPost(url string, data string) ([]byte) {
	return request("POST", url, data)
}

func request(method string, url string, data string) ([]byte) {
	req, err := http.NewRequest(method, url, bytes.NewBufferString(data))
	headerStr := generateHeader(url, method, bytes.NewBufferString(data).Bytes())
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
		return body
	}
	return nil
}

func generateHeader(url string, method string, content []byte) (string) {
	timestamp := string(GetTimestamp())
	nonce := GetUUID()
	url = EncodeUrl(url)

	contentMd5 := GetMd5FromBytes(content)
	contentRes := base64.StdEncoding.EncodeToString(contentMd5)

	sigRaw := Conf.APPID + method + url + timestamp + nonce + contentRes

	secretKeyByteArr, err := base64.StdEncoding.DecodeString(Conf.API_KEY)
	if err != nil {
		log.Fatal("decode API key failed.", err)
	}
	mac := hmac.New(sha256.New, secretKeyByteArr)
	mac.Write([]byte(sigRaw))
	sigRes := base64.StdEncoding.EncodeToString(mac.Sum(nil))

	authStr := Conf.APPID + ":" + sigRes + ":" + nonce + ":" + timestamp
	return "icss " + authStr
}
