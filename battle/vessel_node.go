package battle

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/physics"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/vcgj7-game/assets"
	"github.com/quasilyte/vcgj7-game/gamedata"
)

type vesselNodeConfig struct {
	HP     float64
	Design *gamedata.VesselDesign
}

type vesselNode struct {
	state vesselState

	scene  *ge.Scene
	sprite *scalableSprite

	shield *scalableSprite

	config vesselNodeConfig

	body physics.Body

	pilotOrders vesselPilotOrders

	wrap posWrapper
}

type vesselPilotOrders struct {
	rotateLeft      bool
	rotateRight     bool
	forward         bool
	activateWeapon  bool
	activateSpecial bool
}

func newVesselNode(config vesselNodeConfig) *vesselNode {
	v := &vesselNode{
		config: config,
	}
	v.state.design = config.Design
	v.state.Pos = &v.body.Pos
	v.state.Rotation = &v.body.Rotation
	return v
}

func (v *vesselNode) Init(scene *ge.Scene) {
	state := &v.state

	v.body.InitCircle(v, 16)
	v.body.LayerMask = state.CollisionLayer
	scene.AddBody(&v.body)

	v.scene = scene

	state.Init()
	v.state.hp = v.state.design.MaxHP * v.config.HP
	v.state.Pos = &v.body.Pos
	v.state.Rotation = &v.body.Rotation

	v.sprite = newScalableSprite(v.state.design.Image, &v.body.Pos)
	scene.AddObject(v.sprite)
	v.sprite.s.Rotation = &v.body.Rotation

	v.shield = newScalableSprite(assets.ImageEnergyShield, &v.body.Pos)
	scene.AddObject(v.shield)
	v.shield.s.Rotation = &v.state.shieldRotation

	v.state.shieldRotation = v.body.Rotation
}

func (v *vesselNode) Dispose() {
	v.sprite.Dispose()
	v.shield.Dispose()
	v.body.Dispose()
}

func (v *vesselNode) Destroy() {
	// e := effects.NewExplosion(v.body.Pos)
	// v.scene.AddObject(e)

	v.Dispose()
}

func (v *vesselNode) IsDisposed() bool {
	return v.sprite.IsDisposed()
}

func (v *vesselNode) RotateLeftOrder() {
	v.pilotOrders.rotateLeft = true
}

func (v *vesselNode) RotateRightOrder() {
	v.pilotOrders.rotateRight = true
}

func (v *vesselNode) ForwardOrder() {
	v.pilotOrders.forward = true
}

func (v *vesselNode) ActivateSpecialOrder() {
	v.pilotOrders.activateSpecial = true
}

func (v *vesselNode) ActivateWeaponOrder() {
	v.pilotOrders.activateWeapon = true
}

func (v *vesselNode) OnDamage(weapon *gamedata.WeaponDesign, consumed bool) {
	// scorePos := ge.Pos{Offset: v.body.Pos.Add(gmath.Vec{Y: -48})}
	// score := effects.NewDamageScore(value, scorePos)
	// v.scene.AddObject(score)

	if v.state.hp <= 0 {
		return
	}

	damage := weapon.Damage
	if consumed {
		damage *= 0.25
	}
	v.state.hp = gmath.ClampMin(v.state.hp-damage, 0)
	if v.state.hp <= 0 {
		v.Destroy()
	}
}

