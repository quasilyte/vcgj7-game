package battle

import (
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/vcgj7-game/gamedata"
)

type vesselState struct {
	enemy          *vesselState
	CollisionLayer uint16

	Pos      *gmath.Vec
	Rotation *gmath.Rad

	engineVelocity gmath.Vec
	extraVelocity  gmath.Vec

	hp     float64
	energy float64

	design *gamedata.VesselDesign
}

func (state *vesselState) EnergyPercentage() float64 {
	return state.energy / state.design.MaxEnergy
}

func (state *vesselState) TotalRotationSpeed() gmath.Rad {
	return state.design.RotationSpeed
}

func (state *vesselState) TotalMaxSpeed() float64 {
	return state.design.MaxSpeed
}

func (state *vesselState) TotalAcceleration() float64 {
	return state.design.Acceleration
}

func (state *vesselState) TotalVelocity() gmath.Vec {
	velocity := state.engineVelocity.Add(state.extraVelocity)
	return velocity
}
