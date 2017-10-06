package requests

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"gopkg.in/resty.v0"
)

//go:generate protoc --go_out=. files.proto

const (
	b2ListFiles   = "/b2api/v1/b2_list_file_names"
	b2GetFileInfo = "/b2api/v1/b2_get_file_info"
)

var (
	b2FilesChan chan B2FilesList
)

type B2FilesList struct {
	B2Files  []B2File `json:"files"`
	NextFile string   `json:"nextFileName"`
}

type B2File struct {
	Action          string `json:"action"`
	ContentLength   int    `json:"contentLength"`
	FileId          string `json:"fileId"`
	FileName        string `json:"fileName"`
	Size            int    `json:"size"`
	UploadTimestamp int64  `json:"uploadTimestamp"`
}

func (auth *AuthorizationResponse) ListAllFiles(bucketid string) error {
	b2FilesChan = make(chan B2FilesList, 100)
	go auth.processFileLists(b2FilesChan)
	name, err := os.Hostname()
	log.Println("HostName: ", name)
	if err != nil {
		log.Println(err)
	}
	nextFile, err := listFilesRequest(auth.APIURL, auth.AuthorizationToken, bucketid, "")
	if err != nil {
		return err
	}
	if nextFile != "" {
		for {
			next, err := listFilesRequest(auth.APIURL, auth.AuthorizationToken, bucketid, nextFile)
			if err != nil || next == "" {
				break
			}
		}
	}
	close(b2FilesChan)
	//TODO: handle data paging
	return err
}

//Create the request and executy it.  Adds nexfile if defined
func listFilesRequest(url, token, bucketid, next string) (string, error) {
	params := make(map[string]string)
	params["bucketId"] = bucketid
	if next != "" {
		params["startFileName"] = next
	}

	request, err := resty.R().
		SetQueryParams(params).
		SetHeader("Accept", "application/json").
		SetHeader("Authorization", token).
		Get(url + b2ListFiles)
	fmt.Printf("\nResponse Body: %v", string(request.Body()))
	err = listFilesResponseHandler(request)
	if err != nil {
		return "", err
	}
	nextFile, err := loadFiles(request.Body())
	if err != nil {
		return "", err
	}
	return nextFile, nil
}

func listFilesResponseHandler(response *resty.Response) error {
	switch response.StatusCode() {
	case 200:
		log.Println("Doot doot!  ")
		return nil
	case 400:
		log.Println("Bad Request: ", string(response.Body()))
		return errors.New(string(response.Body()))
	case 401:
		log.Println("Bad Token: ", string(response.Body()))
		return errors.New(string(response.Body()))
	case 503:
		log.Println("Timeout: ", string(response.Body()))
		return errors.New(string(response.Body()))
	}
	return nil
}

//Load the files into json structs sends to processor
//there is a nextfile return a string telling the calling function to load more data
func loadFiles(data []byte) (string, error) {
	var b2Files B2FilesList
	err := json.Unmarshal(data, &b2Files)
	if err != nil {
		return "", err
	}
	b2FilesChan <- b2Files
	if b2Files.NextFile != "" {
		return b2Files.NextFile, nil
	}
	return "", nil
}

type b2FileMeta struct {
	BucketID   string `json:"bucketId"`
	FileLength int64  `json:"contentLength"`
	Sha1       string `json:"contentSha1"`
	FileID     string `json:"fileId"`
	FileName   string `json:"fileName"`
}

//Digest the paged arrays, request the File object from Backblaze and then update bolt with the metadata
func (auth *AuthorizationResponse) processFileLists(ch <-chan B2FilesList) {
	for files := range ch {
		for _, file := range files.B2Files {
			log.Println("Processing file with name: ", file.FileName)
			params := make(map[string]string)
			params["fileId"] = file.FileId

			request, err := resty.R().
				SetQueryParams(params).
				SetHeader("Accept", "application/json").
				SetHeader("Authorization", auth.AuthorizationToken).
				Get(auth.APIURL + b2GetFileInfo)
			if err != nil {
				log.Println(err)
				continue
			}
			if request.StatusCode() == 200 {
				fmt.Printf("\nResponse Body: %v", string(request.Body()))
				var meta b2FileMeta
				err = json.Unmarshal(request.Body(), &meta)
			}
		}
	}
	log.Println("Channel closed")
}
