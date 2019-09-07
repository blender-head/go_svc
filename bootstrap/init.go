package bootstrap

import (
	"log"
	"os"
	"fmt"
	"io"
	"encoding/json"
	"database/sql"
	_ "github.com/go-sql-driver/mysql"
	"path/filepath"
)

type AppConfigData struct {
	App_Version string
	Client_Id string
	Server_Url string
}

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

var AppConfig AppConfigData

func SetupLog() {
	if _, err := os.Stat("./logs"); os.IsNotExist(err) {
    	err := os.Mkdir("." + string(filepath.Separator) + "logs",0777)

	    if err != nil {
	    	log.Fatalf("error creating log dir: %v", err)
	    }
	}

	log_file, err := os.OpenFile("./logs/socket.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	//defer log_file.Close()

	if err != nil {
	    log.Fatalf("error opening file: %v", err)
	}
	
	log_writer := io.MultiWriter(os.Stdout, log_file)
	
	log.SetOutput(log_writer)
}

func InitApp() {

	config_file, err := os.Open("./config/app.json")

	//defer config_file.Close()
	
	if err != nil {
		log.Fatalf("[openAppConfigErr]: %s\n", err)
	}
	
	decoder := json.NewDecoder(config_file)
	
	AppConfig = AppConfigData{}
	
	if err = decoder.Decode(&AppConfig); err != nil {
		log.Fatalf("[decodeAppConfigErr]: %s\n", err)
	}
}

func InitDB() {
	//config_path, _ := filepath.Abs("../go_emenu/config/db.json")

	file, err := os.Open("./config/db.json")

	//defer file.Close()
	
	if err != nil {
		log.Fatalf("[openDBConfigErr]: %s\n", err)
	}
	
	decoder := json.NewDecoder(file)
	
	DBConfig = DBConfigData{}
	
	if err = decoder.Decode(&DBConfig); err != nil {
		log.Fatalf("[decodeDBConfigErr]: %s\n", err)
	}

	conn_detail := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", DBConfig.User, DBConfig.Password, DBConfig.Host, DBConfig.Port, DBConfig.Database)

	//db, err_db_connect = sql.Open("mysql", "root:mtu1500@andre@tcp(172.17.0.4:3306)/go_emenu?charset=utf8")

	if DB, err = sql.Open("mysql", conn_detail); err != nil {
		log.Fatalf("[dbConnErr]: %s\n", err)
	}
}