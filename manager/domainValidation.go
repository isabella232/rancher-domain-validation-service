package manager

import (
	"database/sql"
	"fmt"

	log "github.com/Sirupsen/logrus"
)

//URL for Rancher server
var URL = "http://54.255.182.226:8080/"

//Port for rancher auth
var Port = "8092"

//DatabaseName for domian
var DatabaseName = "domian4"

//DomainTable for domian
var DomainTable = "domain_list"

//CreateDatabase is for CreateDatabase
func CreateDatabase() {

	db, err := sql.Open("mysql", "newuser:mynewpassword@tcp(54.255.182.226:3306)/?charset=utf8") //第一个参数为驱动名
	checkErr(err)

	//create domain Database
	_, err = db.Query(fmt.Sprintf("create database %s", DatabaseName))
	if err != nil {
		log.Errorf("Cannot create database: %v", err)
	}
	_, err = db.Query(fmt.Sprintf("CREATE TABLE `%s`.`%s`  (  `id` INT NOT NULL AUTO_INCREMENT,  `domainid` VARCHAR(45) NULL , `accountid` VARCHAR(45) NULL,  `projectid` VARCHAR(45) NULL,  `state` VARCHAR(45) NULL,  `hashvalue` VARCHAR(45) NULL,  `domain_name` VARCHAR(45) NULL,  PRIMARY KEY (`id`),  UNIQUE INDEX `domainid_UNIQUE` (`domainid` ASC));", DatabaseName, DomainTable))
	if err != nil {
		log.Errorf("Cannot create table: %v", err)
	}
	fmt.Println("Database Created!")
	// db.Query("drop database if exists tmpdb")
	// query, err := db.Query("SELECT * FROM domian2.domain_list;")
	// checkErr(err)
	// v := reflect.ValueOf(query)
	// fmt.Println(v)
	db.Close()
}

func checkErr(errMasg error) {
	if errMasg != nil {
		panic(errMasg)
	}
}
