package game

import "fmt"

type TroopAction func(b *Board, self *Troop, x, y int) (string, error)

type TroopInfo struct {
	Name      string
	ShortName string
	Atk       int
	Secondary int
	MaxHP     int
	Mv        int
	Rng       int
	Act       TroopAction `json:"-"`
}

var TroopList = map[string]TroopInfo{
	"general": {
		Name:      "general",
		ShortName: "g",
		Atk:       2,
		Secondary: 0,
		MaxHP:     5,
		Mv:        2,
		Rng:       1,
		Act:       TroopActions["default"],
	},

	"knight": {
		Name:      "knight",
		ShortName: "k",
		Atk:       4,
		Secondary: 0,
		MaxHP:     12,
		Mv:        2,
		Rng:       1,
		Act:       TroopActions["default"],
	},

	"healer": {
		Name:      "healer",
		ShortName: "h",
		Atk:       -2,
		Secondary: 0,
		MaxHP:     8,
		Mv:        2,
		Rng:       1,
		Act:       TroopActions["default"],
	},

	"archer": {
		Name:      "archer",
		ShortName: "a",
		Atk:       2,
		Secondary: 0,
		MaxHP:     8,
		Mv:        2,
		Rng:       2,
		Act:       TroopActions["default"],
	},

	// "assassin":
	// "ghandi":

}

// returns -1 for invalid move
func CountSteps(b *Board, t *Troop, tox, toy int) int {
	//TODO obstacles
	if b.BoulderMap[tox][toy] != nil || b.TroopMap[tox][toy] != nil {
		return -1
	} else {
		dx := t.X - tox
		if dx < 0 {
			dx = -dx
		}
		dy := t.Y - toy
		if dy < 0 {
			dy = -dy
		}
		return dx + dy
	}
}

func DefaultActTroop(b *Board, t *Troop, x, y int) (string, error) {
	if !b.Inbounds(x, y) {
		return "", fmt.Errorf("Can't move out of bounds")
	}
	what := b.WhatIsAt(x, y)
	if what != -1 {
		if b.InRange(t, x, y) {
			b.Damage(x, y, t.Info.Atk)
			t.CanAct = false
		} else {
			return "", fmt.Errorf("Target out of range")
		}
	}

	//fun stuff
	hitstring := "smacked"
	if t.Info.Atk < 0 {
		hitstring = "healed"
	}
	enemystring := ""
	if what == t.Owner {
		enemystring = "his friend"
	} else if what == 2 {
		enemystring = "a rock"
	} else {
		enemystring = "a bad guy"
	}

	if x == t.X && y == t.Y {
		enemystring = "himself"
	}
	return fmt.Sprintf("%s %s %s", t.Nickname, hitstring, enemystring), nil
}

func DefaultMoveTroop(b *Board, t *Troop, x, y int) (string, error) {
	if !b.Inbounds(x, y) {
		return "", fmt.Errorf("Can't move out of bounds")
	}
	steps := CountSteps(b, t, x, y)
	if steps < 0 || steps > t.Info.Mv-t.Step {
		return "", fmt.Errorf("Can't move that far - steps left: %d", t.Info.Mv-t.Step)
	}

	switch steps {
	case 0:
		return "", fmt.Errorf("Zero step move")
	case 1:
		if b.WhatIsAt(x, y) != -1 {
			return "", fmt.Errorf("Can't move thru something")
		}
	case 2:
		dx := x - t.X
		dy := y - t.Y
		if dx == dy || dx+dy == 0 {
			xwhat := b.WhatIsAt(t.X+dx, t.Y)
			xokay := xwhat == -1 || xwhat == t.Owner
			ywhat := b.WhatIsAt(t.X, t.Y+dy)
			yokay := ywhat == -1 || ywhat == t.Owner
			if !(xokay || yokay) {
				return "", fmt.Errorf("Can't move thru something")
			}
		} else {
			what := b.WhatIsAt(t.X+(dx/2), t.Y+(dy/2))
			okay := what == -1 || what == t.Owner
			if !okay {
				return "", fmt.Errorf("Can't move thru something")
			}
		}

	}
	t.Step += steps

	b.TroopMap[t.X][t.Y] = nil
	b.TroopMap[x][y] = t
	t.X = x
	t.Y = y

	return fmt.Sprintf("%s moved %d steps", t.Nickname, steps), nil
}

var TroopActions = map[string]TroopAction{
	"default": func(b *Board, self *Troop, x, y int) (string, error) {
		if !self.CanAct {
			return "", fmt.Errorf("Troop has no more actions")
		}
		if b.IsEmpty(x, y) {
			return DefaultMoveTroop(b, self, x, y)
		} else {
			return DefaultActTroop(b, self, x, y)
		}
	},
}

type Troop struct {
	Info  TroopInfo
	Owner int

	//mutable info
	Nickname string
	HP       int
	X        int
	Y        int
	Step     int
	CanAct   bool
}

func (t *Troop) Act(b *Board, x, y int) (string, error) {
	return t.Info.Act(b, t, x, y)
}

func (t *Troop) Reset() error {
	t.CanAct = true
	t.Step = 0
	return nil
}

func gennickname(name string) string {
	return "jimmy"
}

func NewTroop(name string, x, y, owner int) (*Troop, error) {
	if _, ok := TroopList[name]; !ok {
		return nil, fmt.Errorf("No such troop %s", name)
	}

	//TODO check x and y
	//TODO check owner
	info := TroopList[name]
	return &Troop{
		Owner:    owner,
		Nickname: gennickname(name),
		Info:     info,

		HP:     info.MaxHP,
		X:      x,
		Y:      y,
		Step:   0,
		CanAct: true,
	}, nil
}

func (t *Troop) String() string {
	return fmt.Sprintf("%-10s Loc:[%2d,%2d] Stats:[%2d,%2d] Active:[%t]\n",
		t.Info.Name, t.X, t.Y, t.Info.Atk, t.HP, t.CanAct)
}
