package service

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net/http"
	"strconv"

	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/gorilla/mux"
	"github.com/rancher/rancher-domain-validaiton-service/manager"
	"github.com/rancher/rancher-domain-validaiton-service/model"
)

var db *sql.DB

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

//DomainRequest is for the JSON output
type DomainRequest struct {
	DomanName string `json:"domainName,omitempty"`
	ProjectID string `json:"projectid,omitempty"`
}

type DomainList struct {
	ID        string `json:"id,omitempty"`
	Accountid string `json:"accountid,omitempty"`
	ProjectID string `json:"projectidid,omitempty"`
	State     string `json:"state,omitempty"`
	Hashvalue string `json:"Hashvalue,omitempty"`
}

//GetDomains return the list of domain
func GetDomains(w http.ResponseWriter, r *http.Request) {
	// // reqestData := RequestData{}
	// // input, err := ioutil.ReadAll(r.Body)

	// db, err := sql.Open("mysql", "newuser:mynewpassword@tcp(54.255.182.226:3306)/?charset=utf8")
	// checkErr(err)

	// // query, err := db.Query("SELECT * FROM domian2.domain_list;")
	// checkErr(err)
	// // v := reflect.ValueOf(query)
	// // printResult(query)
	// // fmt.Println(v)
	// // printResult(query)
	// db.Close()

	// path, _ := mux.CurrentRoute(r).GetPathTemplate()
	// var jsonInput = DomainRequest{}
	// input, err := ioutil.ReadAll(r.Body)
	// if err != nil {
	// 	log.Errorf("Error reading request Body %v for path %v", r, path)
	// 	ReturnHTTPError(w, r, "BadRequest", http.StatusBadRequest, fmt.Sprintf("Error reading request Body , err: %v", err))
	// 	return
	// }
	// err = json.Unmarshal(input, &jsonInput)
	// if err != nil {
	// 	log.Errorf("Error unmarshalling json request body: %v", err)
	// 	ReturnHTTPError(w, r, "BadRequest", http.StatusBadRequest, fmt.Sprintf("Error unmarshalling json request body: %v", err))
	// 	return
	// }
	accountID, err := getAccountID(r)
	if err != nil {
		log.Errorf("Error get the account ID, please check the rancher server URL: %v", err)
		ReturnHTTPError(w, r, "InternalServerError", http.StatusInternalServerError, fmt.Sprintf("Error get the account ID, please check the rancher server URL: %v", err))
		return
	}

	db, err := sql.Open("mysql", "newuser:mynewpassword@tcp(54.255.182.226:3306)/?charset=utf8")
	if err != nil {
		log.Errorf("Error connecting to database: %v", err)
		ReturnHTTPError(w, r, "BadRequest", http.StatusBadRequest, fmt.Sprintf("Error connecting to database: %v", err))
		return
	}
	query, err := db.Query(fmt.Sprintf("SELECT * FROM domian3.domain_list WHERE accountid='%s';", accountID))
	// fmt.Printf("INSERT INTO `%s`.`%s` (`accountid`, `projectid`, `state`, `hashvalue`, `domain_name`) VALUES ('%s', '%s', '%s', '%s', '%s');", manager.DatabaseName, manager.DomainTable, accountID, jsonInput.ProjectID, "Pending", randToken(), jsonInput.DomanName)
	if err != nil {
		log.Errorf("Error inserting the record: %v", err)
		ReturnHTTPError(w, r, "BadRequest", http.StatusBadRequest, fmt.Sprintf("Error inserting the record: %v", err))
		return
	}

	queryResult, err := praseQueryResult(query)

	jsonStr, err := json.Marshal(queryResult)
	if err != nil {
		log.Errorf("Error writing error response %v", err)
		ReturnHTTPError(w, r, "InternalServerError", http.StatusInternalServerError, fmt.Sprintf("Error marshal the result to jsonstring: %v", err))
		return
	}
	w.Write([]byte(jsonStr))
	w.WriteHeader(http.StatusOK)

	db.Close()

}

