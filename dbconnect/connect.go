package dbconnect

import (
	"database/sql"
	"log"
)

var GDBConnection *sql.DB

func DbConnect() error {
	var lErr error
	GDBConnection, lErr = LocalDbConnect(Maria)
	if lErr != nil {
		log.Println("DB01 Error in DB connect")
		return lErr
	}
	return nil
}
