package main

import (
	"github.com/hajimehoshi/go-ebiten"
	"github.com/hajimehoshi/go-ebiten/example/game/blank"
	"github.com/hajimehoshi/go-ebiten/example/game/input"
	"github.com/hajimehoshi/go-ebiten/example/game/monochrome"
	"github.com/hajimehoshi/go-ebiten/example/game/rects"
	"github.com/hajimehoshi/go-ebiten/example/game/rotating"
	"github.com/hajimehoshi/go-ebiten/example/game/sprites"
	"github.com/hajimehoshi/go-ebiten/example/game/testpattern"
	"github.com/hajimehoshi/go-ebiten/graphics"
	"github.com/hajimehoshi/go-ebiten/ui/cocoa"
	"os"
	"runtime"
	"time"
)

type Game interface {
	InitTextures(tf graphics.TextureFactory)
	Update()
	Draw(canvas graphics.Canvas)
}

func main() {
	runtime.GOMAXPROCS(runtime.NumCPU())

	gameName := ""
	if 2 <= len(os.Args) {
		gameName = os.Args[1]
	}

	var game Game
	switch gameName {
	case "blank":
		game = blank.New()
	case "input":
		game = input.New()
	case "monochrome":
		game = monochrome.New()
	case "rects":
		game = rects.New()
	default:
		fallthrough
	case "rotating":
		game = rotating.New()
	case "sprites":
		game = sprites.New()
	case "testpattern":
		game = testpattern.New()
	}

	const screenWidth = 256
	const screenHeight = 240
	const screenScale = 2
	const fps = 60
	const title = "Ebiten Demo"

	var ui ebiten.UI = cocoa.New(screenWidth, screenHeight, screenScale, title)
	ui.InitTextures(game.InitTextures)

	drawing := make(chan *graphics.LazyCanvas)
	go func() {
		inputStateUpdated := ui.ObserveInputStateUpdated()
		screenSizeUpdated := ui.ObserveScreenSizeUpdated()

		frameTime := time.Duration(int64(time.Second) / int64(fps))
		tick := time.Tick(frameTime)
		for {
			select {
			case e, ok := <-inputStateUpdated:
				if ok {
					type Handler interface {
						OnInputStateUpdated(ebiten.InputStateUpdatedEvent)
					}
					if game2, ok := game.(Handler); ok {
						game2.OnInputStateUpdated(e)
					}
				}
				inputStateUpdated = ui.ObserveInputStateUpdated()
			case e, ok := <-screenSizeUpdated:
				if ok {
					type Handler interface {
						OnScreenSizeUpdated(e ebiten.ScreenSizeUpdatedEvent)
					}
					if game2, ok := game.(Handler); ok {
						game2.OnScreenSizeUpdated(e)
					}
				}
				screenSizeUpdated = ui.ObserveScreenSizeUpdated()
			case <-tick:
				game.Update()
			case canvas := <-drawing:
				game.Draw(canvas)
				drawing <- canvas
			}
		}
	}()

	for {
		ui.PollEvents()
		ui.Draw(func(c graphics.Canvas) {
			drawing <- graphics.NewLazyCanvas()
			canvas := <-drawing
			canvas.Flush(c)
		})
	}
}
