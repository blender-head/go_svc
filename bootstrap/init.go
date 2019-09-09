package bootstrap

import (
	"log"
	"os"
	"io"
	"encoding/json"
	"path/filepath"
	"time"
)

type AppConfigData struct {
	App_Version string
	Client_Id string
	Server_Url string
}

var AppConfig AppConfigData

func SetupLog() {

	parent_log := "logs"

	t := time.Now()
	dated_logs := t.Format("2006-01-02")

	log_filename := "socket.log"

	if _, err := os.Stat("./" + parent_log); os.IsNotExist(err) {
    	err := os.Mkdir("." + string(filepath.Separator) + parent_log,0777)

	    if err != nil {
	    	log.Fatalf("error creating parent log dir: %s\n", err)
	    }
	}

	if _, err := os.Stat("./" + parent_log + "/" + dated_logs); os.IsNotExist(err) {
    	err := os.Mkdir("." + string(filepath.Separator) + parent_log + "/" + dated_logs,0777)

	    if err != nil {
	    	log.Fatalf("error creating dated log dir: %s\n", err)
	    }
	}

	log_file, err := os.OpenFile("./" + parent_log + "/" + dated_logs + "/" + log_filename, os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	if err != nil {
	    log.Fatalf("error creating log file: %s\n", err)
	}
	
	log_writer := io.MultiWriter(os.Stdout, log_file)
	
	log.SetOutput(log_writer)
}

func InitApp() {

	config_file, err := os.Open("./config/app.json")

	if err != nil {
		log.Fatalf("error opening app config file: %s\n", err)
	}
	
	decoder := json.NewDecoder(config_file)
	
	AppConfig = AppConfigData{}
	
	if err = decoder.Decode(&AppConfig); err != nil {
		log.Fatalf("error decoding app config file: %s\n", err)
	}
}