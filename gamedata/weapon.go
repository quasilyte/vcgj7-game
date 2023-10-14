package gamedata

import (
	"fmt"

	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/vcgj7-game/assets"
)

type WeaponDesign struct {
	Name string

	FireSound resource.AudioID

	Damage float64
	Reload float64

	Homing float64

	Range           float64
	ProjectileSpeed float64
	ProjectileSize  float64
	ProjectileImage resource.ImageID

	BurstSize                int
	FireOffsets              []gmath.Vec
	ProjectileRotationDeltas []gmath.Rad

	Explosion      resource.ImageID
	ExplosionSound resource.AudioID

	EnergyCost       float64
	EnergyConversion float64
	Blockable        bool
}

func FindWeaponDesign(name string) *WeaponDesign {
	for _, w := range Weapons {
		if w.Name == name {
			return w
		}
	}
	panic(fmt.Sprintf("weapon %q not found", name))
}

var Weapons = []*WeaponDesign{
	{
		Name:             "Pulse Laser",
		FireSound:        assets.AudioPulseLaser1,
		Damage:           8,
		Reload:           0.25,
		EnergyCost:       5,
		EnergyConversion: 2.0,
		Range:            350,
		ProjectileSpeed:  280,
		ProjectileImage:  assets.ImageProjectilePulseLaser,
		ProjectileSize:   6,
		Blockable:        true,
		BurstSize:        1,

		FireOffsets:              []gmath.Vec{{}},
		ProjectileRotationDeltas: []gmath.Rad{0},
	},

	{
		Name:             "Ion Cannon",
		FireSound:        assets.AudioIonCannon1,
		Damage:           10,
		Reload:           0.4,
		EnergyCost:       4,
		EnergyConversion: 0.5,
		Range:            450,
		ProjectileSpeed:  320,
		ProjectileImage:  assets.ImageProjectileIonCannon,
		ProjectileSize:   8,
		Explosion:        assets.ImageIonCannonImpact,
		ExplosionSound:   assets.AudioIonCannonImpact,
		Blockable:        true,
		BurstSize:        1,

		FireOffsets:              []gmath.Vec{{}},
		ProjectileRotationDeltas: []gmath.Rad{0},
	},

	{
		Name:             "Assault Laser",
		FireSound:        assets.AudioAssaultLaser1,
		Damage:           8,
		Reload:           0.2,
		EnergyCost:       4,
		EnergyConversion: 2.5,
		Range:            260,
		ProjectileSpeed:  400,
		ProjectileImage:  assets.ImageProjectileAssaultLaser,
		ProjectileSize:   6,
		Explosion:        assets.ImageAssaultLaserImpact,
		// ExplosionSound:   assets.AudioIonCannonImpact,
		Blockable: true,
		BurstSize: 2,

		FireOffsets: []gmath.Vec{
			{Y: -6},
			{Y: +6},
		},
		ProjectileRotationDeltas: []gmath.Rad{0, 0},
	},

	{
		Name:            "Missile Launcher",
		FireSound:       assets.AudioMissile1,
		Damage:          20,
		Reload:          3.5,
		Range:           700,
		ProjectileSpeed: 250,
		ProjectileImage: assets.ImageProjectileMissile,
		ProjectileSize:  10,
		Explosion:       assets.ImageMissileImpact,
		ExplosionSound:  assets.AudioExplosion1,
		BurstSize:       3,

		FireOffsets: []gmath.Vec{
			{Y: -8},
			{},
			{Y: +8},
		},
		ProjectileRotationDeltas: []gmath.Rad{
			-0.25,
			0,
			+0.25,
		},
	},

	{
		Name:            "Homing Missile Launcher",
		FireSound:       assets.AudioMissile1,
		Damage:          15,
		Reload:          3.5,
		Range:           600,
		ProjectileSpeed: 230,
		ProjectileImage: assets.ImageProjectileHomingMissile,
		ProjectileSize:  10,
		Explosion:       assets.ImageMissileImpact,
		ExplosionSound:  assets.AudioExplosion1,
		BurstSize:       2,
		Homing:          90,

		FireOffsets: []gmath.Vec{
			{Y: -8},
			{Y: +8},
		},
		ProjectileRotationDeltas: []gmath.Rad{
			-0.3,
			+0.3,
		},
	},
}
