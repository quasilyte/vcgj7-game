package battle

import (
	"github.com/quasilyte/gmath"
)

type posWrapper struct {
	wrapDelay float64
}

func (w *posWrapper) Tick(delta float64, pos *gmath.Vec) {
	w.wrapDelay = gmath.ClampMin(w.wrapDelay-delta, 0)

	center := gmath.Vec{X: (1920 / 4), Y: (1080 / 4)}
	if w.wrapDelay == 0 && pos.DistanceTo(center) > ((1080/4)+20) {
		*pos = pos.Sub(center).Mulf(-0.98).Add(center)
		w.wrapDelay = 0.2
	}
}
