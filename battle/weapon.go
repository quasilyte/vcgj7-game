package battle

import (
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/vcgj7-game/gamedata"
)

type weapon struct {
	reload float64
	design *gamedata.WeaponDesign
}

func (w *weapon) Tick(delta float64) {
	w.reload = gmath.ClampMin(w.reload-delta, 0)
}
