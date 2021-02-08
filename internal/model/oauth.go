package model

type (
	OauthRequest struct {
		Code string `json:"code"`
	}

	OauthResponse struct {
		UserId   uint   `json:"user_id"`
		Username string `json:"username"`
	}
)
