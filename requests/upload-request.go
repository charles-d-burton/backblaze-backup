package requests

import (
	"encoding/json"
	"fmt"

	"gopkg.in/resty.v0"
)

const b2UploadURL = "/b2api/v1/b2_get_upload_url"

type UploadUrl struct {
	BucketID           string `json:"bucketId"`
	UploadURL          string `json:"uploadUrl"`
	AuthorizationToken string `json:"authorizationToken"`
}

type BucketID struct {
	BucketId string `json:"bucketId"`
}

/*
 * GetFileUploadRequest ...Retrieve a request containing a unique download URL, butcket, and token
 */
func (auth *AuthorizationResponse) GetFileUploadRequest(bucketid string) (UploadUrl, error) {
	bucketID := BucketID{BucketId: bucketid}
	bucketBytes, err := json.Marshal(&bucketID)
	uploadUrl, err := resty.R().
		SetBody(bucketBytes).
		SetHeader("Accept", "application/json").
		SetHeader("Authorization", auth.AuthorizationToken).
		Post(auth.APIURL + b2UploadURL)

	fmt.Printf("\nError: %v", err)
	fmt.Printf("\nStatus Code: %v", uploadUrl.StatusCode())
	fmt.Printf("\nStatus: %v", uploadUrl.Status())
	fmt.Printf("\nTime: %v", uploadUrl.Time())
	fmt.Printf("\nRecevied At: %v", uploadUrl.ReceivedAt())
	fmt.Printf("\nBody: %v", string(uploadUrl.Body()))

	var upload UploadUrl
	err = json.Unmarshal(uploadUrl.Body(), &upload)
	return upload, err
}
