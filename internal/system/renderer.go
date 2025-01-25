package system

import (
	"bytes"
	"fmt"
	"image"
	_ "image/png"
	"log/slog"

	"github.com/damienfamed75/necs-interp-example/assets"
	"github.com/damienfamed75/necs-interp-example/shared"
	"github.com/hajimehoshi/ebiten/v2"
	"github.com/hajimehoshi/ebiten/v2/ebitenutil"
	"github.com/leap-fish/necs/esync"
	"github.com/yohamta/donburi"
	"github.com/yohamta/donburi/ecs"
	"github.com/yohamta/donburi/filter"
)

type RenderSystem struct {
	query   *donburi.Query
	fishImg *ebiten.Image
}

func NewRenderSystem() *RenderSystem {
	img, _, err := image.Decode(bytes.NewReader(assets.Fish_png))
	if err != nil {
		slog.Error("decode fish image", slog.Any("error", err))
		return nil
	}

	return &RenderSystem{
		fishImg: ebiten.NewImageFromImage(img),
		query:   donburi.NewQuery(filter.Contains(shared.PositionComponent)),
	}
}

func (r *RenderSystem) Draw(e *ecs.ECS, screen *ebiten.Image) {
	op := &ebiten.DrawImageOptions{}
	w, _ := ebiten.WindowSize()

	for entry := range r.query.Iter(e.World) {
		position := shared.PositionComponent.Get(entry)

		op.GeoM.Reset()
		op.GeoM.Scale(2, 2)
		op.GeoM.Translate(position.X, position.Y)

		screen.DrawImage(r.fishImg, op)
		ebitenutil.DebugPrintAt(
			screen,
			fmt.Sprintf("interpolated=%v", entry.HasComponent(esync.InterpComponent)),
			w/2, int(position.Y+10),
		)
	}
}
