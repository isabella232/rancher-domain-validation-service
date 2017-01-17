package model

//DomainValidationErr structure contains the error resource definition
type DomainValidationErr struct {
	Type    string `json:"type"`
	Status  string `json:"status"`
	Code    string `json:"code"`
	Message string `json:"message"`
}

//CreateDatabase for initialize the DomainValidation
func CreateDatabase() {

}
