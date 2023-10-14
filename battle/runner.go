package battle

import (
	"math"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/input"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/gsignal"
	"github.com/quasilyte/vcgj7-game/assets"
	"github.com/quasilyte/vcgj7-game/gamedata"
)

type Runner struct {
	scene *ge.Scene

	pilots []pilot

	player *gamedata.Player

	input *input.Handler

	playerVessel *vesselNode
	enemyVessel  *vesselNode

	enemyDesign *gamedata.VesselDesign

	EventBattleOver gsignal.Event[Results]
}

type Results struct {
	Victory bool
	HP      float64
}

type RunnerConfig struct {
	Input  *input.Handler
	Player *gamedata.Player
	Enemy  *gamedata.VesselDesign
}

func NewRunner(config RunnerConfig) *Runner {
	return &Runner{
		input:       config.Input,
		enemyDesign: config.Enemy,
		player:      config.Player,
	}
}

func (r *Runner) IsDisposed() bool {
	return false
}

func (r *Runner) Init(scene *ge.Scene) {
	r.scene = scene

	bg := scene.NewSprite(assets.ImageBattleBg)
	bg.Pos.Offset.X = 210
	bg.Centered = false
	scene.AddGraphicsBelow(bg, 1)

	v := newVesselNode(vesselNodeConfig{
		HP:     r.player.VesselHP,
		Design: r.player.VesselDesign,
	})
	v.body.Pos = (gmath.Vec{X: 1920 / 4, Y: 1080 / 4}).Sub(gmath.Vec{X: 240})
	v.body.LayerMask = collisionPlayer1
	v.body.Rotation = 0.2
	v.state.CollisionLayer = v.body.LayerMask
	v.EventDestroyed.Connect(nil, r.onDefeat)
	scene.AddObject(v)
	r.playerVessel = v

	p := newHumanPilot(r.input, v)
	r.pilots = append(r.pilots, p)

	{
		v2 := newVesselNode(vesselNodeConfig{
			HP:     1,
			Design: r.enemyDesign,
		})
		v2.body.LayerMask = collisionPlayer2
		v2.body.Pos = (gmath.Vec{X: 1920 / 4, Y: 1080 / 4}).Add(gmath.Vec{X: 240})
		v2.body.Rotation = -math.Pi + 0.2
		v2.state.CollisionLayer = v2.body.LayerMask
		v2.EventDestroyed.Connect(nil, r.onVictory)
		scene.AddObject(v2)
		r.enemyVessel = v2

		v2.state.enemy = v
		v.state.enemy = v2

		p := newComputerPilot(v2, botDummy, scene)
		r.pilots = append(r.pilots, p)
	}

	hud := scene.NewSprite(assets.ImageBattleHUD)
	hud.Centered = false
	scene.AddGraphicsAbove(hud, 1)

	{
		pos := gmath.Vec{X: 178, Y: 50}
		hpBar := newValueBar(pos, &v.state.hp, v.state.design.MaxHP, true)
		scene.AddObject(hpBar)
	}
	{
		pos := gmath.Vec{X: 178 + 494, Y: 50}
		hpBar := newValueBar(pos, &v.state.energy, v.state.design.MaxEnergy, false)
		scene.AddObject(hpBar)
	}
}

func (r *Runner) onDefeat(gsignal.Void) {
	r.enemyVessel.body.LayerMask = 0
	r.EventBattleOver.Emit(Results{
		Victory: false,
	})
}

func (r *Runner) onVictory(gsignal.Void) {
	r.playerVessel.body.LayerMask = 0
	r.EventBattleOver.Emit(Results{
		Victory: true,
		HP:      r.playerVessel.state.HealthPercentage(),
	})
}

func (r *Runner) Update(delta float64) {
	for _, p := range r.pilots {
		p.Update(delta)
	}
}
