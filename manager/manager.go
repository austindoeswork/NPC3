// The manager handles creating games and I/O
package manager

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/austindoeswork/NPC3/game"
)

//TODO game uuid
type Player struct {
	Input  chan []byte
	Output chan []byte
	Index  int
}

type Listener struct {
	Output chan []byte
	Index  int
}

func NewPlayer(input chan []byte, output chan []byte, index int) *Player {
	return &Player{
		Input:  input,
		Output: output,
		Index:  index,
	}
}

func NewListener(output chan []byte, index int) *Listener {
	return &Listener{
		Output: output,
		Index:  index,
	}
}

type Manager struct {
	GameMap       map[string]*game.Game
	PlayerMap     map[string][]*Player
	GameListeners map[string][]*Listener
}

func New() *Manager {
	return &Manager{
		GameMap:       map[string]*game.Game{},
		PlayerMap:     map[string][]*Player{},
		GameListeners: map[string][]*Listener{},
	}
}

func (m *Manager) Exists(name string) bool {
	_, exists := m.GameMap[name]
	return exists
}

func (m *Manager) AddGame(name string) error {
	if m.Exists(name) {
		return fmt.Errorf("Game already exists")
	}

	g := game.New()
	m.GameMap[name] = g
	m.PlayerMap[name] = []*Player{}

	go m.HandleGameOutput(name)
	return nil
}

type GameSummary struct {
	Name    string
	Players int
}

func (m *Manager) ListGames() []*GameSummary {
	arr := []*GameSummary{}
	for k, v := range m.PlayerMap {
		arr = append(arr, &GameSummary{k, len(v)})
	}
	return arr
}

type Ack struct {
	Type    string
	Success bool
	Message string
}

type Prompt struct {
	Type    string
	Message string
}

type State struct {
	Type string
	*game.GameState
}

//OUTPUT
func (m *Manager) HandleGameOutput(gamename string) {
	if !m.Exists(gamename) {
		return
	}
	update := m.GameMap[gamename].Update

	for msg := range update {
		fmt.Println("GOT MESSAGE", msg)
		switch msg {
		case -1: //game over
			for _, l := range m.GameListeners[gamename] {
				close(l.Output)
			}
			//TODO mutex
			delete(m.GameMap, gamename)
			delete(m.PlayerMap, gamename)
			delete(m.GameListeners, gamename)
		case 0, 1: //player turn
			prompt := &Prompt{
				Type:    "PROMPT",
				Message: fmt.Sprintf("YOUR TURN"),
			}
			b, err := json.Marshal(prompt)
			if err != nil {
				log.Println(err)
			}
			for _, l := range m.GameListeners[gamename] {
				if l.Index == msg {
					SendBytesOrTimeout(l.Output, b, 1)
				}
			}
		case 2: //game state updated
			for _, l := range m.GameListeners[gamename] {
				gstate := m.GameMap[gamename].GetState(l.Index)
				state := &State{
					"STATE",
					gstate,
				}
				b, err := json.Marshal(state)
				if err != nil {
					log.Println(err)
				}
				SendBytesOrTimeout(l.Output, b, 1)
			}
		}
	}
}

type Command struct {
	Type string

	Troop int
	X     int
	Y     int

	Message string
}

//INPUT
func (m *Manager) HandleCommand(gamename string, player int, command []byte) ([]byte, error) {
	var cmd Command
	err := json.Unmarshal(command, &cmd)
	if err != nil {
		return []byte{}, err
	}
	if len(m.PlayerMap[gamename]) < 2 {
		res := &Ack{
			Type:    "ACK",
			Success: true,
			Message: "Waiting for another player",
		}
		b, _ := json.Marshal(res)
		return b, nil
	}

	g := m.GameMap[gamename]
	var b []byte

	switch cmd.Type {
	case "ECHO":
		res := &Ack{
			Type:    "ACK",
			Success: true,
			Message: cmd.Message,
		}
		b, _ = json.Marshal(res)
	case "STATE":
		res := &State{
			Type:      "STATE",
			GameState: m.GameMap[gamename].GetState(player),
		}
		b, _ = json.Marshal(res)
	case "MOVE":
		movestring, err := m.GameMap[gamename].MoveTroop(player, cmd.Troop, cmd.X, cmd.Y)
		res := &Ack{
			Type: "ACK",
		}
		if err != nil {
			res.Message = err.Error()
			res.Success = false
		} else {
			res.Success = true
			res.Message = movestring
		}
		b, _ = json.Marshal(res)
	case "END":
		res := &Ack{
			Type: "ACK",
		}
		if m.GameMap[gamename].ActivePlayer() == player {
			res.Success = true
			res.Message = "Thanks."
			g.Step()
		} else {
			res.Success = false
			res.Message = "It's not ur turn bro..."
		}
		b, _ = json.Marshal(res)
	default:
		res := &Ack{
			Type:    "ACK",
			Success: false,
			Message: "INVALID TYPE",
		}
		b, _ = json.Marshal(res)
	}

	// m.GameMap[gamename].Step()
	return b, nil
}

func (m *Manager) AddPlayer(gameinput chan []byte, gameoutput chan []byte, gamename string) (int, error) {
	if !m.Exists(gamename) {
		return -1, fmt.Errorf("Game %s does not exist", gamename)
	}
	if len(m.PlayerMap[gamename]) >= 2 {
		return -1, fmt.Errorf("Too many players")
	}

	playerindex := len(m.PlayerMap[gamename])

	p := NewPlayer(gameinput, gameoutput, playerindex)
	if p == nil {
		return -1, fmt.Errorf("Could not create player")
	}

	l := NewListener(gameoutput, playerindex)
	if l == nil {
		return -1, fmt.Errorf("Could not create listener")
	}

	m.PlayerMap[gamename] = append(m.PlayerMap[gamename], p)
	m.GameListeners[gamename] = append(m.GameListeners[gamename], l)

	// LISTEN TO INPUT
	go func() {
		for msg := range gameinput {
			if !m.Exists(gamename) {
				return
			}
			res, err := m.HandleCommand(gamename, p.Index, msg)
			if err != nil {
				log.Println(err)
			} else {
				fmt.Println("C")
				SendBytesOrTimeout(gameoutput, res, 1)
			}
		}
		// log.Println("ENDED")
	}()

	if playerindex == 1 {
		go func() {
			for _, l := range m.GameListeners[gamename] {
				gstate := m.GameMap[gamename].GetState(l.Index)
				state := &State{
					"STATE",
					gstate,
				}
				b, err := json.Marshal(state)
				if err != nil {
					log.Println(err)
				}
				fmt.Println("D")
				SendBytesOrTimeout(l.Output, b, 1)
			}

			prompt := &Prompt{
				Type:    "PROMPT",
				Message: fmt.Sprintf("YOUR TURN"),
			}
			b, err := json.Marshal(prompt)
			if err != nil {
				log.Println(err)
			}
			fmt.Println("E")
			SendBytesOrTimeout(m.PlayerMap[gamename][0].Output, b, 1)
		}()
	}

	return playerindex, nil
}

func SendBytesOrTimeout(msgchan chan []byte, msg []byte, seconds int) {
	timeout := make(chan struct{})
	go func() {
		time.Sleep(time.Duration(seconds) * time.Second)
		timeout <- struct{}{}
	}()
	go func() {
		defer func() {
			if r := recover(); r != nil {
				fmt.Println("Recovered in SendBytesOrTimeout", r)
			}
		}()
		select {
		case msgchan <- msg:
		case <-timeout:
		}
	}()
}
