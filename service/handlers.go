package service

import (
	"crypto/rand"
	"database/sql"
	"encoding/json"
	"errors"
	"io/ioutil"
	"net"
	"net/http"
	"strconv"
	"strings"

	"fmt"

	log "github.com/Sirupsen/logrus"
	"github.com/deckarep/golang-set"
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

//DomainRequest is for error/success the JSON output
type DomainRequest struct {
	DomanName string `json:"domainName,omitempty"`
	ProjectID string `json:"projectid,omitempty"`
}

//LbPost is for error/success the JSON output
type LbPost struct {
	LbConfig LbConfigData `json:"lbConfig,omitempty"`
}

//LbConfigData is the data inside LbConfig
type LbConfigData struct {
	// CertificateIds string         `json:"certificateIds,omitempty"`
	// Config         string         `json:"config,omitempty"`
	PortRules []PortRuleData `json:"portRules,omitempty"`
	// StickinessPolicy string         `json:"stickinessPolicy,omitempty"`
}

//PortRuleData is the data inside PortRules
type PortRuleData struct {
	Hostname string `json:"hostname,omitempty"`
	// Priority   string `json:"priority,omitempty"`
	// Protocol   string `json:"protocol,omitempty"`
	// SourcePort string `json:"sourcePort,omitempty"`
	// TargetPort string `json:"targetPort,omitempty"`
}

//DomainList is for listing out the domains
type DomainList struct {
	DomainID   string `json:"domainId,omitempty"`
	AccountID  string `json:"accountId,omitempty"`
	ProjectID  string `json:"projectId,omitempty"`
	DomianName string `json:"domainName,omitempty"`
	State      string `json:"state,omitempty"`
	Hashvalue  string `json:"hashvalue,omitempty"`
}

//GetDomains return the list of domain
func GetDomains(w http.ResponseWriter, r *http.Request) {
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
	query, err := db.Query(fmt.Sprintf("SELECT * FROM %s.%s WHERE accountid='%s';", manager.DatabaseName, manager.DomainTable, accountID))
	fmt.Printf(fmt.Sprintf("SELECT * FROM %s.%s WHERE accountid='%s';", manager.DatabaseName, manager.DomainTable, accountID))
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
	// w.WriteHeader(http.StatusOK)

	db.Close()
}

//CreateDomain into the domain list
func CreateDomain(w http.ResponseWriter, r *http.Request) {

	var jsonInput = DomainRequest{}
	input, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("Error reading request Body :%v", err)
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
	db, err := sql.Open("mysql", "newuser:mynewpassword@tcp(54.255.182.226:3306)/?charset=utf8")
	if err != nil {
		log.Errorf("Error connecting to database: %v", err)
		ReturnHTTPError(w, r, "BadRequest", http.StatusBadRequest, fmt.Sprintf("Error connecting to database: %v", err))
		return
	}
	accountfromtoken, err := getAccountID(r)
	if err != nil {
		log.Errorf("Error writing error response %v", err)
		ReturnHTTPError(w, r, "InternalServerError", http.StatusInternalServerError, fmt.Sprintf("Error marshal the result to jsonstring: %v", err))
		return
	}
	vars := mux.Vars(r)
	containerID := vars["id"]
	r.ParseForm()
	if len(r.Form["action"]) > 0 {
		if r.Form["action"][0] != "validate" {
			ReturnHTTPError(w, r, "BadRequest", http.StatusBadRequest, fmt.Sprintf("Action not valid"))
		}
	}

	query, err := db.Query(fmt.Sprintf("SELECT * FROM %s.%s WHERE id='%s';", manager.DatabaseName, manager.DomainTable, decodeDomainID(containerID)))

	if err != nil {
		log.Errorf("Error querying the statement: %v", err)
		ReturnHTTPError(w, r, "NotFound", 404, fmt.Sprintf("Error querying the statement %v", err))
		return
	}
	queryResult, err := praseQueryResult(query)

	if len(queryResult) != 1 {
		log.Errorf("Domain not find")
		ReturnHTTPError(w, r, "NotFound", 404, fmt.Sprintf("Domain not find"))
		return
	}

	//get information from sql query
	accountidfromdb := queryResult[0].AccountID
	hashvaluefromdb := queryResult[0].Hashvalue

	if accountfromtoken != accountidfromdb && accountidfromdb != "" {
		log.Errorf("token unauthorized ")
		ReturnHTTPError(w, r, "Unauthorized", http.StatusInternalServerError, fmt.Sprintf("token unauthorized"))
		return
	}
	//Get the txt DNS Record
	txt, err := net.LookupTXT("_hna-challenge.fiduccia.me")

	if len(txt[0]) >= 1 {
		// fmt.Println("TXT record" + txt[0])
		if txt[0] != hashvaluefromdb {
			if accountfromtoken != hashvaluefromdb {
				log.Errorf("DNS txt record not valid")
				// ReturnHTTPError(w, r, "NotFound", http.StatusNotFound, fmt.Sprintf("DNS txt record not valid"))
				//return
			}
		}
	}

	//validate the acme challenge
	requestURL := "http://" + queryResult[0].DomianName + "/.well-know/hna.txt"
	req, err := http.NewRequest("GET", requestURL, nil)
	if err != nil {
		log.Errorf("Error get the acme challenge: %v", err)
		ReturnHTTPError(w, r, "NotFound", 404, fmt.Sprintf("Error get the acme challenge %v", err))
		return
	}
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("Error get the acme challenge: %v", err)
		ReturnHTTPError(w, r, "NotFound", 404, fmt.Sprintf("Error get the acme challenge %v", err))
		return
	}
	bodyText, err := ioutil.ReadAll(resp.Body)

	trimedbodytext := strings.Trim(string(bodyText), "\n")
	if trimedbodytext != hashvaluefromdb {
		log.Errorf("Error get the acme challenge not match")
		ReturnHTTPError(w, r, "Unauthorized", 401, fmt.Sprintf("Error get the acme challenge"))
		return
	}

	_, err = db.Query(fmt.Sprintf("UPDATE `%s`.`%s` SET `%s`.`%s`.`state` = 'active' WHERE (`id`='%s') LIMIT 1;", manager.DatabaseName, manager.DomainTable, manager.DatabaseName, manager.DomainTable, decodeDomainID(containerID)))

	if err != nil {
		log.Errorf("Error updating the statement: %v", err)
		ReturnHTTPError(w, r, "NotFound", 404, fmt.Sprintf("Error updating the statement %v", err))
		return
	}

	ReturnHTTPSuccess(w, r, "Succeed", http.StatusOK, fmt.Sprintf("update the record succeed"))

}

