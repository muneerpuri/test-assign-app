package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
)

//User is
type User struct {
	Roll           string 	`json:"roll"`
	Words          string	`json:"words"`
	Characters     string	`json:"characters"`
	Wordsperminute string	`json:"wordsperminute"`
}
var client = redis.NewClient(&redis.Options{
	Addr: "localhost:6379",
	Password: "",
	DB: 0,
})
var user User
var user2 User
var puser User
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
var m = map[string]string{}

func homePage(w http.ResponseWriter, r *http.Request) {
	
		
}

func reader(conn *websocket.Conn)	{
	for {
		messageType, p, err := conn.ReadMessage()
		if err != nil {
			log.Println(err)
			return
		}

		res1 := strings.Split(string(p), ",") 
		user.Roll = res1[0]
		user.Words = res1[1]
		user.Characters = res1[2]
		user.Wordsperminute = res1[3]
		

		finuser := &User{Roll: user.Roll,Words:user.Words,Characters: user.Characters,Wordsperminute:user.Wordsperminute }
		e, err := json.Marshal(finuser)
		if err != nil {
			fmt.Println(err)
		}

		 val,err := client.Exists(finuser.Roll).Result()
		 if err != nil {
			fmt.Println(err)
		}

		if val != 0 {
			err = client.Set(finuser.Roll, e, 0).Err()
			if err != nil {
				fmt.Println(err)
			}
		}else{
			err = client.Set(finuser.Roll, e, 0).Err()
			if err != nil {
				fmt.Println(err)
			}
		}
		keys,err := client.Keys("*").Result()
		if err != nil {
			fmt.Println(err)
		}
		// fmt.Println(keys)
		for _,key := range keys {
			val, err := client.Get(fmt.Sprint(key)).Result()
			if err != nil {

				fmt.Println(err)
			}
			// fmt.Println(val)
			m[key] = val



		}
		
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
	defer ws.Close()
	for _, value := range m { // Order not specified 
		err = ws.WriteMessage(1, []byte(value))

    if err != nil {
        log.Println(err)
    }
	}
	
	
	
	reader(ws)

}

func setupRoutes() {
	http.HandleFunc("/", homePage)
	http.HandleFunc("/ws", wsEndpoint)
}

func main() {
	fmt.Println("Hello World")
	

	setupRoutes()
	log.Fatal(http.ListenAndServe(":8080", nil))
}
