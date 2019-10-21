package main

import (
	"log"
	"github.com/sacOO7/gowebsocket"
	"os"
	"os/signal"
	"encoding/base64"
	_ "fmt"
	"encoding/json"
	"time"
	"strings"
	"strconv"
	"github.com/blender-head/go_svc/bootstrap"
	"github.com/blender-head/go_svc/models"
)

type OrderMessage struct {
	Message []interface{} `json:"emit"`
}

//var client_id = "eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJkb3NoaWkgc2VydmVyIiwic3ViIjp7ImZvciI6IkFwcENsaWVudElkIiwiaWQiOjE1OH0sImV4cCI6MTUzMzg3NzMxOX0.TFmgVw8uNr-U21b9y66SMPkx6RzxW-kP32tQd0jzcvA";

//var socket = gowebsocket.New("wss://sandbox.doshii.co/app/socket?auth=" + base64.StdEncoding.EncodeToString([]byte(client_id)))

var client_id string
var server_url string
var socket gowebsocket.Socket
var heartbeat_Interval int

func init() {
	bootstrap.SetupLog()
	bootstrap.InitApp()
	models.InitDB() 
}

func main() {

	log.Println("Service is started")

	client_id = bootstrap.AppConfig.Client_Id
	server_url = bootstrap.AppConfig.Server_Url

	socket = gowebsocket.New(server_url + base64.StdEncoding.EncodeToString([]byte(client_id))) 

	interrupt := make(chan os.Signal, 1)
	signal.Notify(interrupt, os.Interrupt)

	socket.OnConnectError = func(err error, socket gowebsocket.Socket) {
		log.Fatalf("Received connect error - ", err)
	}
  
	socket.OnConnected = func(socket gowebsocket.Socket) {
		log.Println("Connected to server");

		heartbeat_interval := bootstrap.AppConfig.Heartbeat_Interval

		interval, err := strconv.Atoi(heartbeat_interval)

		if err != nil {
			log.Fatalf("error converting heartbeat interval to int: ", err)
		}

		ticker := time.NewTicker(time.Duration(interval) * time.Second)

		go func() {
			for t := range ticker.C {
				log.Println("Ping sent at", t)
				Heartbeat()
			}
		}()
	}
  
	socket.OnTextMessage = func(message string, socket gowebsocket.Socket) {
		is_pong_message := strings.Contains(message, "pong")
		is_checkin_created := strings.Contains(message, "checkin_created")
		is_order_created := strings.Contains(message, "order_created")
		is_order_updated := strings.Contains(message, "order_updated")
		is_transaction_created := strings.Contains(message, "transaction_created")
		is_transaction_updated := strings.Contains(message, "transaction_updated")

		if is_pong_message {
			log.Println("Received PONG message - " + message)
		}

		if is_checkin_created {
			log.Println("Received CHECK-IN CREATED message - " + message)
		}

		if is_order_created {
			log.Println("Received ORDER CREATED message - " + message)
		}

		if is_transaction_created {
			log.Println("Received TRANSACTION CREATED message - " + message)
		}

		if is_order_updated {
			log.Println("Received ORDER UPDATED message - " + message)

			order_message := OrderMessage{}

			err := json.Unmarshal([]byte(message), &order_message)
			
			if err != nil {
				log.Fatalf("error decoding order message: %s\n", err)
			}

			//fmt.Printf("%#v\n", order_message)

			order_data_message, _ := order_message.Message[1].(map[string]interface{})

			doshii_order_id := order_data_message["id"].(string)
			order_status := order_data_message["status"].(string)

			log.Println("ORDER ID - " + doshii_order_id)
			log.Println("STATUS - " + order_status)

			go func() {
				models.UpdateOrderLog(doshii_order_id, order_status)
			}()

			/*
			if order_status == "accepted" || order_status == "complete" {

				go func() {
					models.UpdateOrderLog(doshii_order_id, order_status)
				}()
			}
			*/
		}

		if is_transaction_updated {
			log.Println("Received TRANSACTION UPDATE message - " + message)
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
				log.Println("Service is stopped")
				socket.Close()
				return
		}
	}
}

func Heartbeat() {
	ping_data := make(map[string]interface{})

  	now := time.Now()
    unixtime := now.Unix()

    app_version := bootstrap.AppConfig.App_Version

  	ping_data["doshii"] = map[string]interface{}{"ping": unixtime, "version": app_version}

  	ping_data_json, err := json.Marshal(ping_data)

  	if err != nil {
  		log.Println("error decoding ping data: ", err)
  	}
    
    log.Println("Sent PING message - " + string(ping_data_json))

    socket.SendText(string(ping_data_json))
}