//CreateDomain into the domain list
func CreateDomain(w http.ResponseWriter, r *http.Request) {
	path, _ := mux.CurrentRoute(r).GetPathTemplate()
	var jsonInput = DomainRequest{}
	input, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("Error reading request Body %v for path %v", r, path)
		ReturnHTTPError(w, r, "BadRequest", http.StatusBadRequest, fmt.Sprintf("Error reading request Body , err: %v", err))
		return
	}
	err = json.Unmarshal(input, &jsonInput)
	if err != nil {
		log.Errorf("Error unmarshalling json request body: %v", err)
		ReturnHTTPError(w, r, "BadRequest", http.StatusBadRequest, fmt.Sprintf("Error unmarshalling json request body: %v", err))
		return
	}

	accountID, err := getAccountID(r)
	if err != nil {
		log.Errorf("Error get the account ID, please check the rancher server URL: %v", err)
		ReturnHTTPError(w, r, "InternalServerError", http.StatusInternalServerError, fmt.Sprintf("Error get the account ID, please check the rancher server URL: %v", err))
		return
	}

	db, err := sql.Open("mysql", "newuser:mynewpassword@tcp(54.255.182.226:3306)/?charset=utf8")
	if err != nil {
		log.Errorf("Error connecting to database: %v", err)
		ReturnHTTPError(w, r, "BadRequest", http.StatusBadRequest, fmt.Sprintf("Error connecting to database: %v", err))
		return
	}
	_, err = db.Query(fmt.Sprintf("INSERT INTO `%s`.`%s` (`accountid`, `projectid`, `state`, `hashvalue`, `domain_name`) VALUES ('%s', '%s', '%s', '%s', '%s');", manager.DatabaseName, manager.DomainTable, accountID, jsonInput.ProjectID, "Pending", randToken(), jsonInput.DomanName))
	fmt.Printf("INSERT INTO `%s`.`%s` (`accountid`, `projectid`, `state`, `hashvalue`, `domain_name`) VALUES ('%s', '%s', '%s', '%s', '%s');", manager.DatabaseName, manager.DomainTable, accountID, jsonInput.ProjectID, "Pending", randToken(), jsonInput.DomanName)
	if err != nil {
		log.Errorf("Error inserting the record: %v", err)
		ReturnHTTPError(w, r, "BadRequest", http.StatusBadRequest, fmt.Sprintf("Error inserting the record: %v", err))
		return
	}

	ReturnHTTPSuccess(w, r, "Succeed", http.StatusOK, fmt.Sprintf("inserting the record succeed"))

	// query, err := db.Query(fmt.Sprintf(

	fmt.Println()
	// fmt.Println(praseQueryResult(query))

}

//ActivateDomain is for activating the domain
func ActivateDomain(w http.ResponseWriter, r *http.Request) {
}

//DeleteDomain is for DeleteDomain
func DeleteDomain(w http.ResponseWriter, r *http.Request) {

	
	accountID, err := getAccountID(r)
	if err != nil {
		log.Errorf("Error get the account ID, please check the rancher server URL: %v", err)
		ReturnHTTPError(w, r, "InternalServerError", http.StatusInternalServerError, fmt.Sprintf("Error get the account ID, please check the rancher server URL: %v", err))
		return
	}

	db, err := sql.Open("mysql", "newuser:mynewpassword@tcp(54.255.182.226:3306)/?charset=utf8")
	if err != nil {
		log.Errorf("Error connecting to database: %v", err)
		ReturnHTTPError(w, r, "BadRequest", http.StatusBadRequest, fmt.Sprintf("Error connecting to database: %v", err))
		return
	}
	query, err := db.Query(fmt.Sprintf("DELETE FROM `domian3`.`domain_list` WHERE `id`='21';", accountID))
	// fmt.Printf("INSERT INTO `%s`.`%s` (`accountid`, `projectid`, `state`, `hashvalue`, `domain_name`) VALUES ('%s', '%s', '%s', '%s', '%s');", manager.DatabaseName, manager.DomainTable, accountID, jsonInput.ProjectID, "Pending", randToken(), jsonInput.DomanName)
	if err != nil {
		log.Errorf("Error inserting the record: %v", err)
		ReturnHTTPError(w, r, "BadRequest", http.StatusBadRequest, fmt.Sprintf("Error inserting the record: %v", err))
		return
	}

	queryResult, err := praseQueryResult(query)

	jsonStr, err := json.Marshal(queryResult)
	if err != nil {
		log.Errorf("Error writing error response %v", err)
		ReturnHTTPError(w, r, "InternalServerError", http.StatusInternalServerError, fmt.Sprintf("Error marshal the result to jsonstring: %v", err))
		return
	}
	w.Write([]byte(jsonStr))
	w.WriteHeader(http.StatusOK)

	db.Close()
}

//ValidateDomian filter
func ValidateDomian(w http.ResponseWriter, r *http.Request) {
}

