package battle

import "github.com/quasilyte/ge"

type botKind int

const (
	botDummy botKind = iota
)

type computerPilot struct {
	impl pilot
}

func newComputerPilot(v *vesselNode, kind botKind, scene *ge.Scene) *computerPilot {
	p := &computerPilot{}
	switch kind {
	case botDummy:
		p.impl = newDummyComputerPilot(v, scene)
	}
	return p
}

func (p *computerPilot) Update(delta float64) {
	p.impl.Update(delta)
}
