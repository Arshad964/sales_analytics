package dbconnect

import (
	"fmt"
	"salesAnalysisLumel/common"
)

func (d *AllUsedDatabases) Init() {
	dbconfig := common.ReadTomlConfig("./toml/dbconfig.toml")

	d.Mariadb.Server = fmt.Sprintf("%v", dbconfig.(map[string]interface{})["MariaServer"])
	d.Mariadb.Port = fmt.Sprintf("%v", dbconfig.(map[string]interface{})["MariaPort"])
	d.Mariadb.User = fmt.Sprintf("%v", dbconfig.(map[string]interface{})["MariaUser"])
	d.Mariadb.Password = fmt.Sprintf("%v", dbconfig.(map[string]interface{})["MariaPassword"])
	d.Mariadb.Database = fmt.Sprintf("%v", dbconfig.(map[string]interface{})["MariaDatabase"])
	d.Mariadb.DB = Maria
}
