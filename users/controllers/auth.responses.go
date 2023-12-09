package controllers

type SignInOkRes struct {
	Status       string `json:"status"`
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type ErrResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
