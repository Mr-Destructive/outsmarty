package outsmarty_api

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"strconv"
	"strings"

	"golang.org/x/net/websocket"
	db "outsmarty.sqlc.dev/app/outsmarty"
)

type Server struct {
	conns map[string]map[*websocket.Conn]bool
}

func NewServer() *Server {
	return &Server{
		conns: make(map[string]map[*websocket.Conn]bool),
	}
}

func (s *Server) HandleWS(w *websocket.Conn) {
	path := w.Request().URL.Path
	uid := w.Request().URL.Query().Get("uid")
	fmt.Println(uid)
	prefix := "/ws/room/"
	conn := w
	if !strings.HasPrefix(path, prefix) {
		conn.Close()
		return
	}
	roomSlug := strings.TrimPrefix(path, prefix)

	if s.conns[roomSlug] == nil {
		s.conns[roomSlug] = make(map[*websocket.Conn]bool)
	}
	s.conns[roomSlug][conn] = true
	database, err := sql.Open("sqlite3", "./outsmarty.db")
	if err != nil {
		log.Fatal(err)
	}
	defer database.Close()
	queries := db.New(database)
	userid, _ := strconv.ParseInt(uid, 10, 64)
	username, err := queries.GetUserByID(context.Background(), userid)
	if err != nil {
		log.Fatal(err)
	}

	notification := map[string]string{
		"username": "System",
		"content":  fmt.Sprintf("%s has joined the room", strings.ToTitle(username.Name)),
	}
	notificationJSON, _ := json.Marshal(notification)
	s.broadcast(notificationJSON, roomSlug)

	s.readLoop(conn, roomSlug)
}

func (s *Server) readLoop(ws *websocket.Conn, roomSlug string) {
	buf := make([]byte, 1024)
	for {
		n, err := ws.Read(buf)
		if err != nil {
			if err == io.EOF {
				delete(s.conns[roomSlug], ws)
				break
			}
			continue
		}
		msg := buf[:n]
		s.broadcast(msg, roomSlug)
	}
}

func (s *Server) broadcast(b []byte, roomSlug string) {
	for ws := range s.conns[roomSlug] {
		go func(ws *websocket.Conn) {
			if _, err := ws.Write(b); err != nil {
				fmt.Println("Write error:", err)
			}
		}(ws)
	}
}
