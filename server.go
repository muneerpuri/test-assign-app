package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"github.com/go-redis/redis"
	"github.com/gorilla/websocket"
	"text/template"
)
const doc = `
<!DOCTYPE html>
<html lang="en">
<head>
    <meta charset="UTF-8">
    <meta name="viewport" content="width=device-width, initial-scale=1.0">
	<title>Test Table</title>
	<style>
	.table2{
	width: calc(100% - 20%);
	margin: auto;
	padding: 20px;
	text-align: center;
	border: 1px solid grey;
	}
	th, td {
	padding: 15px;
	text-align: left;
	border-bottom: 1px solid #ddd;
	}
	tr:hover {
		background-color: #f5f5f5;
		}
	tr:nth-child(even) {
		background-color: #f2f2f2;
		}
	th {
	background-color: lightblue;
	color: white;
	}
		</style>
</head>
<body>
<table>
<tr>
<th>Roll No.</th>
<th>total Words</th>
<th>Total Characters</th>
</th><th>Words/Minute</th>
</tr>

`
const doc2 = `
</body>
</html>
`
//User is
type User struct {
	Roll           string 	`json:"roll"`
	Words          string	`json:"words"`
	Characters     string	`json:"characters"`
	Wordsperminute string	`json:"wordsperminute"`
}

var user User
var user2 User
var puser User
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}
var m = map[string]string{}


func homePage(w http.ResponseWriter, r *http.Request) {
	w.Header().Add("Content Type", "text/html")
		// The template name "template" does not matter here
		templates := template.New("template")
		// "doc" is the constant that holds the HTML content
		templates.New("doc").Parse(doc)
		
		templates.Lookup("doc").Execute(w, user2)
		for key, value := range m { // Order not specified 
			json.Unmarshal([]byte(value), &user2)
			fmt.Fprintf(w, `<tr>
			<td>`+key+`</td>
			<td>"`+user2.Words+`</td>
			<td>"`+user2.Characters+`"</td>
			<td>"`+user2.Wordsperminute+`"</td>
			</tr>`)
		}
		fmt.Fprintf(w, doc2)
		
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



		client := redis.NewClient(&redis.Options{
			Addr: "localhost:6379",
			Password: "",
			DB: 0,
		})

		


		 val,err := client.Exists(finuser.Roll).Result()
		 if err != nil {
			fmt.Println(err)
		}
		// fmt.Println(val)

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
	// 	for key, value := range m { // Order not specified 
	// 	fmt.Printf( "key is"+key+"value is"+value)
	// }
		
		




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
	// log.Println("Client Connected")
	// err = ws.WriteMessage(1, []byte("Hi Client!"))
	// if err != nil {
	// 	log.Println(err)
	// }
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
