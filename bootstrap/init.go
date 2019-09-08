package bootstrap

import (
	"log"
	"os"
	"io"
	"encoding/json"
	"path/filepath"
)

type AppConfigData struct {
	App_Version string
	Client_Id string
	Server_Url string
}

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