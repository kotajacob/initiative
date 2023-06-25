package main

import (
	"log"
	"os"
	"os/exec"
	"strings"

	"github.com/faiface/pixel/pixelgl"
)

func main() {
	pixelgl.Run(run)
}

type display struct {
	highlighted int
	lines       []string
}

func run() {
	c, err := loadConfig()
	if err != nil {
		log.Fatalln(err)
	}
	messages := make(chan message)
	go server(messages)

	for {
		select {
		case msg := <-messages:
			if msg.cmd == "start" {
				os.Setenv("DISPLAY", c.Display)
				runCMD(c.Start)
				c.battle(messages)
			}
		default:
		}
	}
}

func runCMD(s string) {
	parts := strings.Split(strings.TrimSpace(s), " ")
	if len(parts) == 0 {
		return
	}
	cmd := exec.Command(parts[0], parts[1:]...)
	err := cmd.Run()
	if err != nil {
		log.Fatalln(err)
	}
}
