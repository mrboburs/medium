package model

type ResponseSuccess struct {
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

type RSUpdateAccount struct {
	Message string      `json:"message"`
	Id      interface{} `json:"id"`
}
