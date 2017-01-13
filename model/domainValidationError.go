package model

//ProxyError structure contains the error resource definition
type DomainValidationErr struct {
	Status  string `json:"status"`
	Message string `json:"message"`
}
