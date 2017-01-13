package service

import "github.com/gorilla/mux"

//NewRouter creates and configures a mux router
func NewRouter() *mux.Router {
	router := mux.NewRouter()
	//list domains
	router.HandleFunc("/v1-domains/domains/", GetDomains).Methods("GET")
	//get domain by id
	router.HandleFunc("/v1-domains/domains/{id}", GetDomains).Methods("GET")
	//get domianby id?="XX"
	router.HandleFunc("/v1-domains/domain", GetDomains).Methods("GET")
	//add new domain
	router.HandleFunc("/v1-domains/domain", AddDomain).Methods("POST")
	//delete domain
	router.HandleFunc("/v1-domains/domain", DeleteDomain).Methods("DELETE")
	//domain filter action?
	router.HandleFunc("/v1-domains/validate", DeleteDomain).Methods("POST")
	return router
}
