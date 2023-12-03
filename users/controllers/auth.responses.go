package controllers

type SignInOkRes struct {
	Status      string `json:"status"`
	AccessToken string `json:"access_token"`
}

type ErrResponse struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
