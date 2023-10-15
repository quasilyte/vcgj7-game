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

	Cost int

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
	Primary          bool
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
		Name:             "Photon Cannon",
		Cost:             100,
		FireSound:        assets.AudioPhotonCannon1,
		Damage:           12,
		Reload:           0.9,
		EnergyCost:       6,
		EnergyConversion: 3.0,
		Range:            1600,
		ProjectileSpeed:  550,
		ProjectileImage:  assets.ImageProjectilePhotonCannon,
		ProjectileSize:   10,
		Explosion:        assets.ImagePhotonCannonImpact,
		// ExplosionSound:   assets.AudioIonCannonImpact,
		Blockable: true,
		Primary:   true,
		BurstSize: 1,

		FireOffsets:              []gmath.Vec{{}},
		ProjectileRotationDeltas: []gmath.Rad{0},
	},

	{
		Name:             "Pulse Laser",
		Cost:             210,
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
		Primary:          true,
		BurstSize:        1,

		FireOffsets:              []gmath.Vec{{}},
		ProjectileRotationDeltas: []gmath.Rad{0},
	},

	{
		Name:             "Ion Cannon",
		Cost:             150,
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
		Primary:          true,
		BurstSize:        1,

		FireOffsets:              []gmath.Vec{{}},
		ProjectileRotationDeltas: []gmath.Rad{0},
	},

	{
		Name:             "Assault Laser",
		Cost:             440,
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
		Primary:   true,
		BurstSize: 2,

		FireOffsets: []gmath.Vec{
			{Y: -6},
			{Y: +6},
		},
		ProjectileRotationDeltas: []gmath.Rad{0, 0},
	},

	{
		Name:             "Scatter Gun",
		Cost:             520,
		FireSound:        assets.AudioScatterGun1,
		Damage:           7,
		Reload:           0.4,
		EnergyCost:       9,
		EnergyConversion: 1.2,
		Range:            350,
		ProjectileSpeed:  450,
		ProjectileImage:  assets.ImageProjectileScatterGun,
		ProjectileSize:   6,
		Explosion:        assets.ImageScatterGunImpact,
		// ExplosionSound:   assets.AudioIonCannonImpact,
		Blockable: true,
		Primary:   true,
		BurstSize: 7,

		FireOffsets: []gmath.Vec{
			{Y: -9},
			{Y: -6},
			{Y: -3},
			{},
			{Y: +3},
			{Y: +6},
			{Y: +9},
		},
		ProjectileRotationDeltas: []gmath.Rad{-0.45, -0.3, -0.15, 0, +0.15, +0.3, +0.45},
	},

	{
		Name:             "Trident",
		Cost:             650,
		FireSound:        assets.AudioTrident1,
		Damage:           14,
		Reload:           0.5,
		EnergyCost:       14,
		EnergyConversion: 1.8,
		Range:            380,
		ProjectileSpeed:  350,
		ProjectileImage:  assets.ImageProjectileTrident,
		ProjectileSize:   8,
		Explosion:        assets.ImageTridentImpact,
		Primary:          true,
		Blockable:        true,
		BurstSize:        3,

		FireOffsets: []gmath.Vec{
			{Y: -25},
			{},
			{Y: +25},
		},
		ProjectileRotationDeltas: []gmath.Rad{+0.025, 0, -0.025},
	},

	{
		Name:             "Lance",
		Cost:             900,
		FireSound:        assets.AudioLance1,
		Damage:           25,
		Reload:           0.7,
		EnergyCost:       12,
		EnergyConversion: 2.5,
		Range:            2600,
		ProjectileSpeed:  600,
		ProjectileImage:  assets.ImageProjectileLance,
		ProjectileSize:   10,
		Explosion:        assets.ImageLanceImpact,
		// ExplosionSound:   assets.AudioIonCannonImpact,
		Blockable: true,
		Primary:   true,
		BurstSize: 1,

		FireOffsets:              []gmath.Vec{{}},
		ProjectileRotationDeltas: []gmath.Rad{0},
	},

	{
		Name:            "Missile Launcher",
		Cost:            200,
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
		Cost:            280,
		FireSound:       assets.AudioMissile1,
		Damage:          15,
		Reload:          4.0,
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

	{
		Name:            "Torpedo Launcher",
		Cost:            400,
		FireSound:       assets.AudioMissile1,
		Damage:          35,
		Reload:          6.0,
		Range:           1900,
		ProjectileSpeed: 185,
		ProjectileImage: assets.ImageProjectileTorpedo,
		ProjectileSize:  10,
		Explosion:       assets.ImageMissileImpact,
		ExplosionSound:  assets.AudioExplosion1,
		BurstSize:       1,
		Homing:          120,

		FireOffsets:              []gmath.Vec{{}},
		ProjectileRotationDeltas: []gmath.Rad{0},
	},
}
