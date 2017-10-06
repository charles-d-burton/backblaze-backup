package requests

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"os"

	"gopkg.in/resty.v0"

	"github.com/spf13/viper"
)

const (
	b2CreateBucket = "/b2api/v1/b2_create_bucket"
	b2ListBucket   = "/b2api/v1/b2_list_buckets"
)

type B2Buckets struct {
	B2Buckets []B2Bucket `json:"buckets"`
}

type B2Bucket struct {
	AccountID      string          `json:"accountId"`
	BucketID       string          `json:"bucketId"`
	BucketInfo     json.RawMessage `json:"bucketInfo"`
	BucketName     string          `json:"bucketName"`
	BucketType     string          `json:"bucketType"`
	LifeCycleRules json.RawMessage `json:"lifecycleRules"`
}

//TODO:  Add in handler to check response status code

//CreateBackblazeBucket ... Create the bucket based on hostnamet to put file data, if bucket exists return bucketid
func (auth *AuthorizationResponse) CreateBackblazeBucket() (string, error) {

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
		Get(auth.APIURL + b2CreateBucket)
	fmt.Printf("\nResponse Status Code: %v", request.StatusCode())
	fmt.Printf("\nResponse Body: %v", string(request.Body()))
	if request.StatusCode() != 200 {
		request, err = resty.R().
			SetQueryParams(map[string]string{
			"accountId": viper.GetString("account-id"),
		}).
			SetHeader("Accept", "application/json").
			SetHeader("Authorization", auth.AuthorizationToken).
			Get(auth.APIURL + b2ListBucket)
		if request.StatusCode() == 200 {
			var buckets B2Buckets
			err := json.Unmarshal(request.Body(), &buckets)
			if err != nil {
				return "", err
			}
			for _, bucket := range buckets.B2Buckets {
				if bucket.BucketName == name {
					return bucket.BucketID, nil
				}
			}
		}
		return "", errors.New("Bucket not found")
	} else {
		var bucket B2Bucket
		err := json.Unmarshal(request.Body(), &bucket)
		if err != nil {
			return "", err
		}
		return bucket.BucketID, nil
	}
}
