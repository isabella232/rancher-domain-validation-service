package service

import (
	"database/sql"
	"encoding/json"
	"net/http"
	"strconv"

	log "github.com/Sirupsen/logrus"
	"github.com/rancher/api-filter-proxy/model"
)

//RequestData is for the JSON output
type RequestData struct {
	Headers map[string][]string    `json:"headers,omitempty"`
	Body    map[string]interface{} `json:"body,omitempty"`
}

//AuthorizeData is for the JSON output
type AuthorizeData struct {
	Message string `json:"message,omitempty"`
}

//MessageData is for the JSON output
type MessageData struct {
	Data []interface{} `json:"data,omitempty"`
}

//GetDomains return the list of domain
func GetDomains(w http.ResponseWriter, r *http.Request) {
	db, err := sql.Open("mysql", "user:password@/dbname")
	checkErr(err)
}

//AddDomain into the domain list
func AddDomain(w http.ResponseWriter, r *http.Request) {
}

//ActivateDomain is for activating the domain
func ActivateDomain(w http.ResponseWriter, r *http.Request) {
}

//ActivateDomain is for activating the domain
func DeleteDomain(w http.ResponseWriter, r *http.Request) {
}

//ValidateDomian filter
func ValidateDomian(w http.ResponseWriter, r *http.Request) {
}

//ReturnHTTPError handles sending out CatalogError response
func ReturnHTTPError(w http.ResponseWriter, r *http.Request, httpStatus int, errorMessage string) {
	svcError := model.ProxyError{
		Status:  strconv.Itoa(httpStatus),
		Message: errorMessage,
	}
	writeError(w, svcError)
}

func writeError(w http.ResponseWriter, svcError model.ProxyError) {
	status, err := strconv.Atoi(svcError.Status)
	if err != nil {
		log.Errorf("Error writing error response %v", err)
		w.Write([]byte(svcError.Message))
		return
	}
	w.WriteHeader(status)

	jsonStr, err := json.Marshal(svcError)
	if err != nil {
		log.Errorf("Error writing error response %v", err)
		w.Write([]byte(svcError.Message))
		return
	}
	w.Write([]byte(jsonStr))
}
