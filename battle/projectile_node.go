package battle

import (
	"math"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/physics"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/vcgj7-game/gamedata"
)

type projectileNode struct {
	collisionLayer uint16

	weapon *gamedata.WeaponDesign

	wrap posWrapper

	hp       float64
	velocity gmath.Vec

	body physics.Body

	sprite *scalableSprite
	scene  *ge.Scene
}

func newProjectileNode(collisionLayer uint16, weapon *gamedata.WeaponDesign, pos gmath.Vec, rotation gmath.Rad) *projectileNode {
	p := &projectileNode{
		collisionLayer: collisionLayer,
		weapon:         weapon,
	}
	p.body.Pos = pos
	p.body.Rotation = rotation
	return p
}

func (p *projectileNode) Init(scene *ge.Scene) {
	p.scene = scene
	p.hp = scene.Rand().FloatRange(0.9, 1.1) * p.weapon.Range

	p.body.InitCircle(p, math.Round(p.weapon.ProjectileSize*0.5))
	p.body.LayerMask = p.collisionLayer
	scene.AddBody(&p.body)

	p.sprite = newScalableSprite(p.weapon.ProjectileImage, &p.body.Pos)
	scene.AddObjectBelow(p.sprite, 1)
	p.sprite.s.Rotation = &p.body.Rotation

	p.velocity = gmath.RadToVec(p.body.Rotation).Mulf(p.weapon.ProjectileSpeed)
}

func (p *projectileNode) IsDisposed() bool {
	return p.sprite.IsDisposed()
}

func (p *projectileNode) Update(delta float64) {
	p.hp -= delta * p.weapon.ProjectileSpeed
	if p.hp <= 0 {
		p.Destroy(false)
		return
	}

	// accel := p.seek()
	// p.velocity = p.velocity.Add(accel.Mulf(delta)).ClampLen(p.design.projectileSpeed)
	// p.body.Rotation = p.velocity.Angle()
	// p.body.Pos = p.body.Pos.Add(p.velocity.Mulf(delta))
	p.body.Pos = p.body.Pos.Add(p.velocity.Mulf(delta))

	p.wrap.Tick(delta, &p.body.Pos)
}

// func (p *rocketProjectile) seek() gmath.Vec {
// 	dst := p.target.Sub(p.body.Pos).Normalized().Mulf(p.design.projectileSpeed)
// 	return dst.Sub(p.velocity).Normalized().Mulf(p.design.homingSteer)
// }

func (p *projectileNode) Dispose() {
	p.sprite.Dispose()
	p.body.Dispose()
}

func (p *projectileNode) Destroy(impact bool) {
	if p.weapon.Explosion != 0 && impact {
		e := newEffectNode(p.body.Pos, normalEffectLayer, p.weapon.Explosion)
		p.scene.AddObject(e)
		e.anim.SetSecondsPerFrame(0.035)
		if p.weapon.ExplosionSound != 0 {
			playSound(p.scene, p.weapon.ExplosionSound)
		}
	}

	p.Dispose()
}
