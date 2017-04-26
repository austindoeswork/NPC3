package game

type Boulder struct {
	X     int
	Y     int
	MaxHP int
	HP    int
}

func NewBoulder(x, y int) *Boulder {
	return &Boulder{
		X:     x,
		Y:     y,
		MaxHP: 10,
		HP:    10,
	}
}
