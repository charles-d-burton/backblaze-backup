package requests

import (
	"crypto/sha1"
	"encoding/hex"
	"fmt"
	"io/ioutil"
	"log"

	"gopkg.in/resty.v0"
)

func (upload *UploadUrl) UploadFile(file string) error {
	fileBytes, _ := ioutil.ReadFile(file)
	log.Println("File length: ", len(fileBytes))
	hasher := sha1.New()
	hasher.Write(fileBytes)
	sha := hex.EncodeToString(hasher.Sum(nil))
	fmt.Print("\n" + sha)
	upLoadFile, err := resty.R().
		SetBody(fileBytes).
		SetContentLength(true).
		SetHeader("Authorization", upload.AuthorizationToken).
		SetHeader("X-Bz-Content-Sha1", sha).
		SetHeader("X-Bz-File-Name", "testfile.txt").
		Post(upload.UploadURL)

	fmt.Printf("\nError: %v", err)
	fmt.Printf("\nResponse Status Code: %v", upLoadFile.StatusCode())
	fmt.Printf("\nResponse Status: %v", upLoadFile.Status())
	fmt.Printf("\nResponse Time: %v", upLoadFile.Time())
	fmt.Printf("\nResponse Recevied At: %v", upLoadFile.ReceivedAt())
	fmt.Printf("\nResponse Body: %v", string(upLoadFile.Body()))
	return err
}
