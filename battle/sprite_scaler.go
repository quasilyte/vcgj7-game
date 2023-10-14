package battle

import (
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
)

var screenCenter = gmath.Vec{
	X: 1920 / 4,
	Y: 1080 / 4,
}

type scalableSprite struct {
	s   *ge.Sprite
	img resource.ImageID
	pos *gmath.Vec
}

func calcSpriteScaling(pos gmath.Vec) float64 {
	centerDist := pos.DistanceTo(screenCenter)
	spriteScale := 1.0 - gmath.ClampMin(centerDist-140, 0)*0.0015
	return spriteScale
}

func newScalableSprite(img resource.ImageID, pos *gmath.Vec) *scalableSprite {
	return &scalableSprite{
		img: img,
		pos: pos,
	}
}

func (s *scalableSprite) Init(scene *ge.Scene) {
	s.s = scene.NewSprite(s.img)
	s.s.Pos.Base = s.pos
	scene.AddGraphics(s.s)
	s.resize()
}

func (s *scalableSprite) Update(delta float64) {
	s.resize()
}

func (s *scalableSprite) resize() {
	spriteScale := calcSpriteScaling(s.s.Pos.Resolve())
	s.s.SetScale(spriteScale, spriteScale)
}

func (s *scalableSprite) IsDisposed() bool {
	return s.s.IsDisposed()
}

func (s *scalableSprite) Dispose() {
	s.s.Dispose()
}
