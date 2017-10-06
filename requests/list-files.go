package requests

import (
	"fmt"
	"log"
	"os"

	"github.com/spf13/viper"
	"gopkg.in/resty.v0"
)

const b2ListFiles = "/b2api/v1/b2_list_file_names"

type B2FilesList struct {
	B2Files []B2File `json:"files"`
}

type B2File struct {
	Action          string `json:"action"`
	ContentLength   int    `json:"contentLength"`
	FileId          string `json:"fileId"`
	FileName        string `json:"fileName"`
	Size            int    `json:"size"`
	UploadTimestamp int64  `json:"uploadTimestamp"`
}

func (auth *AuthorizationResponse) ListAllFiles(bucketid string) {
	name, err := os.Hostname()
	log.Println("HostName: ", name)
	if err != nil {
		log.Println(err)
	}
	request, err := resty.R().
		SetQueryParams(map[string]string{
		"accountId":  viper.GetString("account-id"),
		"bucketName": name,
		"bucketType": "allPrivate",
	}).
		SetHeader("Accept", "application/json").
		SetHeader("Authorization", auth.AuthorizationToken).
		Get(auth.APIURL + b2ListFiles)
	fmt.Printf("\nResponse Body: %v", string(request.Body()))

	return err
}
