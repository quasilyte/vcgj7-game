package worldsim

import (
	"github.com/quasilyte/gmath"
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
	for _, p := range r.world.Planets {
		p.MineralsDelay = gmath.ClampMin(p.MineralsDelay-delta, 0)
	}
}
