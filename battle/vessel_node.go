package battle

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/physics"
	"github.com/quasilyte/gmath"
)

type vesselNode struct {
	state vesselState

	scene  *ge.Scene
	sprite *ge.Sprite

	body physics.Body

	pilotOrders vesselPilotOrders

	wrap posWrapper
}

type vesselPilotOrders struct {
	rotateLeft      bool
	rotateRight     bool
	forward         bool
	activateSpecial bool
}

func newVesselNode() *vesselNode {
	v := &vesselNode{}
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

	v.sprite = scene.NewSprite(v.state.design.Image)
	v.sprite.Pos.Base = &v.body.Pos
	v.sprite.Rotation = &v.body.Rotation
	scene.AddGraphics(v.sprite)
}

func (v *vesselNode) Dispose() {
	v.sprite.Dispose()
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

func (v *vesselNode) OnDamage(value float64) {
	// scorePos := ge.Pos{Offset: v.body.Pos.Add(gmath.Vec{Y: -48})}
	// score := effects.NewDamageScore(value, scorePos)
	// v.scene.AddObject(score)
	// v.state.hp = gmath.ClampMin(v.state.hp-value, 0)

	// if v.state.hp == 0 {
	// 	v.Destroy()
	// }
}

func (v *vesselNode) Update(delta float64) {
	// for _, collision := range v.scene.GetCollisions(&v.body) {
	// 	switch obj := collision.Body.Object.(type) {
	// 	case *simpleProjectile:
	// 		obj.Destroy()
	// 		v.OnDamage(obj.design.damage)
	// 	case *rocketProjectile:
	// 		obj.Destroy()
	// 		v.OnDamage(obj.design.damage)
	// 	case *AfterburnerFlame:
	// 		obj.Destroy()
	// 		v.OnDamage(15)
	// 	}
	// }

	if v.state.energy < v.state.energyRegenThreshold {
		v.state.energy = gmath.ClampMax(v.state.energy+v.state.design.EnergyRegen*delta, v.state.energyRegenThreshold)
	}

	pilotOrders := v.pilotOrders
	// autoOrders := v.systems.Tick(delta)
	// state := &v.state
	v.pilotOrders = vesselPilotOrders{}

	if pilotOrders.forward {
		v.sprite.FrameOffset.X = v.sprite.FrameWidth
	} else {
		v.sprite.FrameOffset.X = 0
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
