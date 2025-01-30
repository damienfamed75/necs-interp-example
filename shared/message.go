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

func lerpVector2(from, to dmath.Vec2, t float64) *dmath.Vec2 {
	return &dmath.Vec2{
		X: lerp(from.X, to.X, t),
		Y: lerp(from.Y, to.Y, t),
	}
}

func RegisterComponents() {
	esync.RegisterComponent(1, esync.NetworkId(0), esync.NetworkIdComponent)
	esync.RegisterComponent(2, esync.InterpData{}, esync.InterpComponent)
	// Specify a lerping function to use when interpolating client-side
	esync.RegisterComponent(10, dmath.Vec2{}, PositionComponent, esync.WithInterpFn(1, lerpVector2))
}