//DeleteDomain is for DeleteDomain
func DeleteDomain(w http.ResponseWriter, r *http.Request) {

	accountID, err := getAccountID(r)
	if err != nil {
		log.Errorf("Error get the account ID, please check the rancher server URL: %v", err)
		ReturnHTTPError(w, r, "InternalServerError", http.StatusInternalServerError, fmt.Sprintf("Error get the account ID, please check the rancher server URL: %v", err))
		return
	}
	vars := mux.Vars(r)
	containerID := vars["id"]

	db, err := sql.Open("mysql", "newuser:mynewpassword@tcp(54.255.182.226:3306)/?charset=utf8")
	if err != nil {
		log.Errorf("Error connecting to database: %v", err)
		ReturnHTTPError(w, r, "BadRequest", http.StatusBadRequest, fmt.Sprintf("Error connecting to database: %v", err))
		return
	}
	_, err = db.Query(fmt.Sprintf("DELETE FROM `%s`.`%s` WHERE `id`='%s' and `accountid`='%s' ;", manager.DatabaseName, manager.DomainTable, decodeDomainID(containerID), accountID))
	fmt.Printf("DELETE FROM `%s`.`%s` WHERE `id`='%s' and `accountid`='%s' ;", manager.DatabaseName, manager.DomainTable, decodeDomainID(containerID), accountID)
	if err != nil {
		log.Errorf("Error delete the record: %v", err)
		ReturnHTTPError(w, r, "BadRequest", http.StatusBadRequest, fmt.Sprintf("Error inserting the record: %v", err))
		return
	}
	ReturnHTTPSuccess(w, r, "Succeed", http.StatusOK, fmt.Sprintf("Delete the record succeed"))
	db.Close()
}

