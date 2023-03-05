package schemas

type SignUpInput struct {
	Username string `json:"username" validate:"required,min=3,max=32"`
	Email    string `json:"email" validate:"required,email,min=6,max=32"`
	Password string `json:"password" validate:"required,min=6,max=32"`
}

type SignInInput struct {
	Email        string `json:"email" validate:"required_if=GrantType password"`
	Password     string `json:"password" validate:"required_if=GrantType password"`
	GrantType    string `json:"grantType" validate:"oneof=password refreshToken"`
	RefreshToken string `json:"refreshToken" validate:"required_if=GrantType refreshToken"`
}
