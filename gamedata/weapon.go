package gamedata

import (
	"fmt"

	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/vcgj7-game/assets"
)

type WeaponDesign struct {
	Name string

	FireSound resource.AudioID

	Damage float64
	Reload float64

	Range           float64
	ProjectileSpeed float64
	ProjectileSize  float64
	ProjectileImage resource.ImageID

	Explosion resource.ImageID

	EnergyCost       float64
	EnergyConversion float64
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
		Reload:           0.4,
		EnergyCost:       6,
		EnergyConversion: 2.0,
		Range:            250,
		ProjectileSpeed:  280,
		ProjectileImage:  assets.ImageProjectilePulseLaser,
		ProjectileSize:   6,
	},

	{
		Name:             "Ion Cannon",
		FireSound:        assets.AudioIonCannon1,
		Damage:           10,
		Reload:           0.9,
		EnergyCost:       5,
		EnergyConversion: 0.5,
		Range:            300,
		ProjectileSpeed:  320,
		ProjectileImage:  assets.ImageProjectileIonCannon,
		ProjectileSize:   8,
		Explosion:        assets.ImageIonCannonImpact,
	},
}
