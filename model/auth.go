package model

type ResponseSignUp struct {
	Id    int    `json:"id" `
	Token string `json:"token"`
}

type SignInInput struct {
	Email    string `json:"email"  default:"babdusalom72@gmail.com"`
	Password string `json:"password" default:"1996"`
}

type VerificationCode struct {
	Code string `json:"code"`
}
