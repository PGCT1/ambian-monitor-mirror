package main

import (
  "github.com/gorilla/websocket"
  "github.com/go-martini/martini"
  "github.com/martini-contrib/cors"
  "github.com/pgct1/ambian-monitor/notification"
  "github.com/pgct1/ambian-monitor/connection"
  "net/http"
)

const source = "wss://ambianmonitordev-projectgemini.rhcloud.com:8443/json"

type AuthorizationPacket struct {
	Password string
	DesiredStreams notification.MetaData
}

var sourceStream chan notification.Packet

func websocketConnectionHandler(ws *websocket.Conn) {

  //TODO mirror incoming connections to the source

}

func main(){

  sourceStream = make(chan notification.Packet)

  go InitSourceStreamConnection(sourceStream)
  go connection.InitializeConnectionManager(subscriptionPassword,sourceStream)

  martiniServerSetup := martini.Classic()

	martiniServerSetup.Use(cors.Allow(&cors.Options{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET"},
		AllowHeaders:     []string{"Origin"},
		AllowCredentials: true,
	}))

	martiniServerSetup.Get("/stream", func(w http.ResponseWriter, r *http.Request) {

		ws, err := websocket.Upgrade(w, r, nil, 1024, 1024)

		if _, ok := err.(websocket.HandshakeError); ok {

			http.Error(w, "Invalid websocket handshake", 400)

			return
		} else if err != nil {
			return
		}

		websocketConnectionHandler(ws)

	})

  martiniServerSetup.Use(martini.Static("web"))

  martiniServerSetup.Run()
}
