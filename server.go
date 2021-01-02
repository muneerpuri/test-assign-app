package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/websocket"
	// badger "github.com/dgraph-io/badger/v2"
)

//User is
type User struct {
	Roll           int32
	Words          int32
	Characters     int32
	Wordsperminute int32
}

var user User
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}


	func homePage(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "<table><tr><th>Roll No.</th><th>total Words</th><th>Total Characters</th></th><th>Words/Minute</th></tr><tr><td>"+fmt.Sprint(user.Roll)+"</td><td>"+fmt.Sprint(user.Words)+"</td><td>"+fmt.Sprint(user.Characters)+"</td><td>"+fmt.Sprint(user.Wordsperminute)+"</td></tr></table>")
		
	}

func reader(conn *websocket.Conn) {
	for {
		// read in a message
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}
		// print out that message for clarity
		// fmt.Println(string(p))
		m := string(p)
		json.Unmarshal([]byte(m), &user)
		fmt.Printf("%+v", user)

		if err := conn.WriteMessage(messageType, p); err != nil {
			log.Println(err)
			return
		}

	}
}

func wsEndpoint(w http.ResponseWriter, r *http.Request) {
	upgrader.CheckOrigin = func(r *http.Request) bool { return true }

	// upgrade this connection to a WebSocket
	// connection
	ws, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println(err)
	}

	log.Println("Client Connected")
	err = ws.WriteMessage(1, []byte("Hi Client!"))
	if err != nil {
		log.Println(err)
	}
	// listen indefinitely for new messages coming
	// through on our WebSocket connection
	reader(ws)
}

func setupRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", wsEndpoint)
}

func main() {
	fmt.Println("Hello World")
	// db, err := badger.Open(badger.DefaultOptions("/tmp/badger"))
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer db.Close()

	setupRoutes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