func (v *vesselNode) Update(delta float64) {
	for _, collision := range v.scene.GetCollisions(&v.body) {
		switch obj := collision.Body.Object.(type) {
		case *projectileNode:
			consumed := false
			if obj.weapon.Blockable {
				projectileAngle := v.body.Pos.AngleToPoint(obj.body.Pos).Normalized()
				projectileAngleDelta := v.state.shieldRotation.Normalized().AngleDelta(projectileAngle)
				if projectileAngleDelta.Abs() < 1 {
					consumed = true
				}
			}
			obj.Destroy(!consumed)
			if consumed {
				energyGain := obj.weapon.EnergyCost * obj.weapon.EnergyConversion
				v.state.energy = gmath.ClampMax(v.state.energy+energyGain, v.state.design.MaxEnergy)
				playSound(v.scene, assets.AudioShieldAbsorb)
				v.OnDamage(obj.weapon, true)
			} else {
				v.OnDamage(obj.weapon, false)
			}
		}
	}

	v.state.Tick(delta)

	pilotOrders := v.pilotOrders
	v.pilotOrders = vesselPilotOrders{}

	if pilotOrders.forward {
		v.sprite.s.FrameOffset.X = v.sprite.s.FrameWidth
	} else {
		v.sprite.s.FrameOffset.X = 0
	}

	if pilotOrders.activateWeapon {
		if v.state.CanFire() {
			v.state.Fire()
			v.createProjectiles(v.state.design.MainWeapon)
		}
	}
	if pilotOrders.activateSpecial {
		if v.state.CanFireSecondary() {
			v.state.FireSecondary()
			v.createProjectiles(v.state.design.SecondaryWeapon)
		}
	}

	// if state.specialWeapon != nil {
	// 	state.specialWeapon.Update(delta)
	// 	if pilotOrders.activateSpecial && state.specialWeapon.CanActivate() {
	// 		state.specialWeapon.Activate(gmath.Vec{})
	// 	}
	// }

	// for i, w := range state.weapons {
	// 	if w == nil {
	// 		continue
	// 	}
	// 	w.Update(delta)
	// 	if autoOrders.activateWeapon && i == int(autoOrders.weaponIndex) && w.CanActivate() {
	// 		w.Activate(autoOrders.weaponTarget)
	// 	}
	// }

	v.applyMovement(delta, pilotOrders)

	v.wrap.Tick(delta, &v.body.Pos)
}

func (v *vesselNode) createProjectiles(weapon *gamedata.WeaponDesign) {
	playSound(v.scene, weapon.FireSound)

	targetPos := &v.state.enemy.body.Pos
	for i := 0; i < weapon.BurstSize; i++ {
		firePos := v.body.Pos
		offset := weapon.FireOffsets[i]
		if !offset.IsZero() {
			// Translate the offset.
			translatedOffset := offset.Rotated(v.body.Rotation)
			firePos = firePos.Add(translatedOffset)
		}
		projectileRotation := v.body.Rotation
		projectileRotation += weapon.ProjectileRotationDeltas[i]
		p := newProjectileNode(enemyCollisionMask(v.state.CollisionLayer), weapon, firePos, projectileRotation, targetPos)
		v.scene.AddObject(p)
	}

}

func (v *vesselNode) applyMovement(delta float64, orders vesselPilotOrders) {
	deceleration := 0.05
	state := &v.state

	rotationMultiplier := 1.0
	if orders.forward {
		rotationMultiplier = 0.7
	}

	// Adjust vessel rotation.
	var rotationDelta gmath.Rad
	if orders.rotateLeft {
		rotationDelta -= state.TotalRotationSpeed()
	}
	if orders.rotateRight {
		rotationDelta += state.TotalRotationSpeed()
	}
	if rotationDelta != 0 {
		r := gmath.Rad(float64(rotationDelta) * delta * rotationMultiplier)
		v.body.Rotation = (v.body.Rotation + r).Normalized()
		deceleration = 0.2
	}

	if orders.forward {
		accel := state.TotalAcceleration() * delta
		accelVector := gmath.RadToVec(v.body.Rotation).Mulf(accel)
		state.engineVelocity = state.engineVelocity.Add(accelVector)
		state.engineVelocity = state.engineVelocity.ClampLen(state.TotalMaxSpeed())
	} else if !state.engineVelocity.IsZero() {
		state.engineVelocity = state.engineVelocity.Mulf(1.0 - (deceleration * delta))
	}

	if !state.extraVelocity.IsZero() {
		state.extraVelocity = state.extraVelocity.Mulf(1.0 - (0.3 * delta))
	}

	v.body.Pos = v.body.Pos.Add(state.TotalVelocity().Mulf(delta))
}
