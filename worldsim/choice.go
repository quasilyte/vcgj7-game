package worldsim

import (
	"github.com/quasilyte/vcgj7-game/gamedata"
)

type Choice struct {
	Time       int // In hours
	Text       string
	Mode       gamedata.Mode // In-process mode
	OnResolved func() gamedata.Mode
}

const MaxChoices = 6
