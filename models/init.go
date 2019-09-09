package models

import (
	"log"
	"os"
	"fmt"
	"encoding/json"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
)

type DBConfigData struct {
	Host		string
	Port		string
	User		string
	Password	string
	Database	string
	Charset 	string
}

var DB *sql.DB
var DBConfig DBConfigData

func InitDB() {
	file, err := os.Open("./config/db.json")

	if err != nil {
		log.Fatalf("error opening db config file: %s\n", err)
	}
	
	decoder := json.NewDecoder(file)
	
	DBConfig = DBConfigData{}
	
	if err = decoder.Decode(&DBConfig); err != nil {
		log.Fatalf("error decoding db config file: %s\n", err)
	}

	conn_detail := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", DBConfig.User, DBConfig.Password, DBConfig.Host, DBConfig.Port, DBConfig.Database)

	//db, err_db_connect = sql.Open("mysql", "root:mtu1500@andre@tcp(172.17.0.4:3306)/go_emenu?charset=utf8")

	if DB, err = sql.Open("mysql", conn_detail); err != nil {
		log.Fatalf("error connecting to database server: %s\n", err)
	}
}