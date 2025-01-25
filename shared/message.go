package shared

import (
	"github.com/leap-fish/necs/esync"
	"github.com/yohamta/donburi"
	dmath "github.com/yohamta/donburi/features/math"
)

type TestMessage struct {
	Message string
}

var PositionComponent = donburi.NewComponentType[dmath.Vec2]()

func lerp(a, b, t float64) float64 {
	return (1.0-t)*a + b*t
}

func RegisterComponents() {
	esync.RegisterComponent(1, esync.NetworkId(0), esync.NetworkIdComponent)
	esync.RegisterComponent(2, esync.InterpData{}, esync.InterpComponent)
	esync.RegisterComponent(10, dmath.Vec2{}, PositionComponent)

	// Interpolation components use a separate ids
	esync.RegisterInterpolated(1, PositionComponent, func(from, to dmath.Vec2, t float64) *dmath.Vec2 {
		return &dmath.Vec2{
			X: lerp(from.X, to.X, t),
			Y: lerp(from.Y, to.Y, t),
		}
	})
}
