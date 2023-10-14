package battle

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/vcgj7-game/assets"
)

type valueBar struct {
	sprite *ge.Sprite

	pos gmath.Vec

	value    *float64
	maxValue float64

	isHP bool
}

func newValueBar(pos gmath.Vec, value *float64, maxValue float64, isHP bool) *valueBar {
	return &valueBar{
		pos:      pos,
		isHP:     isHP,
		value:    value,
		maxValue: maxValue,
	}
}

func (b *valueBar) Init(scene *ge.Scene) {
	img := assets.ImageBattleBarEnergy
	if b.isHP {
		img = assets.ImageBattleBarHP
	}
	s := scene.NewSprite(img)
	s.Centered = false
	s.FlipVertical = true
	s.Pos.Base = &b.pos
	b.sprite = s
	scene.AddGraphicsAbove(s, 1)

	b.updateFrame()
}

func (b *valueBar) IsDisposed() bool { return false }

func (b *valueBar) updateFrame() {
	percent := *b.value / b.maxValue
	if percent >= 0.99 {
		b.sprite.FrameTrimBottom = 0
		return
	}
	b.sprite.FrameTrimBottom = (1.0 - percent) * b.sprite.FrameHeight
}

func (b *valueBar) Update(delta float64) {
	b.updateFrame()
}
