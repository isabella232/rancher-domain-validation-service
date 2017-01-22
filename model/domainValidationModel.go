package model

//DomainValidationErr structure contains the error resource definition
type DomainValidationErr struct {
	Type    string `json:"type"`
	Status  string `json:"code"`
	Code    string `json:"status"`
	Message string `json:"message"`
}

//CreateDatabase for initialize the DomainValidation
func CreateDatabase() {

}
