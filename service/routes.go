package service

import "github.com/gorilla/mux"

//NewRouter creates and configures a mux router
func NewRouter() *mux.Router {
	router := mux.NewRouter()
	//list domains
	router.HandleFunc("/v1-domains/domains", GetDomains).Methods("GET")
	//get domain by id
	router.HandleFunc("/v1-domains/domains/{id}", GetDomains).Methods("GET")
	//get domianby id?="XX"
	router.HandleFunc("/v1-domains/domain", GetDomains).Methods("GET")
	//add new domain
	router.HandleFunc("/v1-domains/domain", CreateDomain).Methods("POST")
	//validate domain
	router.HandleFunc("/v1-domains/domain/{id}", ActivateDomain).Methods("POST")
	//delete domain
	router.HandleFunc("/v1-domains/domain/{id}", DeleteDomain).Methods("DELETE")
	//domain filter id=envid
	router.HandleFunc("/v1-domains/filter/projects/{id}/loadbalancerservice", ValidateDomian).Methods("POST")

	return router
}
