package main

import (
  "github.com/gorilla/websocket"
  "github.com/pgct1/ambian-monitor/notification"
  "fmt"
  "time"
  "crypto/tls"
  "net/url"
  "net/http"
  "encoding/json"
)

func InitSourceStreamConnection(stream chan notification.Packet){

  for {
    connectToSource(stream)  // returns on connection error
    time.Sleep(10 * time.Second)
  }

}

func connectToSource(stream chan notification.Packet){

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

  sources := notification.Sources{true,true,true}

  metaData := notification.MetaData{
    AmbianStreamIds:streamIds,
    Sources:sources,
  }

  authorizationPacket := AuthorizationPacket{
    Password:sourcePassword,
    DesiredStreams:metaData,
  }

  ws.SetWriteDeadline(time.Now().Add(2 * time.Second))

  err = ws.WriteJSON(authorizationPacket)

  if err != nil {
    fmt.Println("ERROR")
    fmt.Println(err)
    return
  }

  // connection successful; now push notifications through the stream

  for {

    var packet notification.Packet

    _, jsonNotification, err := ws.ReadMessage()

    if err == nil {
      err = json.Unmarshal(jsonNotification,&packet)

      if err == nil {

        stream <- packet

      }else{
        fmt.Println("Unable to parse JSON.")
      }

    }else{
      fmt.Println("Connection error. Closing source connection.")
      fmt.Println(err)
      break
    }

  }

}
