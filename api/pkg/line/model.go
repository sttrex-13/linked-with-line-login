package line

type LINEWebhookRequest struct {
	Destination string                    `json:"destination"`
	Events      []LINEWebhookEventRequest `json:"events"`
}

type LINEWebhookEventRequest struct {
	ReplyToken string `json:"replyToken"`
	// i am focus only field that i need..
}
