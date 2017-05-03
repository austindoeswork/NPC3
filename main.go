package main

import (
	"flag"
	"fmt"
	"log"
	"math/rand"
	"time"

	"github.com/austindoeswork/NPC3/manager"
	"github.com/austindoeswork/NPC3/server"
)

const (
	port = ":80"
)

var (
	versionFlag = flag.Bool("v", false, "git commit hash")
	commithash  string
)

func main() {
	flag.Parse()
	if *versionFlag {
		fmt.Println(commithash)
		return
	}
	rand.Seed(time.Now().UTC().UnixNano())

	m := manager.New()
	s := server.New(port, "./static/", m)

	fmt.Println("version: " + commithash)
	fmt.Printf("blastoff @ %s\n", port)
	log.Fatal(s.Start())
}
