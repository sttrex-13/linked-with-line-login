package line

import (
	"api/config"
	"api/pkg/util"
	"fmt"
	"io"
	"log"
	"net/http"
)

type LINEHandler struct {
	lineClient LINEClientInterface
}

func NewLINEHandler() *LINEHandler {
	config := config.New()

	return &LINEHandler{
		lineClient: &LINEClient{
			AppURL:                   config.AppURL,
			LINELoginClientID:        config.LINELoginClientID,
			LINELoginChanelSecret:    config.LINELoginChanelSecret,
			LINEMessagingAccessToken: config.LINEMessagingAccessToken,
		},
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
	requestURL := h.lineClient.RequestAuthenticationCode()

	http.Redirect(w, r, requestURL, http.StatusSeeOther)
}
