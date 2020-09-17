package helpers

import (
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// TypeNewUser -
	TypeNewUser = iota
	// TypeNewMessage -
	TypeNewMessage
	// TypeUserLeaveChat -
	TypeUserLeaveChat
)

// PayloadRequest -
type PayloadRequest struct {
	Message string
}

// PayloadResponse -
type PayloadResponse struct {
	Type    int    `json:"type"`
	Sender  string `json:"sender"`
	Time    int64  `json:"time"`
	Message string `json:"message"`
}

// UserConnection -
type UserConnection struct {
	*websocket.Conn
	Username string
	Room     string
}

var connections = make([]*UserConnection, 0)

func handleReadWrite(currConn *UserConnection) {
	defer func() {
		if err := recover(); err != nil {
			log.Println("ERR:", err)
		}
	}()

	broadcast(currConn, PayloadResponse{
		Message: fmt.Sprintf("%s join the chat", currConn.Username),
		Sender:  currConn.Username,
		Time:    time.Now().Unix(),
		Type:    TypeNewUser,
	})

	for {
		p := PayloadRequest{}
		err := currConn.ReadJSON(&p)
		if err != nil {
			if err != nil {
				if strings.Contains(err.Error(), "websocket: close") {
					broadcast(currConn, PayloadResponse{
						Message: fmt.Sprintf("%s has left the chat", currConn.Username),
						Sender:  currConn.Username,
						Time:    time.Now().Unix(),
						Type:    TypeUserLeaveChat,
					})
					deleteConnection(currConn)
					return
				}
				log.Println("Err:", err.Error())
				continue
			}
		}

		broadcast(currConn, PayloadResponse{
			Message: p.Message,
			Sender:  currConn.Username,
			Time:    time.Now().Unix(),
			Type:    TypeNewMessage,
		})
	}
}

func broadcast(currConn *UserConnection, pr PayloadResponse) {
	for _, conn := range connections {
		if conn.Room == currConn.Room && conn.Username != currConn.Username {
			conn.WriteJSON(pr)
		}
	}
}

func deleteConnection(currConn *UserConnection) {
	filteredConn := make([]*UserConnection, 0)
	for _, conn := range connections {
		if currConn != conn {
			filteredConn = append(filteredConn, conn)
		}
	}
	connections = filteredConn
}

// WS -
func WS(w http.ResponseWriter, r *http.Request) {
	room := r.URL.Query().Get("room")
	username := r.URL.Query().Get("username")
	if len(room) < 4 || len(username) < 4 {
		http.Error(w, "Err: Status Unpsrocessable Entity", http.StatusUnprocessableEntity)
		return
	}
	conn, err := websocket.Upgrade(w, r, w.Header(), 1024, 1024)
	if err != nil {
		http.Error(w, "Failed open ws connection", http.StatusBadRequest)
		return
	}

	currConn := UserConnection{Conn: conn, Username: username, Room: room}
	connections = append(connections, &currConn)

	go handleReadWrite(&currConn)
}
