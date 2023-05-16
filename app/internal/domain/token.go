package domain

type RefreshTokenInput struct {
	RefreshToken string `json:"refresh_token"`
}

type RefreshTokenOutput struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}
