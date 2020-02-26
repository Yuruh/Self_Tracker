package spotify

type Album struct {
	Artists []Artist `json:"artists"`
	Name string      `json:"name"`
	Uri string       `json:"uri"`
}

type Artist struct {
	Name string `json:"name"`
	Id string `json:"id"`
	Uri string `json:"uri"`
}

type Track struct {
	Name    string   `json:"name"`
	Id      string   `json:"id"`
	Uri     string   `json:"uri"`
	Artists []Artist `json:"artists"`
	Album   Album    `json:"album"`
}

type TokenResponse struct {
	AccessToken string `json:"access_token"`
	TokenType string `json:"token_type"`
	ExpiresIn uint32 `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

type RegisterTokenRequest struct {
	Code string `json:"code"`
	State string `json:"state"`
}