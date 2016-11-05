/*
Package main implements a simple server with websockets.
*/
package main 

import (
	"log"
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"text/template"
	"sync"
	"time"
	"math/rand"
	"encoding/json"
	"os"
)

var (
	
	letters = "0123456789ABCDEF"
	
	connections = struct {
		sync.RWMutex
		m map[*websocket.Conn]Player
	}{ m: make(map[*websocket.Conn]Player)}

	upgrader = websocket.Upgrader {
		ReadBufferSize: 1024,
		WriteBufferSize: 1024,
	}

)

type Player struct {
	ID string
	Pos Position
}

type Position struct {
	X float32
	Y float32
	Z float32
	R_X float32
	R_Y float32
	R_Z float32
}

// 
type Message struct {
	Title string
	Data Player
}

type PlayerMessage struct {
	Title string
	Pos Position
}

func wsHandler(w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)

	if _, ok := err.(websocket.HandshakeError); ok {
		http.Error(w, "Not a websocket handshake", 400)
		return
	} else if err != nil {
		log.Println(err)
		return
	}

	log.Println("Successfully upgraded connection")
	
	connections.Lock()
	connections.m[conn] = Player{ ID: generateID() }
	connections.Unlock()
	
	go createPlayer( conn )
	
	for oldconn := range connections.m {
		
		// show other players to new player
		go addOtherPlayers( conn, oldconn )
		
		// show new player to other players
		go addNewPlayer( conn, oldconn )
		
	}

	
	for {
		_, msg, err := conn.ReadMessage()

		if err != nil {
			closeConnection( conn )
			return
		}
		
		if string(msg) == "Keep connection alive!" {
			continue
		}
		
		var result PlayerMessage
		
		err = json.Unmarshal([]byte(msg), &result)
	    if err != nil {
	        fmt.Println(err)
	        fmt.Printf("%+v\n", result)
	    }
	    
	    switch result.Title {
	    	case "updatePlayer": 
	    		connections.m[conn] = Player { ID: connections.m[conn].ID, Pos: result.Pos }
	    		
		    	for oldconn := range connections.m {
		    		go updatePlayer( conn, oldconn, connections.m[conn] )
	    		}
	    }
		
	}
}

func updatePlayer( conn *websocket.Conn, oldconn *websocket.Conn, data Player) {
	if conn != oldconn {
		message, _ := json.Marshal(&Message{ Title: "updatePlayer", Data: connections.m[conn] })
		
		if err := oldconn.WriteMessage(websocket.TextMessage, message); err != nil {
			closeConnection( oldconn )
		}
	}
}

func removePlayer( conn *websocket.Conn, oldconn *websocket.Conn ) {
	
	if conn != oldconn {
		message, _ := json.Marshal(&Message{ Title: "removePlayer", Data: connections.m[conn] })

		if err := oldconn.WriteMessage(websocket.TextMessage, message); err != nil {
			closeConnection( oldconn )
		}
	}
	
}

func createPlayer( conn *websocket.Conn ) {
	message, _ := json.Marshal(&Message{ Title: "createPlayer", Data: connections.m[conn] })

	if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
		closeConnection( conn )
	}
}

func addOtherPlayers( conn *websocket.Conn, oldconn *websocket.Conn) {
	if conn != oldconn {
		message, _ := json.Marshal(&Message{ Title: "addPlayer", Data: connections.m[oldconn] })

		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			closeConnection( conn )
		}
	}
}

func addNewPlayer( conn *websocket.Conn, oldconn *websocket.Conn) {
	if conn != oldconn {
		message, _ := json.Marshal(&Message{ Title: "addPlayer", Data: connections.m[conn] })

		if err := oldconn.WriteMessage(websocket.TextMessage, message); err != nil {
			closeConnection( oldconn )
		}
	}
}

func closeConnection( conn *websocket.Conn ) {
	for oldconn := range connections.m {
		removePlayer( conn, oldconn )
	}
	connections.Lock()
	delete(connections.m, conn)
	connections.Unlock()
	conn.Close()
}

func generateID() string {
	
	str := ""
	
	rand.Seed(time.Now().UnixNano())
	
	for i := 0; i < 6; i++ {
		num := rand.Intn(5)
		str += letters[num: num + 1]
	}
	
	return str
}

func main() {
	port := GetPort()

	http.HandleFunc("/", serveIndex)
	http.HandleFunc("/ws", wsHandler)
	http.HandleFunc("/css/", serveStatic)
	http.HandleFunc("/js/", serveStatic)
	http.HandleFunc("/images/", serveStatic)

	log.Printf("Running on port" + port)

	err := http.ListenAndServe(port, nil)
	log.Println(err.Error())
}

func serveIndex(w http.ResponseWriter, r *http.Request) {
	template.Must(template.ParseFiles("html/client.html")).Execute(w, nil)
}

func serveStatic(w http.ResponseWriter, r *http.Request) {
	template.Must(template.ParseFiles(r.URL.Path[1:])).Execute(w, nil)
}

func GetPort() string {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8081"
		log.Println("[-] No PORT environment variable detected. Setting to ", port)
	}
	return ":" + port
}
