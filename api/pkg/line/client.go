package line

import (
	"bytes"
	"encoding/json"
	"log"
	"net/http"
	"net/url"
	"strings"

	"github.com/google/uuid"
)

const (
	LINERequestAccessApiURL string = "https://access.line.me/oauth2/v2.1/authorize"
	LINEMessagingApiURL     string = "https://api.line.me/v2/bot/message"
	LINEOauthApiURL         string = "https://api.line.me/oauth2/v2.1"
	LINEAccountApiURL       string = "https://api.line.me/v2/"
)

type LINEClientInterface interface {
	RequestAuthenticationCode() (url string)
	GetAccessToken(authenticationCode string) (resp GetAccessTokenResponse, err error)
	GetProfile(accessToken string) (resp GetProfileResponse, err error)
	PushTextMessage(userId string, message string) (err error)
}

type LINEClient struct {
	AppURL                   string
	LINELoginClientID        string
	LINELoginChanelSecret    string
	LINEMessagingAccessToken string
}

func (c *LINEClient) RequestAuthenticationCode() string {
	queryParams := "?response_type=code"
	queryParams += "&client_id=" + c.LINELoginClientID
	// Redirect URL: URL-encoded callback URL
	queryParams += "&redirect_uri=" + url.QueryEscape(c.AppURL+"/api/line/login-callback")
	// State: A unique alphanumeric string
	queryParams += "&state=" + uuid.New().String()
	// Scope: Permissions requested from the user see https://developers.line.biz/en/docs/line-login/integrate-line-login/#scopes
	queryParams += "&scope=profile"

	requestURL := LINERequestAccessApiURL + queryParams
	return requestURL
}

func (c *LINEClient) PushTextMessage(userId string, message string) (err error) {
	var (
		requestURL  = LINEMessagingApiURL + "/push"
		requestBody = map[string]interface{}{
			"to": userId,
			"messages": []map[string]interface{}{
				{
					"type": "text",
					"text": message,
				},
			},
		}
	)

	jsonRequestBody, _ := json.Marshal(requestBody)
	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewBuffer(jsonRequestBody))
	if err != nil {
		log.Println("can not create request for replying webhook: ", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+c.LINEMessagingAccessToken)

	client := &http.Client{}
	_, err = client.Do(req)
	if err != nil {
		log.Println("can not request to LINE api: ", err)
		return
	}

	return
}

type GetAccessTokenResponse struct {
	AccessToken  string `json:"access_token"`
	TokenType    string `json:"token_type"`
	ExpiresIn    int    `json:"expires_in"`
	Scope        string `json:"scope"`
	RefreshToken string `json:"refresh_token"`
}

func (c *LINEClient) GetAccessToken(authenticationCode string) (resp GetAccessTokenResponse, err error) {
	var (
		requestURL  = LINEOauthApiURL + "/token"
		requestBody = url.Values{
			"grant_type":    {"authorization_code"},
			"code":          {authenticationCode},
			"redirect_uri":  {c.AppURL + "/api/line/login-callback"},
			"client_id":     {c.LINELoginClientID},
			"client_secret": {c.LINELoginChanelSecret},
		}
		encodedRequestBody = requestBody.Encode()
	)

	// new POST request
	req, err := http.NewRequest(http.MethodPost, requestURL, strings.NewReader(encodedRequestBody))
	if err != nil {
		log.Println("can not create request for getting token: ", err)
		return
	}

	// set header
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	// starting request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Println("can not request to LINE api: ", err)
		return
	}

	// read response
	err = json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		log.Println("can not read get token api response: ", err)
		return
	}

	return resp, nil
}

type GetProfileResponse struct {
	UserID        string `json:"userId"`
	DisplayName   string `json:"displayName"`
	StatusMessage string `json:"statusMessage"`
	PictureUrl    string `json:"pictureUrl"`
}

func (c *LINEClient) GetProfile(accessToken string) (resp GetProfileResponse, err error) {
	var (
		requestURL = LINEAccountApiURL + "/profile"
	)

	// new POST request
	req, err := http.NewRequest(http.MethodGet, requestURL, nil)
	if err != nil {
		log.Println("can not create request for getting token: ", err)
		return
	}

	// set header
	req.Header.Set("Authorization", "Bearer "+accessToken)

	// starting request
	client := &http.Client{}
	res, err := client.Do(req)
	if err != nil {
		log.Println("can not request to LINE api: ", err)
		return
	}

	// read response
	err = json.NewDecoder(res.Body).Decode(&resp)
	if err != nil {
		log.Println("can not read get token api response: ", err)
		return
	}

	return resp, nil
}