//ValidateDomian filter
func ValidateDomian(w http.ResponseWriter, r *http.Request) {

	vars := mux.Vars(r)
	ProjectID := vars["id"]

	csrf := r.Header.Get("x-api-csrf")
	if csrf == "" {
		log.Errorf("Error reading request Body , err: no x-api-csrf ")
		ReturnHTTPError(w, r, "BadRequest", http.StatusBadRequest, fmt.Sprintf("Error reading request Body , err: no x-api-csrf "))
		return
	}
	var jsonInput = LbPost{}
	input, err := ioutil.ReadAll(r.Body)
	if err != nil {
		log.Errorf("Error reading request Body :%v", err)
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
	query, err := db.Query(fmt.Sprintf("SELECT domain_name FROM %s.%s WHERE accountid='%s' AND state='active' AND projectid='%s';", manager.DatabaseName, manager.DomainTable, accountID, ProjectID))
	fmt.Printf("SELECT domain_name FROM %s.%s WHERE accountid='%s' AND state='active' AND projectid='%s';", manager.DatabaseName, manager.DomainTable, accountID, ProjectID)

	if err != nil {
		log.Errorf("Error inserting the record: %v", err)
		ReturnHTTPError(w, r, "BadRequest", http.StatusBadRequest, fmt.Sprintf("Error inserting the record: %v", err))
		return
	}

	// fmt.Println(jsonInput.LbConfig.PortRules[0].Hostname)
	if len(jsonInput.LbConfig.PortRules) >= 1 {

	}
	queryResult, err := praseQueryResult(query)
	validDomainlist := mapset.NewSet()
	fmt.Println("extract domain name from db")
	for _, v := range queryResult {
		if v.DomainID != "" {
			fmt.Println(v.DomainID)
			validDomainlist.Add(v.DomainID)
		}
	}

	//validation the route rule in domain list
	fmt.Println("compare domain name from db")
	for _, v := range jsonInput.LbConfig.PortRules {
		fmt.Printf(v.Hostname)
		if validDomainlist.Contains(v.Hostname) {
			fmt.Println(v.Hostname)
		} else {
			log.Errorf("Domain %s is not valid", v.Hostname)
			ReturnHTTPError(w, r, "Forbidden", 403, fmt.Sprintf("Domain %s is not valid", v.Hostname))
			return
		}
	}
	//if all the domain name in the route rule is valid, then redirect the post to GLB
	cookie, err := r.Cookie("token")
	if err != nil {
		log.Errorf("Token not found in the request")
		ReturnHTTPError(w, r, "BadRequest", 400, fmt.Sprintf("Token not found in the request"))
		return
	}

	client := &http.Client{}
	requestURL := fmt.Sprintf("%s/v2-beta/projects/%s/loadbalancerservice", manager.URL, ProjectID)
	req, err := http.NewRequest("GET", requestURL, nil)
	// postCookie := http.Cookie{Name: "token", cookie.Value}
	postcookie := http.Cookie{Name: "token", Value: cookie.Value}
	req.AddCookie(&postcookie)
	req.Header.Add("x-api-csrf", csrf)
	resp, err := client.Do(req)
	if err != nil {
		log.Errorf("Error connecting to rancher server and create LB: %v", err)
		ReturnHTTPError(w, r, "NotFound", http.StatusNotFound, fmt.Sprintf("Error connecting to rancher server and create LB: %v", err))
		return
	}
	bodyText, err := ioutil.ReadAll(resp.Body)
	w.Write(bodyText)

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
		if accountIDData[0] != "ID_NOT_FIND" && accountIDData[0] != "Unauthorized" {
			return accountIDData[0], nil
		}
		err = errors.New("Cannot extract account id")
		return "", err

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

func praseQueryResult(query *sql.Rows) ([]DomainList, error) {
	column, _ := query.Columns()              //读出查询出的列字段名
	values := make([][]byte, len(column))     //values是每个列的值，这里获取到byte里
	scans := make([]interface{}, len(column)) //因为每次查询出来的列是不定长的，用len(column)定住当次查询的长度
	for i := range values {                   //让每一行数据都填充到[][]byte里面
		scans[i] = &values[i]
	}
	results := []DomainList{} // 最后得到的map
	for query.Next() {        //循环，让游标往下移动
		if err := query.Scan(scans...); err != nil { //query.Scan查询出来的不定长值放到scans[i] = &values[i],也就是每行都放在values里
			return nil, err
		}
		row := DomainList{}        //每行数据
		for k, v := range values { //每行数据是放在values里面，现在把它挪到row里
			switch column[k] {
			case "domainid":
				row.DomainID = string(v)
			case "accountid":
				row.AccountID = string(v)
			case "projectid":
				row.ProjectID = string(v)
			case "domain_name":
				row.DomianName = string(v)
			case "hashvalue":
				row.Hashvalue = string(v)
			case "state":
				row.State = string(v)
			default:

			}
		}
		results = append(results, row) //装入结果集中

	}

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

//convert the cattle format to the domain id in database
func decodeDomainID(cattleid string) string {
	return cattleid[2:]
}

//convert the domain id in database to the cattle format
func encodeDomainID(id string) string {
	return "1d" + id
}
