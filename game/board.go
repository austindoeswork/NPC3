package game

import (
	"fmt"
)

type Board struct {
	Width       int
	Height      int
	LastMessage string

	Boulders []*Boulder
	Troops   [][]*Troop

	TroopMap   [][]*Troop
	BoulderMap [][]*Boulder
}

func (b *Board) Inbounds(x, y int) bool {
	return x >= 0 && x < b.Width && y >= 0 && y < b.Height
}

func (b *Board) InRange(t *Troop, x, y int) bool {
	dx := t.X - x
	if dx < 0 {
		dx = -dx
	}
	dy := t.Y - y
	if dy < 0 {
		dy = -dy
	}

	return dx <= t.Info.Rng && dy <= t.Info.Rng
}

// -1 nothing 0 p0troop 1 p1troop 2 boulder
func (b *Board) WhatIsAt(x, y int) int {
	if !b.Inbounds(x, y) {
		return -1
	} else if b.BoulderMap[x][y] != nil {
		return 2
	} else if b.TroopMap[x][y] != nil {
		return b.TroopMap[x][y].Owner
	} else {
		return -1
	}
}

func (b *Board) IsEmpty(x, y int) bool {
	if !b.Inbounds(x, y) {
		return true
	}

	if b.TroopMap[x][y] == nil && b.BoulderMap[x][y] == nil {
		return true
	}
	return false
}

//returns true if the damaged thing got killed
func (b *Board) Damage(x, y, dmg int) (bool, error) {
	killed := false

	if !b.Inbounds(x, y) {
		return false, fmt.Errorf("Out of bounds")
	}
	what := b.WhatIsAt(x, y)
	if what == -1 {
		return false, fmt.Errorf("Nothing to damage at x:%d y:%d", x, y)
	} else if what == 2 {
		boul := b.BoulderMap[x][y]
		if boul == nil {
			return false, fmt.Errorf("No boulder")
		}
		boul.HP -= dmg
		if boul.HP > boul.MaxHP {
			boul.HP = boul.MaxHP
		}
		if boul.HP <= 0 {
			killed = true
			b.BoulderMap[x][y] = nil
			index := -1
			for i := 0; i < len(b.Boulders); i++ {
				if b.Boulders[i] == boul {
					index = i
					break
				}
			}
			if index != -1 {
				b.Boulders = append(b.Boulders[:index], b.Boulders[index+1:]...)
			}
		}
	} else {
		t := b.TroopMap[x][y]
		if t == nil {
			return false, fmt.Errorf("No troop")
		}
		t.HP -= dmg
		if t.HP > t.Info.MaxHP {
			t.HP = t.Info.MaxHP
		}
		if t.HP <= 0 {
			killed = true
			index := -1
			for i := 0; i < len(b.Troops[t.Owner]); i++ {
				if b.Troops[t.Owner][i] == t {
					index = i
					break
				}
			}
			if index != -1 {
				if index != 0 {
					b.TroopMap[x][y] = nil
					//delete the unit
					b.Troops[t.Owner] = append(b.Troops[t.Owner][:index], b.Troops[t.Owner][index+1:]...)
				}
			}
		}
	}

	return killed, nil
}
