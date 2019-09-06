//echo_websocket.go

package main

import (
	"log"
	"github.com/sacOO7/gowebsocket"
	"os"
	"os/signal"
	"encoding/base64"
	"fmt"
	"encoding/json"
	"time"
	"strings"
	"io"
)

type OrderMessage struct {
	Message []interface{} `json:"emit"`
}

type AppConfig struct {
	AppVersion string
	ClientID string
	ServerUrl string
}

var app_config AppConfig

//var client_id = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJkb3NoaWkgc2VydmVyIiwic3ViIjp7ImZvciI6IkFwcENsaWVudElkIiwiaWQiOjE1OH0sImV4cCI6MTUzMzg3NzMxOX0.TFmgVw8uNr-U21b9y66SMPkx6RzxW-kP32tQd0jzcvA";

//var socket = gowebsocket.New("wss://sandbox.doshii.co/app/socket?auth=" + base64.StdEncoding.EncodeToString([]byte(client_id)))

var client_id string

var socket gowebsocket.Socket

func main() {

	InitApp()

	SetupLog()

	client_id = app_config.ClientID

	socket = gowebsocket.New(app_config.ServerUrl + base64.StdEncoding.EncodeToString([]byte(client_id))) 

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	socket.OnConnectError = func(err error, socket gowebsocket.Socket) {
		log.Fatalf("Received connect error - ", err)
	}
  
	socket.OnConnected = func(socket gowebsocket.Socket) {
		log.Println("Connected to server");

		ticker := time.NewTicker(5 * time.Second)

		go func() {

			for t := range ticker.C {
				log.Println("Ping sent at", t)
				Heartbeat()
			}

		}()
	
	}
  
	socket.OnTextMessage = func(message string, socket gowebsocket.Socket) {
		
		is_pong_message := strings.Contains(message, "pong")
		is_checkin_message := strings.Contains(message, "checkin_created")
		is_order_message := strings.Contains(message, "order_created")

		if is_pong_message {
			log.Println("Received PONG message - " + message)
		}

		if is_checkin_message {
			log.Println("Received CHECK-IN message - " + message)
		}

		if is_order_message {

			log.Println("Received ORDER message - " + message)

			order_message := OrderMessage{}

			err := json.Unmarshal([]byte(message), &order_message)
			
			if err != nil {
				panic(err)
			}

			fmt.Printf("%#v\n", order_message)

			order_data_message, _ := order_message.Message[1].(map[string]interface{})
			log.Println("Processed ORDER ID - " + order_data_message["id"].(string))
		}
	}
  
	socket.OnPingReceived = func(data string, socket gowebsocket.Socket) {
		log.Println("Received ping - " + data)
	}
  
    socket.OnPongReceived = func(data string, socket gowebsocket.Socket) {
		log.Println("Received pong - " + data)
	}

	socket.OnDisconnected = func(err error, socket gowebsocket.Socket) {
		log.Println("Disconnected from server ")
		return
	}
  
	socket.Connect()

  	for {
		select {
			case <-interrupt:
				log.Println("interrupt")
				socket.Close()
				return
		}
	}
}

func InitApp() {

	config_file, err := os.Open("./config/app.json")

	defer config_file.Close()
	
	if err != nil {
		log.Fatalf("[openAppConfigErr]: %s\n", err)
	}
	
	decoder := json.NewDecoder(config_file)
	
	app_config = AppConfig{}
	
	if err = decoder.Decode(&app_config); err != nil {
		log.Fatalf("[decodeAppConfigErr]: %s\n", err)
	}
}

func SetupLog() {
	log_file, err := os.OpenFile("socket.log", os.O_RDWR|os.O_CREATE|os.O_APPEND, 0666)

	//defer log_file.Close()

	if err != nil {
	    log.Fatalf("error opening file: %v", err)
	}
	
	log_writer := io.MultiWriter(os.Stdout, log_file)
	
	log.SetOutput(log_writer)
}

func Heartbeat() {
	ping_data := make(map[string]interface{})

  	now := time.Now()
    unixtime := now.Unix()

  	ping_data["doshii"] = map[string]interface{}{"ping":unixtime, "version":"12345"}

  	ping_data_json, _ := json.Marshal(ping_data)
    
    log.Println("Sent PING message - " + string(ping_data_json))

    socket.SendText(string(ping_data_json))
}