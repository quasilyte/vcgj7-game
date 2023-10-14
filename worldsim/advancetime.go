package worldsim

import (
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/vcgj7-game/gamedata"
)

func (r *Runner) AdvanceTime(hours int) {
	for i := 0; i < hours; i++ {
		r.world.GameTime++
		// One in-game hour is simulated during 1 second in delta time terms.
		for j := 0; j < 5; j++ {
			r.updateWorld(0.2)
		}
	}
}

func (r *Runner) updateWorld(delta float64) {
	r.world.UpgradeRerollDelay = gmath.ClampMin(r.world.UpgradeRerollDelay-delta, 0)
	r.world.NextUpgradeDelay = gmath.ClampMin(r.world.NextUpgradeDelay-delta, 0)
	if r.world.UpgradeRerollDelay == 0 {
		r.world.UpgradeRerollDelay = float64(r.scene.Rand().IntRange(5, 15))
		r.world.UpgradeAvailable = gamedata.UpgradeKind(r.scene.Rand().IntRange(int(gamedata.FirstUpgrade), int(gamedata.LastUpgrade)))
	}

	for _, p := range r.world.Planets {
		p.MineralsDelay = gmath.ClampMin(p.MineralsDelay-delta, 0)
	}
}
