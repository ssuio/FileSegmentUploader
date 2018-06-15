package utils

import (
	"crypto/md5"
	"os"
	"io"
	"log"
	"encoding/hex"
	"strings"
	"github.com/google/uuid"
	"time"
	"fmt"
	"encoding/json"
)

const Conf_PATH = "resource/secret.json"

func GetMd5FromFile(file *os.File) (result string) {
	h := md5.New()
	if _, err := io.Copy(h, file); err != nil {
		log.Fatal(err)
	}
	return hex.EncodeToString(h.Sum(nil))
}


func GetMd5FromBytes(b []byte) (result []byte) {
	h := md5.New()
	h.Write(b)
	return h.Sum(nil)
}

func GetUUID() (uuidStr string) {
	return strings.Replace(uuid.New().String(), "-", "", -1)
}

func GetTimestamp() (timestamp int64) {
	return time.Now().Unix()
}

type Configuration struct {
	APPID     string
	API_KEY   string
	IFT_URL   string
	GetLogInfoUrl string
	UploadResumeUrl string
	ServiceId string
}

func LoadConfiguration() Configuration {
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