package game

import "fmt"

type TroopAction func(b *Board, self *Troop, x, y int) (string, error)

type TroopInfo struct {
	Name        string
	ShortName   string
	Description string
	Quote       string
	Atk         int
	Secondary   int
	MaxHP       int
	Mv          int
	Rng         int
	Act         TroopAction `json:"-"`
}

var TroopList = map[string]TroopInfo{
	"king": {
		Name:        "king",
		ShortName:   "K",
		Description: "If he dies, you lose",
		Quote:       `"If you die in the Kranch, you die in real life. -Caleb"`,
		Atk:         2,
		Secondary:   0,
		MaxHP:       5,
		Mv:          2,
		Rng:         1,
		Act:         TroopActions["default"],
	},

	"knight": {
		Name:        "knight",
		ShortName:   "n",
		Description: "Strong, well rounded soldier",
		Quote:       "",
		Atk:         4,
		Secondary:   0,
		MaxHP:       10,
		Mv:          2,
		Rng:         1,
		Act:         TroopActions["default"],
	},

	"healer": {
		Name:        "healer",
		ShortName:   "h",
		Description: "Deals negative damage",
		Quote:       "",
		Atk:         -2,
		Secondary:   0,
		MaxHP:       8,
		Mv:          2,
		Rng:         1,
		Act:         TroopActions["default"],
	},

	"ranger": {
		Name:        "ranger",
		ShortName:   "r",
		Description: "Can shoot over boulders and people",
		Quote:       `"Boom. Headshot. -jimmy"`,
		Atk:         2,
		Secondary:   0,
		MaxHP:       8,
		Mv:          2,
		Rng:         2,
		Act:         TroopActions["default"],
	},

	"cannibal": {
		Name:        "cannibal",
		ShortName:   "c",
		Description: "Heals and buffs himself when he kills something",
		Quote:       `"I must confess to you, I'm giving very serious thought... to eating your wife. -Hannibal Lecter"`,
		Atk:         2,
		Secondary:   0,
		MaxHP:       6,
		Mv:          2,
		Rng:         2,
		Act:         TroopActions["cannibal"],
	},

	"assassin": {
		Name:        "assassin",
		ShortName:   "a",
		Description: "Can backstab to do quadruple damage",
		Quote:       `"All warfare is based on deception. -Sun Tsu"`,
		Atk:         2,
		Secondary:   8,
		MaxHP:       8,
		Mv:          3,
		Rng:         2,
		Act:         TroopActions["assassin"],
	},
	// "assassin":
	// "ghandi":

}

// returns -1 for invalid move
func CountSteps(startx, starty, endx, endy int) int {
	dx := startx - endx
	if dx < 0 {
		dx = -dx
	}
	dy := starty - endy
	if dy < 0 {
		dy = -dy
	}
	return dx + dy
}

func CanMove(b *Board, owner, startx, starty, endx, endy, mv int, pass bool) bool {
	if !b.Inbounds(startx, starty) || !b.Inbounds(endx, endy) {
		return false
	}

	what := b.WhatIsAt(endx, endy)
	if what == owner && pass {
	} else if what != -1 {
		return false
	}

	if CountSteps(startx, starty, endx, endy) > mv {
		return false
	}

	dx := endx - startx
	dy := endy - starty

	if abs(dx)+abs(dy) == 1 {
		return true
	} else {
		u := CanMove(b, owner, startx, starty, startx, starty-1, 1, true) &&
			CanMove(b, owner, startx, starty-1, endx, endy, mv-1, false)
		d := CanMove(b, owner, startx, starty, startx, starty+1, 1, true) &&
			CanMove(b, owner, startx, starty+1, endx, endy, mv-1, false)
		l := CanMove(b, owner, startx, starty, startx-1, starty, 1, true) &&
			CanMove(b, owner, startx-1, starty, endx, endy, mv-1, false)
		r := CanMove(b, owner, startx, starty, startx+1, starty, 1, true) &&
			CanMove(b, owner, startx+1, starty, endx, endy, mv-1, false)
		return u || d || l || r
	}
}
func abs(input int) int {
	if input < 0 {
		return -1 * input
	}
	return input
}

func AssassinActTroop(b *Board, t *Troop, x, y int) (string, error) {
	if !b.Inbounds(x, y) {
		return "", fmt.Errorf("Can't move out of bounds")
	}

	hitstring := "sliced"

	what := b.WhatIsAt(x, y)
	if what != -1 {
		if b.InRange(t, x, y) {
			damage := t.Info.Atk
			if t.Owner == 0 {
				if t.X == x+1 && t.Y == y {
					hitstring = "backstabbed"
					damage = t.Info.Secondary
				}
			} else {
				if t.X == x-1 && t.Y == y {
					hitstring = "backstabbed"
					damage = t.Info.Secondary
				}
			}
			b.Damage(x, y, damage)
			t.CanAct = false
		} else {
			return "", fmt.Errorf("Target out of range")
		}
	}

	//fun stuff
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

	actionmessage := fmt.Sprintf("%s %s %s", t.Nickname, hitstring, enemystring)
	b.LastMessage = actionmessage
	return actionmessage, nil
}

func CannibalActTroop(b *Board, t *Troop, x, y int) (string, error) {
	if !b.Inbounds(x, y) {
		return "", fmt.Errorf("Can't move out of bounds")
	}

	hitstring := "bit"

	what := b.WhatIsAt(x, y)
	if what != -1 {
		if b.InRange(t, x, y) {
			killed, _ := b.Damage(x, y, t.Info.Atk)
			if killed {
				t.HP = t.Info.MaxHP
				t.Info.Atk += 1
				hitstring = "consumed"
			}
			t.CanAct = false
		} else {
			return "", fmt.Errorf("Target out of range")
		}
	}

	//fun stuff
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

	actionmessage := fmt.Sprintf("%s %s %s", t.Nickname, hitstring, enemystring)
	b.LastMessage = actionmessage
	return actionmessage, nil
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

	actionmessage := fmt.Sprintf("%s %s %s", t.Nickname, hitstring, enemystring)
	b.LastMessage = actionmessage
	return actionmessage, nil
}

func DefaultMoveTroop(b *Board, t *Troop, x, y int) (string, error) {
	if !b.Inbounds(x, y) {
		return "", fmt.Errorf("Can't move out of bounds")
	}

	stepsleft := t.Info.Mv - t.Step
	if !CanMove(b, t.Owner, t.X, t.Y, x, y, stepsleft, false) {
		return "", fmt.Errorf("Can't seem to get there")
	}

	steps := CountSteps(t.X, t.Y, x, y)
	t.Step += steps

	b.TroopMap[t.X][t.Y] = nil
	b.TroopMap[x][y] = t
	t.X = x
	t.Y = y

	actionmessage := fmt.Sprintf("%s moved %d steps", t.Nickname, steps)
	b.LastMessage = actionmessage
	return actionmessage, nil
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
	"cannibal": func(b *Board, self *Troop, x, y int) (string, error) {
		if !self.CanAct {
			return "", fmt.Errorf("Troop has no more actions")
		}
		if b.IsEmpty(x, y) {
			return DefaultMoveTroop(b, self, x, y)
		} else {
			return CannibalActTroop(b, self, x, y)
		}
	},
	"assassin": func(b *Board, self *Troop, x, y int) (string, error) {
		if !self.CanAct {
			return "", fmt.Errorf("Troop has no more actions")
		}
		if b.IsEmpty(x, y) {
			return DefaultMoveTroop(b, self, x, y)
		} else {
			return AssassinActTroop(b, self, x, y)
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
