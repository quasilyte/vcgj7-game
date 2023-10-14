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

	hp                   float64
	energy               float64
	energyRegenThreshold float64

	weapon *weapon

	design *gamedata.VesselDesign
}

func (state *vesselState) Init() {
	state.hp = state.design.MaxHP
	state.energy = state.design.MaxEnergy

	state.energyRegenThreshold = state.design.MaxEnergy * 0.5

	if state.design.MainWeapon != nil {
		state.weapon = &weapon{
			design: state.design.MainWeapon,
		}
	}
}

func (state *vesselState) Tick(delta float64) {
	if state.weapon != nil {
		state.weapon.Tick(delta)
	}

	if state.energy < state.energyRegenThreshold {
		state.energy = gmath.ClampMax(state.energy+state.design.EnergyRegen*delta, state.energyRegenThreshold)
	}
}

func (state *vesselState) CanFire() bool {
	if state.weapon == nil {
		return false
	}
	if state.weapon.reload > 0 {
		return false
	}
	if state.energy < state.weapon.design.EnergyCost {
		return false
	}
	return true
}

func (state *vesselState) Fire() {
	if state.weapon.design.EnergyCost != 0 {
		state.energy -= state.weapon.design.EnergyCost
	}
	state.weapon.reload = state.weapon.design.Reload
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
