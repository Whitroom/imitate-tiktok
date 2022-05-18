package sql

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var DB *gorm.DB

func InitDatabase() {

	type connection struct {
		Account  string `json:"account"`
		Password string `json:"password"`
		Database string `json:"database"`
	}

	conf, _ := os.Open("./confs/database.json")
	defer conf.Close()
	value, _ := ioutil.ReadAll(conf)
	var conn connection
	json.Unmarshal([]byte(value), &conn)

	dsn := fmt.Sprintf("%s:%s@tcp(localhost:3306)/%s?charset=utf8&parseTime=true",
		conn.Account, conn.Password, conn.Database) + "&loc=Asia%2fShanghai"

	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		fmt.Println("failed creating database:%w", err)
	}
	db.AutoMigrate(&User{}, &Video{}, &Comment{})
	DB = db
}
