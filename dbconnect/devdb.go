package dbconnect

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
)

type DatabaseType struct {
	Server   string
	Port     string
	User     string
	Password string
	Database string
	// DBType   string
	DB string
}

type AllUsedDatabases struct {
	Mariadb DatabaseType
}

var Maria = "Maria"

func LocalDbConnect(DBtype string) (*sql.DB, error) {
	DbDetails := new(AllUsedDatabases)
	DbDetails.Init()

	//log.Println(DbDetails)

	connString := ""
	// localDBtype := ""

	var db *sql.DB
	var err error
	var dataBaseConnection DatabaseType
	log.Println("DBtype", DBtype)
	// get connection details
	if DBtype == DbDetails.Mariadb.DB {
		dataBaseConnection = DbDetails.Mariadb
		// localDBtype = DbDetails.Mariadb.DB
		connString = `` + dataBaseConnection.User + `:` + dataBaseConnection.Password + `@tcp(` + dataBaseConnection.Server + `:` + dataBaseConnection.Port + `)/` + dataBaseConnection.Database + `?parseTime=true`
		// connString = `user=` + dataBaseConnection.User + ` password=` + dataBaseConnection.Password + ` port=` + dataBaseConnection.Port + ` dbname=` + dataBaseConnection.Database + ` host=` + dataBaseConnection.Server + ` sslmode=disable`
		db, err = sql.Open("mysql", connString)
		if err != nil {
			log.Println("Open connection failed:", err.Error())
		}
	} else {
		return db, fmt.Errorf(" Invalid DB Details")
	}

	/* // Prepare connection string
	if localDBtype == Maria {
		log.Println("IN", localDBtype)
		// connString = `user=` + dataBaseConnection.User + ` password=` + dataBaseConnection.Password + ` port=` + fmt.Sprintf("%v", dataBaseConnection.Port) + ` dbname=` + dataBaseConnection.Database + ` host=` + dataBaseConnection.Server + ` sslmode=disable`
	}

	log.Println(localDBtype, "localDBtype")
	log.Println("connString - ", connString)
	//make a connection to db
	if localDBtype != "" {
		db, err = sql.Open(localDBtype, connString)
		if err != nil {
			log.Println("Open connection failed:", err.Error())
		}
	} else {
		return db, fmt.Errorf(" Invalid DB Details")
	} */

	return db, err
}

func ExecuteBulkStatement(db *sql.DB, sqlStringValues string, sqlString string) error {
	log.Println("ExecuteBulkStatement+")
	//trim the last ,
	log.Println("query :", sqlString+sqlStringValues)
	sqlStringValues = sqlStringValues[0 : len(sqlStringValues)-1]
	_, err := db.Exec(sqlString + sqlStringValues)
	if err != nil {
		log.Println(err)
		log.Println("ExecuteBulkStatement-")
		return err
	} else {
		log.Println("inserted Sucessfully")
	}
	log.Println("ExecuteBulkStatement-")
	return nil
}
