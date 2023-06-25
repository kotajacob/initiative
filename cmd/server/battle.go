package main

import (
	"fmt"
	"log"
	"strconv"

	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/faiface/pixel/text"
	"golang.org/x/image/colornames"
)

const (
	WIDTH  = 1280
	HEIGHT = 720
)

func (c *config) battle(messages <-chan message) {
	var d display
	cfg := pixelgl.WindowConfig{
		Title:     "Initiative",
		Bounds:    pixel.R(0, 0, WIDTH, HEIGHT),
		Resizable: true,
		VSync:     true,
	}

	win, err := pixelgl.NewWindow(cfg)
	if err != nil {
		log.Fatalln(err)
	}

	face, err := loadFont()
	if err != nil {
		log.Fatalln(err)
	}

	atlas := text.NewAtlas(face, text.ASCII)
	lineheight := atlas.LineHeight() * 1.5
	txt := text.New(pixel.V(WIDTH/2, HEIGHT), atlas)
	txt.LineHeight = lineheight

	for !win.Closed() {
		win.Clear(colornames.Black)
		textheight := lineheight * float64(len(d.lines)-1)
		topmargin := (HEIGHT - textheight) / 2
		y := HEIGHT - topmargin
		x := WIDTH / 2.0
		txt = text.New(pixel.V(x, y), atlas)
		txt.LineHeight = lineheight

		for i, line := range d.lines {
			txt.Dot.X -= txt.BoundsOf(line).W() / 2
			if i == d.highlighted {
				txt.Color = colornames.Red
			} else {
				txt.Color = colornames.White
			}
			fmt.Fprintln(txt, line)
		}
		txt.Draw(win, pixel.IM)
		win.Update()

		select {
		case msg := <-messages:
			switch msg.cmd {
			case "battle":
				d.lines = msg.options
			case "highlight":
				if len(msg.options) == 0 {
					continue
				}
				i, err := strconv.Atoi(msg.options[0])
				if err != nil {
					continue
				}
				d.highlighted = i
			case "end":
				runCMD(c.End)
				win.SetClosed(true)
				win.Destroy()
			}
		default:
		}
	}
}
