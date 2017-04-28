package main

import (
	"fmt"
	"log"

	"github.com/austindoeswork/NPC3/manager"
	"github.com/austindoeswork/NPC3/server"
)

const (
	port = ":80"
)

func main() {
	// g := game.New()
	// gs := g.GetState(0)
	// b, _ := json.Marshal(gs)
	// fmt.Printf("%s", b)

	m := manager.New()
	s := server.New(port, "./static/", m)
	fmt.Printf("blastoff @ %s\n", port)
	log.Fatal(s.Start())

	// g := game.New()
	// gs := g.GetState(0)
	// b, err := json.Marshal(gs)
	// if err != nil {
	// fmt.Println(err)
	// } else {
	// fmt.Printf("%s\n", b)
	// }

	// reader := bufio.NewReader(os.Stdin)

	// for ; g.Status != game.OVER; g.Step() {
	// currentplayer := g.ActivePlayer()
	// fmt.Printf("player %d's turn\n", currentplayer+1)
	// g.PPrint()

	// ended := false
	// for g.CanAct(currentplayer) && !ended && g.Status != game.OVER {
	// fmt.Printf("player %d >> ", currentplayer+1)
	// movestr, _ := reader.ReadString('\n')
	// movestr = movestr[0 : len(movestr)-1]

	// input := strings.Split(movestr, " ")
	// if len(input) == 1 {
	// action := input[0]
	// switch action {
	// case "end":
	// ended = true
	// case "ls":
	// for k, t := range g.Board.Troops[currentplayer] {
	// fmt.Print(fmt.Sprintf("[%d] ", k) + t.String())
	// }
	// }
	// } else if len(input) == 3 {
	// troopindex, err := strconv.Atoi(input[0])
	// tox, err := strconv.Atoi(input[1])
	// toy, err := strconv.Atoi(input[2])
	// if err != nil {
	// fmt.Println("INVALID INPUT")
	// } else {
	// err := g.MoveTroop(currentplayer, troopindex, tox, toy)
	// if err != nil {
	// fmt.Println(err)
	// } else {
	// g.PPrint()
	// }
	// }
	// }
	// }
	// }
	// g.PPrint()
}
