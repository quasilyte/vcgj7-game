package battle

import (
	"math"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/input"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/vcgj7-game/assets"
	"github.com/quasilyte/vcgj7-game/gamedata"
)

type Runner struct {
	scene *ge.Scene

	pilots []pilot

	input *input.Handler
}

func NewRunner(h *input.Handler) *Runner {
	return &Runner{
		input: h,
	}
}

func (r *Runner) IsDisposed() bool {
	return false
}

func (r *Runner) Init(scene *ge.Scene) {
	r.scene = scene

	v := newVesselNode()
	v.body.Pos = (gmath.Vec{X: 1920 / 4, Y: 1080 / 4}).Sub(gmath.Vec{X: 200})
	v.body.LayerMask = collisionPlayer1
	v.state.CollisionLayer = v.body.LayerMask
	v.state.design = &gamedata.VesselDesign{
		Image:         assets.ImageVesselRaider,
		MaxHP:         100,
		MaxEnergy:     100,
		EnergyRegen:   1,
		MaxSpeed:      250,
		Acceleration:  100,
		RotationSpeed: 4,
		MainWeapon:    gamedata.FindWeaponDesign("Ion Cannon"),
	}

	scene.AddObject(v)

	p := newHumanPilot(r.input, v)
	r.pilots = append(r.pilots, p)

	{
		v2 := newVesselNode()
		v2.body.LayerMask = collisionPlayer2
		v2.body.Pos = (gmath.Vec{X: 1920 / 4, Y: 1080 / 4}).Add(gmath.Vec{X: 200})
		v2.body.Rotation = -math.Pi
		v2.state.CollisionLayer = v2.body.LayerMask
		v2.state.design = &gamedata.VesselDesign{
			Image:         assets.ImageVesselMarauder,
			MaxHP:         150,
			MaxEnergy:     120,
			EnergyRegen:   1.5,
			MaxSpeed:      190,
			Acceleration:  120,
			RotationSpeed: 2.5,
		}
		scene.AddObject(v2)
	}

	hud := scene.NewSprite(assets.ImageBattleHUD)
	hud.Centered = false
	scene.AddGraphicsAbove(hud, 1)

	{
		v.state.hp = v.state.design.MaxHP * 0.6
		pos := gmath.Vec{X: 178, Y: 50}
		hpBar := newValueBar(pos, &v.state.hp, v.state.design.MaxHP, true)
		scene.AddObject(hpBar)
	}
	{
		v.state.energy = v.state.design.MaxEnergy * 0.9
		pos := gmath.Vec{X: 178 + 494, Y: 50}
		hpBar := newValueBar(pos, &v.state.energy, v.state.design.MaxEnergy, false)
		scene.AddObject(hpBar)
	}
}

func (r *Runner) Update(delta float64) {
	for _, p := range r.pilots {
		p.Update(delta)
	}
}
