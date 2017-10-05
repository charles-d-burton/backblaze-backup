package requests

import (
	"log"

	"github.com/spf13/viper"
)

func CreateBackblazeBucket() error {
	authorization, err := GetAuthorization(viper.GetString("account-id"), viper.GetString("application-key"))
	if err != nil {
		log.Println(err)
	}
	log.Println("Token: ", authorization.AuthorizationToken)

	return err
}
