package handlers

import (
	"fmt"
	"log"
	"net/http"
	"sync"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool { return true },
}

type Message struct {
	Username string `json:"username"`
	Text     string `json:"text"`
	Type     string `json:"type"`
	CourseID string `json:"course_id"`
}

type Classroom struct {
	clients   map[*websocket.Conn]string
	broadcast chan Message
	mu        sync.Mutex
}

var (
	classrooms   = make(map[string]*Classroom)
	classroomsMu sync.Mutex
)

func getClassroom(courseID string) *Classroom {
	classroomsMu.Lock()
	defer classroomsMu.Unlock()
	if room, exists := classrooms[courseID]; exists {
		return room
	}
	room := &Classroom{
		clients:   make(map[*websocket.Conn]string),
		broadcast: make(chan Message, 256),
	}
	classrooms[courseID] = room
	go room.run()
	return room
}

func (room *Classroom) run() {
	for msg := range room.broadcast {
		room.mu.Lock()
		for conn := range room.clients {
			if err := conn.WriteJSON(msg); err != nil {
				conn.Close()
				delete(room.clients, conn)
			}
		}
		room.mu.Unlock()
	}
}

func ClassroomWS(c *gin.Context) {
	courseID := c.Param("courseId")
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("WebSocket upgrade hatası:", err)
		return
	}
	room := getClassroom(courseID)
	username := fmt.Sprintf("User_%v", c.GetUint("user_id"))

	room.mu.Lock()
	room.clients[conn] = username
	room.mu.Unlock()

	room.broadcast <- Message{Username: username, Text: "sınıfa katıldı", Type: "system", CourseID: courseID}

	defer func() {
		room.mu.Lock()
		delete(room.clients, conn)
		room.mu.Unlock()
		room.broadcast <- Message{Username: username, Text: "ayrıldı", Type: "system", CourseID: courseID}
		conn.Close()
	}()

	for {
		var msg Message
		if err := conn.ReadJSON(&msg); err != nil {
			break
		}
		msg.Username = username
		msg.Type = "chat"
		msg.CourseID = courseID
		room.broadcast <- msg
	}
	return
}
