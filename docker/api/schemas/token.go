package schemas

type Token struct {
	AccessToken  string `json:"accessToken"`
	RefreshToken string `json:"refreshToken"`
	Expire       int64  `json:"exp"`
	TokenType    string `json:"tokenType"`
}

func NewToken(access_token, refreshToke string, exp int64) *Token {
	return &Token{
		AccessToken:  access_token,
		RefreshToken: refreshToke,
		Expire:       exp,
		TokenType:    "bearer",
	}
}
