package requests

import (
	"encoding/json"
	"fmt"
	"log"
	"os"

	"gopkg.in/resty.v0"
)

const b2ListFiles = "/b2api/v1/b2_list_file_names"

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
	name, err := os.Hostname()
	log.Println("HostName: ", name)
	if err != nil {
		log.Println(err)
	}
	err = listFilesRequest(auth.APIURL, auth.AuthorizationToken, bucketid, "")
	if err != nil {
		return err
	}
	//TODO: handle data paging
	return err
}

//Create the request and executy it.  Adds nexfile if defined
func listFilesRequest(url, token, bucketid, next string) error {
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
		return err
	}
	return nil
}

func listFilesResponseHandler(response *resty.Response) error {
	switch response.StatusCode() {
	case 200:
	case 400:
	case 401:
	case 503:
	}
	return nil
}

//Load the files into json structs, there is a nextfile return a bool telling the calling function to load more data
func loadFiles(data []byte) (string, error) {
	var b2Files B2FilesList
	err := json.Unmarshal(data, &b2Files)
	if err != nil {
		return "", err
	}
	if b2Files.NextFile != "" {
		return b2Files.NextFile, nil
	}
	return "", nil
}
