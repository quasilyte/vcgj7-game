package battle

import (
	"github.com/quasilyte/ge/input"
	"github.com/quasilyte/vcgj7-game/controls"
)

type pilot interface {
	Update(delta float64)
}

type humanPilot struct {
	input  *input.Handler
	vessel *vesselNode
}

func newHumanPilot(h *input.Handler, v *vesselNode) *humanPilot {
	return &humanPilot{
		input:  h,
		vessel: v,
	}
}

func (p *humanPilot) Update(delta float64) {
	if p.input.ActionIsPressed(controls.ActionForward) {
		p.vessel.ForwardOrder()
	}
	if p.input.ActionIsPressed(controls.ActionLeft) {
		p.vessel.RotateLeftOrder()
	}
	if p.input.ActionIsPressed(controls.ActionRight) {
		p.vessel.RotateRightOrder()
	}
}
