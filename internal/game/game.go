package game

import (
	"log/slog"

	"github.com/damienfamed75/necs-interp-example/internal/system"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/leap-fish/necs/esync/clisync"
	"github.com/yohamta/donburi/ecs"
)

var _ ebiten.Game = &Game{}

type Game struct {
	w, h   int
	ecs    *ecs.ECS
	logger *slog.Logger
}

func New(e *ecs.ECS) *Game {
	e.AddRenderer(0, system.NewRenderSystem().Draw).
		AddSystem(clisync.NewInterpolateSystem())

	return &Game{
		ecs:    e,
		logger: slog.Default().With("layer", "game"),
	}
}

func (g *Game) Layout(ow, oh int) (int, int) {
	g.w, g.h = ow, oh
	return ow, oh
}

func (g *Game) Update() error {
	g.ecs.Update()

	return nil
}

func (g *Game) Draw(screen *ebiten.Image) {
	g.ecs.Draw(screen)
}
