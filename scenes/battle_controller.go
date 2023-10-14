package scenes

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/vcgj7-game/battle"
	"github.com/quasilyte/vcgj7-game/session"
)

type BattleController struct {
	state *session.State

	runner *battle.Runner
}

func NewBattleController(state *session.State) *BattleController {
	return &BattleController{state: state}
}

func (c *BattleController) Init(scene *ge.Scene) {
	c.runner = battle.NewRunner(c.state.Input)
	scene.AddObject(c.runner)
}

func (c *BattleController) Update(delta float64) {}
