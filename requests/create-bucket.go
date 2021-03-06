package requests

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/resty.v0"

	"github.com/spf13/viper"
)

const b2CreateBucket = "/b2api/v1/b2_create_bucket"

//CreateBackblazeBucket ... Create the bucket based on hostnamet to put file data
func (auth *AuthorizationResponse) CreateBackblazeBucket() error {

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
	fmt.Printf("\nResponse Body: %v", string(request.Body()))

	return err
}
