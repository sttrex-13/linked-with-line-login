package line

import (
	"api/pkg/util"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"

	"github.com/google/uuid"
)

type LINEHandler struct {
	lineClient LINEClientInterface
}

func NewLINEHandler() *LINEHandler {
	return &LINEHandler{
		lineClient: &LINEClient{},
	}
}

func (h *LINEHandler) WebHook(w http.ResponseWriter, r *http.Request) {
	bodyBytes, err := io.ReadAll(r.Body)
	if err != nil {
		log.Println("can not read webhook's request body: ", err)
		return
	}
	fmt.Println("webhook request: ", string(bodyBytes))

	util.WriteJSONResponse(w, http.StatusOK, nil)
}

func (h *LINEHandler) LoginCallback(w http.ResponseWriter, r *http.Request) {
	q := r.URL.Query()
	code := q.Get("code")

	res, err := h.lineClient.GetAccessToken(code)
	if err != nil {
		log.Println("can not get access token: ", err)
		return
	}

	profile, err := h.lineClient.GetProfile(res.AccessToken)
	if err != nil {
		log.Println("can not get profile: ", err)
		return
	}

	fmt.Println("line profile: ", util.ToJsonString(profile))

	err = h.lineClient.PushTextMessage(profile.UserID, "You have logged in successfully")
	if err != nil {
		log.Println("can not push message: ", err)
		return
	}

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
	queryParams += "&scope=profile"

	requestURL := LINERequestAccessApiURL + queryParams

	http.Redirect(w, r, requestURL, http.StatusSeeOther)
}
