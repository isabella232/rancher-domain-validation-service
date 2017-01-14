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

	// dbtemp, err := sql.Open("mysql", "root:password@/dbname?charset=utf8")
	// db := dbtemp
	// if err != nil {
	// 	fmt.Println(err)
	// }
	// stmt, err := db.Prepare("create table if not exists dev(id int UNSIGNED AUTO_INCREMENT NOT NULL PRIMARY KEY,uid varchar(64),did varchar(64),name varchar(64),qid varchar(64),status char DEFAULT 'u')")
	// if stmt != nil {
	// 	stmt.Exec()
	// 	stmt.Close()
	// }
	// stmt, err := db.Prepare("alter table dev convert to character set utf8 collate utf8_general_ci") //要修改一下编码
	// if stmt != nil {
	// 	stmt.Exec()
	// 	stmt.Close()
	// }
}
