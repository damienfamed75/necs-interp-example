package main

import (
	"log"

	"github.com/coder/websocket"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leap-fish/necs/esync/clisync"
	"github.com/leap-fish/necs/router"
	"github.com/leap-fish/necs/transports"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"

	"github.com/damienfamed75/necs-interp-example/internal/game"
	"github.com/damienfamed75/necs-interp-example/shared"
)

const (
	screenWidth  = 640
	screenHeight = 480
)

func main() {
	client := transports.NewWsClientTransport("ws://localhost:7373")

	router.OnError(func(sender *router.NetworkClient, err error) {
		log.Printf("Message Error: %s", err.Error())
	})

	e := ecs.NewECS(donburi.NewWorld())
	clisync.RegisterClient(e.World)

	shared.RegisterComponents()

	go func() {
		err := client.Start(func(conn *websocket.Conn) {
			// If you want to use the connection for other purposes, this is where you might want to
			// store it for later use.
		})
		if err != nil {
			log.Fatalf("Unable to dial: %s", err)
		}
	}()

	ebiten.SetWindowSize(screenWidth, screenHeight)
	ebiten.SetWindowTitle("Interpolation (NECS Example)")

	g := game.New(e)
	if err := ebiten.RunGame(g); err != nil {
		panic(err)
	}
}