func checkErr(errMasg error) {
	if errMasg != nil {
		panic(errMasg)
	}
}

func getAccountID(r *http.Request) (string, error) {

	//get accountid from token
	cookie, err := r.Cookie("token")
	if err != nil {
		return "", err
	}
	log.Infof("token:" + cookie.Value)
	accountIDData := getValue(manager.URL, "accounts", cookie.Value)
	if len(accountIDData) == 1 {
		return accountIDData[0], nil
	}
	err = errors.New("Cannot extract account id")
	return "", err

}

// func praseStruct(queryResult map[int]map[string]string) map[string]string {
// 	var result = make(map[string]string)
// 	for _, v := range queryResult {
// 		for key, value := range v {
// 			result[key] = value

// 		}
// 	}
// 	return result
// }

//ReturnHTTPError handles sending out CatalogError response
func ReturnHTTPError(w http.ResponseWriter, r *http.Request, status string, httpStatus int, errorMessage string) {
	svcError := model.DomainValidationErr{
		Type:    "error",
		Code:    status,
		Status:  strconv.Itoa(httpStatus),
		Message: errorMessage,
	}
	writeError(w, svcError)
}

func writeError(w http.ResponseWriter, svcError model.DomainValidationErr) {
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

//ReturnHTTPSuccess handles sending out CatalogError response
func ReturnHTTPSuccess(w http.ResponseWriter, r *http.Request, status string, httpStatus int, errorMessage string) {
	svcError := model.DomainValidationErr{
		Type:    "sucess",
		Code:    status,
		Status:  strconv.Itoa(httpStatus),
		Message: errorMessage,
	}
	writeError(w, svcError)
}

func praseQueryResult(query *sql.Rows) ([]map[string]string, error) {
	column, _ := query.Columns()              //读出查询出的列字段名
	values := make([][]byte, len(column))     //values是每个列的值，这里获取到byte里
	scans := make([]interface{}, len(column)) //因为每次查询出来的列是不定长的，用len(column)定住当次查询的长度
	for i := range values {                   //让每一行数据都填充到[][]byte里面
		scans[i] = &values[i]
	}

	results := make([]map[string]string, 0)
	count := 0
	for query.Next() { //循环，让游标往下移动
		row := make(map[string]string) //每行数据
		for k, v := range values {     //每行数据是放在values里面，现在把它挪到row里
			key := column[k]
			row[key] = string(v)
		}
		results = append(results, row) //装入结果集中
		count++
	}
	fmt.Println(count)
	return results, nil
}

//get the projectID and accountID from rancher API
func getValue(host string, path string, token string) []string {
	var result []string
	client := &http.Client{}
	requestURL := host + "v2-beta/" + path
	req, err := http.NewRequest("GET", requestURL, nil)
	cookie := http.Cookie{Name: "token", Value: token}
	req.AddCookie(&cookie)
	resp, err := client.Do(req)
	if err != nil {
		log.Infof("Cannot connect to the rancher server. Please check the rancher server URL")
		result = []string{"ID_NOT_FIND"}
		return result
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	authMessage := AuthorizeData{}
	err = json.Unmarshal(bodyText, &authMessage)
	if err != nil {
		log.Info(err)
		log.Infof("Cannot parse the authorization data.")
		result = []string{"ID_NOT_FIND"}
		return result
	}
	if authMessage.Message == "Unauthorized" {
		result = []string{"Unauthorized"}
	} else {
		messageData := MessageData{}
		err = json.Unmarshal(bodyText, &messageData)
		if err != nil {
			log.Info(err)
			log.Infof("Cannot parse the id.")
			result = []string{"ID_NOT_FIND"}

		}
		//get id from the data

		for i := 0; i < len(messageData.Data); i++ {

			idData, suc := messageData.Data[i].(map[string]interface{})
			if suc {
				id, suc := idData["id"].(string)
				name, namesuc := idData["uuid"].(string)
				if suc && namesuc {
					result = append(result, id)
					//if the token belongs to admin, only return the admin token
					if name == "admin" && path == "accounts" {
						result = []string{id}
						break
					}
				} else {
					log.Infof("No id find")
					result = []string{"ID_NOT_FIND"}
				}
			}

		}
		//get the admin user id. admin token will list all the ids. Need to just keep admin id.

	}

	return result
}
func randToken() string {
	b := make([]byte, 40)
	rand.Read(b)
	return fmt.Sprintf("%x", b)
}
