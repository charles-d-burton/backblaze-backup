package requests

import (
	"fmt"
	"log"
	"os"

	"gopkg.in/resty.v0"

	"github.com/spf13/viper"
)

const b2CreateBucket = "/b2api/v1/b2_create_bucket"

func CreateBackblazeBucket() error {
	authorization, err := GetAuthorization(viper.GetString("account-id"), viper.GetString("application-key"))
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
		SetHeader("Authorization", authorization.AuthorizationToken).
		Get(authorization.APIURL + b2CreateBucket)
	fmt.Printf("\nResponse Body: %v", string(request.Body()))

	return err
}
