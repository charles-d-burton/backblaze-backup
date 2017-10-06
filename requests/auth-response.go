package requests

import (
	"encoding/base64"
	"encoding/json"
	"fmt"

	"gopkg.in/resty.v0"
)

type AuthorizationResponse struct {
	AccountID          string `json:"accountId"`
	APIURL             string `json:"apiUrl"`
	AuthorizationToken string `json:"authorizationToken"`
	DownloadURL        string `json:"downloadUrl"`
	MininumPartSize    int    `json:"minimumPartSize"`
}

func GetAuthorization(id, token string) (AuthorizationResponse, error) {
	keyString := id + ":" + token
	sEnc := "Basic " + base64.StdEncoding.EncodeToString([]byte(keyString))
	getToken, err := resty.R().
		SetHeader("Authorization", sEnc).
		Get("https://api.backblazeb2.com/b2api/v1/b2_authorize_account")

	fmt.Printf("\nError: %v", err)
	fmt.Printf("\nResponse Status Code: %v", getToken.StatusCode())
	fmt.Printf("\nResponse Status: %v", getToken.Status())
	fmt.Printf("\nResponse Time: %v", getToken.Time())
	fmt.Printf("\nResponse Recevied At: %v", getToken.ReceivedAt())
	fmt.Printf("\nResponse Body: %v", string(getToken.Body()))

	var authResponse AuthorizationResponse
	err = json.Unmarshal(getToken.Body(), &authResponse)
	return authResponse, err
}
