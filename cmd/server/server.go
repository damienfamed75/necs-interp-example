package main

import (
	"log"
	"log/slog"
	"math"
	"time"

	"github.com/leap-fish/necs/esync/srvsync"
	"github.com/leap-fish/necs/router"
	"github.com/leap-fish/necs/transports"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"

	"github.com/damienfamed75/necs-interp-example/shared"
)

const (
	_tickRate = 15
)

func main() {
	router.OnError(func(sender *router.NetworkClient, err error) {
		log.Printf("Message Error: %s", err.Error())
	})

	e := ecs.NewECS(donburi.NewWorld())
	srvsync.UseEsync(e.World)

	// Register a contract that both the client and server understand
	shared.RegisterComponents()
	e.AddSystem(newMoveSystem())

	nonInterpEntity := e.World.Create(shared.PositionComponent)
	shared.PositionComponent.Get(e.World.Entry(nonInterpEntity)).Y = 40
	srvsync.NetworkSync(e.World, &nonInterpEntity, shared.PositionComponent)

	// Create an entity normally and specify that we want to tell the client
	// to interpolate its Position component.
	interpEntity := e.World.Create(shared.PositionComponent)
	shared.PositionComponent.Get(e.World.Entry(interpEntity)).Y = 220
	srvsync.NetworkSync(e.World, &interpEntity, shared.PositionComponent)
	srvsync.NetworkInterp(e.World, &interpEntity, shared.PositionComponent)

	go func() {
		for range time.NewTicker(time.Second / _tickRate).C {
			e.Update()

			if err := srvsync.DoSync(); err != nil {
				slog.Error("syncronize", slog.Any("error", err))
			}
		}
	}()

	server := transports.NewWsServerTransport(7373, "", nil)
	err := server.Start()
	if err != nil {
		log.Fatalf("Unable to dial: %s", err)
	}
}

func newMoveSystem() func(*ecs.ECS) {
	query := donburi.NewQuery(filter.Contains(shared.PositionComponent))

	time := 0.0
	return func(e *ecs.ECS) {
		time++
		x := math.Sin(time*0.25)*80 + 200

		for entry := range query.Iter(e.World) {
			pos := shared.PositionComponent.Get(entry)
			pos.X = x
		}
	}
}
