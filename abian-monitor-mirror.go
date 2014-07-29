package main

import (
  "fmt"
  "github.com/gorilla/websocket"
  "github.com/go-martini/martini"
  "github.com/martini-contrib/cors"
  "crypto/tls"
  "net/http"
	"net/url"
  "time"
)

const source = "wss://ambianmonitordev-projectgemini.rhcloud.com:8443/json"


type Sources struct {
	Corporate bool
	SocialMedia bool
	Aggregate bool
}

type NotificationMetaData struct {
	AmbianStreamIds []int
	Sources Sources
}

const (
	cNotificationTypeDefault = iota
	cNotificationTypeTweet = iota
	cNotificationTypeOfficialNews = iota
)

// packets

type NotificationPacket struct {
	Type int
	Content string
	MetaData NotificationMetaData
}

type AuthorizationPacket struct {
	Password string
	DesiredStreams NotificationMetaData
}

func websocketConnectionHandler(ws *websocket.Conn) {

}

func main(){

  //sourceStream := make(chan NotificationPacket)


	streamUrl, _ := url.Parse(source)

  connection, err := tls.Dial("tcp", streamUrl.Host, &tls.Config{})

  if err != nil {
    fmt.Println(err)
    return
	}

	ws, response, err := websocket.NewClient(connection, streamUrl, http.Header{"Origin": {"ambianmonitordev-projectgemini.rhcloud.com"}}, 1024, 1024)
  
  if err != nil {
    fmt.Println(err)
    return
  }

  if response.StatusCode > 500 {
    fmt.Println("500 error received when trying to connect.")
    return
  }

	defer ws.Close()

  // send authorization packet

  streamIds := make([]int,0,1)

  streamIds = append(streamIds,1)

  sources := Sources{true,true,true}

  metaData := NotificationMetaData{
    AmbianStreamIds:streamIds,
    Sources:sources,
  }

  authorizationPacket := AuthorizationPacket{
    Password:SubscriptionPassword,
    DesiredStreams:metaData,
  }

	ws.SetWriteDeadline(time.Now().Add(2 * time.Second))

	err = ws.WriteJSON(authorizationPacket)

	if err != nil {
    fmt.Println("ERROR")
    fmt.Println(err)
    return
	}

  fmt.Println("nice")

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
