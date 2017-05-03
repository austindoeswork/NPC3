package server

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strings"
	"time"

	"github.com/austindoeswork/NPC3/game"
	"github.com/austindoeswork/NPC3/manager"
	"github.com/gorilla/websocket"
)

type Server struct {
	port     string
	static   string
	mux      *http.ServeMux
	gmanager *manager.Manager
}

func (s *Server) AddHandleFunc(path string, handle func(http.ResponseWriter, *http.Request)) {
	s.mux.HandleFunc(path, handle)
}

func (s *Server) AddStaticHandler(dirpath string) {
	s.mux.Handle("/", http.FileServer(http.Dir(dirpath)))
}

func New(port string, static string, gm *manager.Manager) *Server {
	s := &Server{
		port:     port,
		static:   static,
		mux:      http.NewServeMux(),
		gmanager: gm,
	}
	s.AddStaticHandler(s.static)

	s.AddHandleFunc("/ws", s.JoinWS)
	s.AddHandleFunc("/wswatch", s.WatchWS)
	s.AddHandleFunc("/games", s.ListGames)

	return s
}

func (s *Server) Start() error {
	return http.ListenAndServe(s.port, s.mux)
}

func (s *Server) ListGames(w http.ResponseWriter, r *http.Request) {
	arr := s.gmanager.ListGames()
	b, err := json.Marshal(arr)
	if err != nil {
		log.Println(err)
		return
	}
	w.Write(b)
}

//TODO eliminate redundancy here and in manager
type Ack struct {
	Type    string
	Success bool
	Message string
}

type Prompt struct {
	Type    string
	Message string
}

func cleanstring(input string) string {
	if input == "doge" {
		return `<img src="http://i0.kym-cdn.com/entries/icons/facebook/000/013/564/aP2dv.jpg">`
	}
	invalid := []string{"<", ">", `/`, `\`, `"`, `'`, `;`}
	output := input
	for _, v := range invalid {
		output = strings.Replace(output, v, "?", -1)
	}
	return output
}

type State struct {
	Type string
	*game.GameState
}

func (s *Server) WatchWS(w http.ResponseWriter, r *http.Request) {
	gamename := r.FormValue("game")
	if len(gamename) <= 0 {
		w.Write([]byte("No gamename provided"))
		return
	}
	gamename = cleanstring(gamename)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Cannot upgrade websocket", err)
		return
	}

	gameoutput := make(chan []byte)

	// ADD WATCHER
	err = s.gmanager.AddWatcher(gameoutput, gamename)
	if err != nil {
		ack := &Ack{
			Type:    "ACK",
			Success: false,
			Message: err.Error(),
		}
		conn.WriteJSON(ack)
		log.Println(err)
		return
	}

	gs := s.gmanager.GameMap[gamename].GetState(-1)
	st := &State{
		Type:      "STATE",
		GameState: gs,
	}
	conn.WriteJSON(st)

	ack := &Ack{
		Type:    "ACK",
		Success: true,
		Message: fmt.Sprintf("Watching %s", gamename),
	}

	err = conn.WriteJSON(ack)
	if err != nil {
		fmt.Println(err)
	}

	//OUTPUT
	go func() {
		defer conn.Close()
		for msg := range gameoutput {
			err = conn.SetWriteDeadline(time.Now().Add(1 * time.Minute))
			if err != nil {
				log.Println(err)
				return
			}
			fmt.Printf("[w] %6.6s >> %-20.20s\n", gamename, msg)
			conn.WriteMessage(1, msg)
		}
	}()
}

func (s *Server) JoinWS(w http.ResponseWriter, r *http.Request) {
	gamename := r.FormValue("game")
	if len(gamename) <= 0 {
		w.Write([]byte("No gamename provided"))
		return
	}
	gamename = cleanstring(gamename)
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Println("Cannot upgrade websocket", err)
		return
	}

	defer conn.Close()

	prompt := &Prompt{
		Type:    "PROMPT",
		Message: "SEND DEVKEY",
	}
	conn.WriteJSON(prompt)

	// ADD GAME
	if !s.gmanager.Exists(gamename) {
		err := s.gmanager.AddGame(gamename)
		if err != nil {
			//TODO send err to WS
			log.Println(err)
			return
		}
	}

	// CREATE I/O CHANS
	gameinput := make(chan []byte)
	defer close(gameinput)
	gameoutput := make(chan []byte)

	// ADD PLAYER
	index, err := s.gmanager.AddPlayer(gameinput, gameoutput, gamename)
	if err != nil {
		ack := &Ack{
			Type:    "ACK",
			Success: false,
			Message: err.Error(),
		}
		conn.WriteJSON(ack)
		log.Println(err)
		return
	}

	ack := &Ack{
		Type:    "ACK",
		Success: true,
		Message: fmt.Sprintf("%d %s", index, gamename),
	}
	conn.WriteJSON(ack)

	//OUTPUT
	go func() {
		defer conn.Close()
		for msg := range gameoutput {
			err = conn.SetReadDeadline(time.Now().Add(5 * time.Minute))
			if err != nil {
				log.Println(err)
				return
			}
			fmt.Printf("[%d] %6.6s >> %-20.20s\n", index, gamename, msg)
			conn.WriteMessage(1, msg)
		}
	}()

	//INPUT
	for {
		mt, message, err := conn.ReadMessage()
		if err != nil || mt == CloseMessage {
			log.Println("Closing websocket", err)
			return
		}
		err = conn.SetReadDeadline(time.Now().Add(5 * time.Minute))
		if err != nil {
			log.Println(err)
			return
		}

		fmt.Printf("[%d] %6.6s << %-20.20s\n", index, gamename, message)

		timeout := make(chan struct{})
		go func() {
			time.Sleep(1 * time.Second)
			timeout <- struct{}{}
		}()

		select {
		case gameinput <- message:
		case <-timeout:
			return
		}
		// SendBytesOrTimeout(gameinput, message, 1)
	}
}

// Upgrades a regular ResponseWriter to WebSocketResponseWriter
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// from https://godoc.org/github.com/gorilla/websocket
const (
	// TextMessage denotes a text data message. The text message payload is
	// interpreted as UTF-8 encoded text data.
	TextMessage = 1

	// BinaryMessage denotes a binary data message.
	BinaryMessage = 2

	// CloseMessage denotes a close control message. The optional message
	// payload contains a numeric code and text. Use the FormatCloseMessage
	// function to format a close message payload.
	CloseMessage = 8

	// PingMessage denotes a ping control message. The optional message payload
	// is UTF-8 encoded text.
	PingMessage = 9

	// PongMessage denotes a ping control message. The optional message payload
	// is UTF-8 encoded text.
	PongMessage = 10
)
