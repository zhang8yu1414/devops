package model

type ApiResMessage struct {
	Errors []info `json:"errors"`
}

type info struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}
