package game

import (
	"fmt"
	"strings"
)

const (
	OVER    = -1
	RUNNING = 1
)

type Game struct {
	Board  *Board
	Turn   int
	Status int
	Winner int
}

func (g *Game) AddTroop(t *Troop) error {
	if t == nil {
		return fmt.Errorf("Troop is nil")
	}
	if !g.Board.Inbounds(t.X, t.Y) {
		return fmt.Errorf("Troop is out of bounds")
	}
	if old := g.Board.TroopMap[t.X][t.Y]; old != nil {
		return fmt.Errorf("Troop already exists at x:%d y:%d", t.X, t.Y)
	}
	g.Board.TroopMap[t.X][t.Y] = t
	g.Board.Troops[t.Owner] = append(g.Board.Troops[t.Owner], t)
	return nil
}

func (g *Game) AddBoulder(b *Boulder) error {
	if b == nil {
		return fmt.Errorf("Boulder is nil")
	}
	if !g.Board.Inbounds(b.X, b.Y) {
		return fmt.Errorf("Boulder is out of bounds")
	}
	if old := g.Board.BoulderMap[b.X][b.Y]; old != nil {
		return fmt.Errorf("Boulder already exists at x:%d y:%d", b.X, b.Y)
	}
	g.Board.BoulderMap[b.X][b.Y] = b
	g.Board.Boulders = append(g.Board.Boulders, b)
	return nil
}

func (g *Game) MoveTroop(player, troopindex, x, y int) error {
	if player < 0 || player > 1 {
		return fmt.Errorf("No such player %d", player)
	}
	if troopindex < 0 || troopindex >= len(g.Board.Troops[player]) {
		return fmt.Errorf("No such troop %d", troopindex)
	}

	t := g.Board.Troops[player][troopindex]
	if t == nil {
		return fmt.Errorf("Nil troop at player: %d troopindex: %d", player, troopindex)
	}
	err := t.Act(g.Board, x, y)

	if err == nil {
		hp0 := g.Board.Troops[0][0].HP
		hp1 := g.Board.Troops[1][0].HP
		if hp0 <= 0 {
			g.Status = OVER
			g.Winner = 1
		} else if hp1 <= 0 {
			g.Status = OVER
			g.Winner = 0
		}
	}

	return err
}

func (g *Game) ActivePlayer() int {
	return g.Turn % 2
}

func (g *Game) CanAct(player int) bool {
	if player < 0 || player > 1 {
		return false
	}
	if player != g.ActivePlayer() {
		return false
	}
	for i := 0; i < len(g.Board.Troops[player]); i++ {
		if g.Board.Troops[player][i].CanAct {
			return true
		}
	}
	return false
}

func (g *Game) Step() {
	nextactiveplayer := (g.Turn + 1) % 2
	for i := 0; i < len(g.Board.Troops[nextactiveplayer]); i++ {
		g.Board.Troops[nextactiveplayer][i].Reset()
	}
	g.Turn = g.Turn + 1
}

func New() *Game {
	g := &Game{
		Board: &Board{
			Width: 10,
			// Width:    20,
			Height: 5,
			// Height:   15,
			Boulders: []*Boulder{},
			Troops:   [][]*Troop{[]*Troop{}, []*Troop{}},
			//			 x y
			TroopMap:   [][]*Troop{},
			BoulderMap: [][]*Boulder{},
		},
		Turn:   0,
		Status: RUNNING,
	}

	//initialize troop and boulder maps
	for i := 0; i < g.Board.Width; i++ {
		tarr := []*Troop{}
		barr := []*Boulder{}
		for j := 0; j < g.Board.Height; j++ {
			tarr = append(tarr, nil)
			barr = append(barr, nil)
		}
		g.Board.TroopMap = append(g.Board.TroopMap, tarr)
		g.Board.BoulderMap = append(g.Board.BoulderMap, barr)
	}

	// TODO rando gen boulders
	// for i := 7; i < 15; i++ {
	// b := NewBoulder(i, 8)
	// g.AddBoulder(b)
	// }
	p0general, _ := NewTroop("general", 0, 2, 0)
	p0knight, _ := NewTroop("knight", 1, 3, 0)

	p1general, _ := NewTroop("general", 9, 2, 1)
	p1knight, _ := NewTroop("knight", 8, 2, 1)

	b0 := NewBoulder(0, 1)
	b1 := NewBoulder(1, 2)

	g.AddBoulder(b0)
	g.AddBoulder(b1)

	g.AddTroop(p0general)
	g.AddTroop(p0knight)
	g.AddTroop(p1general)
	g.AddTroop(p1knight)

	return g
}

func (g *Game) PPrint() {
	b := [][]byte{}
	info := []string{}

	for i := 0; i < g.Board.Height+2; i++ {
		inner := []byte{}
		for j := 0; j < g.Board.Width+2; j++ {
			info = append(info, "")
			if i == 0 || i == g.Board.Height+1 {
				if j == 0 || j == g.Board.Width+1 {
					inner = append(inner, '+')
				} else {
					if j%5 == 0 {
						inner = append(inner, '|')
					} else {
						inner = append(inner, '-')
					}
				}
			} else if j == 0 || j == g.Board.Width+1 {
				if i%5 == 0 {
					inner = append(inner, '-')
				} else {
					inner = append(inner, '|')
				}
			} else {
				if (i+j)%2 == 0 {
					inner = append(inner, '.')
				} else {
					inner = append(inner, ' ')
				}
			}
		}
		b = append(b, inner)
	}
	for i := 0; i < len(g.Board.Boulders); i++ {
		boul := g.Board.Boulders[i]
		b[boul.Y+1][boul.X+1] = '*'
	}
	for i := 0; i < len(g.Board.TroopMap); i++ {
		for j := 0; j < len(g.Board.TroopMap[i]); j++ {
			t := g.Board.TroopMap[i][j]
			if t != nil {
				str := t.Info.ShortName
				if t.Owner == 1 {
					str = strings.ToUpper(str)
				}
				b[j+1][i+1] = []byte(str)[0]
				info[j+1] += fmt.Sprintf("%s(%d,%d)[%d,%d] ", str, t.X, t.Y, t.Info.Atk, t.HP)
			}
		}
	}
	for i := 0; i < g.Board.Height+2; i++ {
		fmt.Printf("%s %s\n", b[i], info[i])
	}
}

type BoardState struct {
	Width  int
	Height int
}

type GameState struct {
	You      int
	Opponent int
	Turn     int
	Status   int
	Winner   int

	Board    BoardState
	Troops   [][]*Troop
	Boulders []*Boulder
}

func (g *Game) GetState(player int) *GameState {
	opponent := -1
	if player == 0 {
		opponent = 1
	} else {
		opponent = 0
	}

	gs := &GameState{
		You:      player,
		Opponent: opponent,
		Turn:     g.Turn,
		Status:   g.Status,
		Winner:   g.Winner,
		Board: BoardState{
			Width:  g.Board.Width,
			Height: g.Board.Height,
		},
		Troops:   g.Board.Troops,
		Boulders: g.Board.Boulders,
	}
	return gs
}
