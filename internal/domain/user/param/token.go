package param

type Tokens struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

type RefreshAccessTokenRequest struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshAccessTokenResponse struct {
	AccessToken string `json:"access_token"`
}
