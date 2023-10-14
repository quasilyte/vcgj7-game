package battle

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
)

type dummyComputerPilot struct {
	vessel *vesselNode
	enemy  *vesselNode
	scene  *ge.Scene

	screenCenter gmath.Vec
	centerOffset gmath.Vec

	noTurnTime     float64
	randTurnTime   float64
	randTurnLeft   bool
	alignTurnTime  float64
	centerTurnTime float64

	thrustTime  float64
	thrustDelay float64

	noAttackDelay float64
	agressiveTime float64

	angleToTarget    gmath.Rad
	targetAngleDelta gmath.Rad
}

func newDummyComputerPilot(v *vesselNode, scene *ge.Scene) *dummyComputerPilot {
	return &dummyComputerPilot{
		vessel: v,
		enemy:  v.state.enemy,
		scene:  scene,
		screenCenter: gmath.Vec{
			X: 1920 / 4,
			Y: 1080 / 4,
		},
	}
}

func (p *dummyComputerPilot) Update(delta float64) {
	p.angleToTarget = p.vessel.body.Pos.AngleToPoint(p.enemy.body.Pos).Normalized()
	p.targetAngleDelta = p.vessel.body.Rotation.Normalized().AngleDelta(p.angleToTarget)

	p.navigate(delta)
	p.attack(delta)
}

func (p *dummyComputerPilot) attack(delta float64) {
	enemyDist := p.vessel.body.Pos.DistanceTo(p.enemy.body.Pos)

	if p.agressiveTime == 0 {
		noAttackDecay := delta
		if enemyDist < 175 {
			noAttackDecay *= 2
		}
		if enemyDist < 80 {
			noAttackDecay *= 3
		}
		p.noAttackDelay = gmath.ClampMin(p.noAttackDelay-noAttackDecay, 0)
		if p.noAttackDelay > 0 {
			return
		}

		noAttackChance := (p.vessel.state.hp / p.vessel.state.design.MaxHP) * 0.3
		if enemyDist < 150 {
			noAttackChance *= 0.3
		}
		if p.scene.Rand().Chance(noAttackChance) {
			p.noAttackDelay = p.scene.Rand().FloatRange(0.8, 2)
			return
		}
	}

	if p.agressiveTime == 0 {
		aggressiveChance := 0.05
		if enemyDist < 150 {
			aggressiveChance += 0.04
		}
		if p.vessel.state.HealthPercentage() < 0.5 {
			aggressiveChance *= 2
		}
		if p.scene.Rand().Chance(aggressiveChance) {
			p.agressiveTime = p.scene.Rand().FloatRange(0.5, 2)
			if p.scene.Rand().Bool() {
				p.agressiveTime *= 2
			}
		}
	}

	maxAngleDelta := 0.15
	switch {
	case enemyDist > 100:
		maxAngleDelta = 0.2
	case enemyDist > 200:
		maxAngleDelta = 0.3
	case enemyDist > 300:
		maxAngleDelta = 0.4
	}
	maxAngleDelta *= p.scene.Rand().FloatRange(0.8, 1.4)
	if p.targetAngleDelta.Abs() <= maxAngleDelta {
		// TODO: don't fire if facing the shield.
		if p.vessel.state.CanFireSecondary() && p.scene.Rand().Chance(0.4) {
			p.vessel.ActivateSpecialOrder()
		} else {
			p.vessel.ActivateWeaponOrder()
		}
	}
}

func (p *dummyComputerPilot) navigate(delta float64) {
	if p.centerTurnTime > 0 {
		angleToCenter := p.vessel.body.Pos.AngleToPoint(p.centerOffset).Normalized()
		angleDelta := p.vessel.body.Rotation.Normalized().AngleDelta(angleToCenter)
		if angleDelta.Abs() < 0.3 {
			p.vessel.ForwardOrder()
		}
	} else {
		if p.thrustTime > 0 {
			p.thrustTime = gmath.ClampMin(p.thrustTime-delta, 0)
			p.vessel.ForwardOrder()
		} else {
			p.thrustDelay = gmath.ClampMin(p.thrustDelay-delta, 0)
			if p.thrustDelay == 0 {
				if p.scene.Rand().Chance(0.4) {
					p.thrustDelay = p.scene.Rand().FloatRange(0.4, 3.8)
					if p.scene.Rand().Chance(0.4) {
						p.thrustDelay *= 2
					}
				} else {
					p.thrustTime = p.scene.Rand().FloatRange(0.4, 3.8)
				}
			}
		}
	}

	if p.noTurnTime > 0 {
		p.noTurnTime = gmath.ClampMin(p.noTurnTime-delta, 0)
		return
	}

	if p.randTurnTime > 0 {
		p.randTurnTime = gmath.ClampMin(p.randTurnTime-delta, 0)
		if p.randTurnLeft {
			p.vessel.RotateLeftOrder()
		} else {
			p.vessel.RotateRightOrder()
		}
		return
	}

	if p.alignTurnTime > 0 {
		p.alignTurnTime = gmath.ClampMin(p.alignTurnTime-delta, 0)
		angleDelta := p.targetAngleDelta
		if angleDelta.Abs() < 0.2 {
			return
		}
		if angleDelta < 0 {
			p.vessel.RotateLeftOrder()
		} else {
			p.vessel.RotateRightOrder()
		}
		return
	}

	if p.centerTurnTime > 0 {
		p.centerTurnTime = gmath.ClampMin(p.centerTurnTime-delta, 0)
		angleToCenter := p.vessel.body.Pos.AngleToPoint(p.centerOffset).Normalized()
		angleDelta := p.vessel.body.Rotation.Normalized().AngleDelta(angleToCenter)
		if angleDelta.Abs() < 0.2 {
			return
		}
		if angleDelta < 0 {
			p.vessel.RotateLeftOrder()
		} else {
			p.vessel.RotateRightOrder()
		}
		return
	}

	switch roll := p.scene.Rand().Float(); {
	case roll < 0.45:
		if p.vessel.body.Pos.DistanceTo(p.screenCenter) > 200 && p.scene.Rand().Chance(0.4) {
			p.centerTurnTime = p.scene.Rand().FloatRange(0.7, 1.6)
			p.centerOffset = p.screenCenter.Add(p.scene.Rand().Offset(-64, 64))
		} else {
			p.alignTurnTime = p.scene.Rand().FloatRange(0.5, 1)
		}
	case roll < 0.75:
		if p.vessel.state.HealthPercentage() < 0.5 && p.scene.Rand().Chance(0.5) {
			p.alignTurnTime = p.scene.Rand().FloatRange(0.6, 1.1)
		} else {
			p.noTurnTime = p.scene.Rand().FloatRange(0.3, 1.7)
		}
	default:
		p.randTurnTime = p.scene.Rand().FloatRange(0.2, 0.5)
		p.randTurnLeft = p.scene.Rand().Bool()
	}
}
