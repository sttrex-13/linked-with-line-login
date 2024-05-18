package line

import (
	"api/pkg/util"
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/google/uuid"
)

const (
	AppURL string = "https://polished-pheasant-explicitly.ngrok-free.app"

	LINERequestAccessApiURL string = "https://access.line.me/oauth2/v2.1/authorize"
	LINEMessagingApiURL     string = "https://api.line.me/v2/bot/message"
	LINEOauthApiURL         string = "https://api.line.me/oauth2/v2.1"

	ClientID           string = "2005138101" // LINE Login's Channel ID
	ChannelAccessToken string = "fEd+xwkjCLrKX0Y3kY6i3J1HXpPrv3MN825KF95ES6R1AvARfom5LY2FhgsdkXGpleUzR07LVM8rqMihUK4RVAexkD2hyp0aQ33fDhvTxM0qbQ/0Hw/qwtyO9EOeio0yUHKhKhNs8Ik9TQyeGbZ3JAdB04t89/1O/w1cDnyilFU="
	ChanelSecret       string = "e6f25e7be49afa72300489c4a8241347"
)

type LINEHandler struct {
}

func (h *LINEHandler) WebHook(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("can not read webhook's request body: ", err)
		return
	}
	fmt.Println("webhook request: ", string(bodyBytes))

	var request LINEWebhookRequest
	if err := json.Unmarshal(bodyBytes, &request); err != nil {
		log.Println("can not unmarshal webhook's request body: ", err)
		return
	}

	// reply message attached with login link
	requestURL := LINEMessagingApiURL + "/reply"
	requestBody := map[string]interface{}{
		"replyToken": request.Events[0].ReplyToken, // fixed reply token
		"messages": []map[string]interface{}{
			{
				"type": "text",
				"text": fmt.Sprintf("Here is login link: %s", AppURL+"/api/line/request-login"),
			},
		},
	}
	requestBodyJson, _ := json.Marshal(requestBody)
	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewBuffer(requestBodyJson))
	if err != nil {
		log.Println("can not create request for replying webhook: ", err)
		return
	}
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+ChannelAccessToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("can not request to LINE api: ", err)
		return
	}

	responseBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("can not read reply api response: ", err)
		return
	}
	log.Println("reply webhook response: ", string(responseBodyBytes))

	util.WriteJSONResponse(w, http.StatusOK, nil)
}

func (h *LINEHandler) LoginCallback(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	code := q.Get("code")

	fmt.Println("path is: ", r.URL.Path)

	// Getting an access token
	requestURL := LINEOauthApiURL + "/token"
	requestBody := map[string]interface{}{
		"grant_type":    "authorization_code",
		"code":          code,
		"redirect_uri":  url.QueryEscape(AppURL + "/api/line/login-callback"),
		"client_id":     ClientID,
		"client_secret": ChanelSecret,
	}
	requestBodyJson, _ := json.Marshal(requestBody)
	req, err := http.NewRequest(http.MethodPost, requestURL, bytes.NewBuffer(requestBodyJson))
	if err != nil {
		log.Println("can not create request for getting token: ", err)
		return
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Println("can not request to LINE api: ", err)
		return
	}
	responseBodyBytes, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println("can not read get token api response: ", err)
		return
	}
	log.Println("get token response: ", string(responseBodyBytes))

	util.WriteJSONResponse(w, http.StatusOK, nil)
}

func (h *LINEHandler) RequestLogin(w http.ResponseWriter, r *http.Request) {
	queryParams := "?response_type=code"
	queryParams += "&client_id=" + ClientID
	// Redirect URL: URL-encoded callback URL
	queryParams += "&redirect_uri=" + url.QueryEscape(AppURL+"/api/line/login-callback")
	// State: A unique alphanumeric string
	queryParams += "&state=" + uuid.New().String()
	// Scope: Permissions requested from the user see https://developers.line.biz/en/docs/line-login/integrate-line-login/#scopes
	queryParams += "&scope=profile%20openid"

	requestURL := LINERequestAccessApiURL + queryParams

	fmt.Println("request login: ", requestURL)

	http.Redirect(w, r, requestURL, http.StatusSeeOther)
}
