package main

import (
	"os"
	"encoding/hex"
	"crypto/md5"
	"io"
	"log"
	"fmt"
	"github.com/google/uuid"
	"strings"
	"time"
	"encoding/base64"
	"crypto/hmac"
	"crypto/sha256"
	"net/http"
	"bytes"
	"io/ioutil"
	"net/url"
	"encoding/json"
)

var conf Configuration
var getLogInfoObj GetLogInfoRequest

const LOG_PATH = "resource/test.txt"
const Conf_PATH = "resource/secret.json"

type Configuration struct {
	APPID     string
	API_KEY   string
	IFT_URL   string
	GetLogInfoUrl string
	ServiceId string
}

type GetLogInfoRequest struct {
	ServiceId string `json:"Serviceid"`
}

func init(){
	conf = loadConfiguration()
	getLogInfoObj = GetLogInfoRequest{
		ServiceId: conf.ServiceId,
	}
}

func main() {
	executePost("https://"+conf.IFT_URL+conf.GetLogInfoUrl, toJsonStr(getLogInfoObj))
}

func loadConfiguration() Configuration {
	var config Configuration
	configFile, err := os.Open(Conf_PATH)
	defer configFile.Close()
	if err != nil {
		fmt.Println(err.Error())
	}
	jsonParser := json.NewDecoder(configFile)
	jsonParser.Decode(&config)
	return config
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

func executePost(url string, data string) {
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

	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ := ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}

func generateHeader(url string, method string, content []byte) (string) {
	timestamp := string(getTimestamp())
	nonce := getUUID()
	url = encodeUrl(url)

	contentMd5 := getMd5FromBytes(content)
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

func getMd5FromFile(file *os.File) (result string) {
	h := md5.New()
	if _, err := io.Copy(h, file); err != nil {
		log.Fatal(err)
	}
	return hex.EncodeToString(h.Sum(nil))
}

func getMd5FromBytes(b []byte) (result []byte) {
	h := md5.New()
	h.Write(b)
	return h.Sum(nil)
}

func getUUID() (uuidStr string) {
	return strings.Replace(uuid.New().String(), "-", "", -1)
}

func getTimestamp() (timestamp int64) {
	return time.Now().Unix()
}

func encodeUrl(urlStr string) (string) {
	return strings.ToLower(url.QueryEscape(urlStr))
}
