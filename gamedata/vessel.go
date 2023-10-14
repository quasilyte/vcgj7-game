package gamedata

import (
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/gmath"
)

type VesselDesign struct {
	Image resource.ImageID

	Faction Faction

	MaxHP     float64
	MaxEnergy float64

	EnergyRegen float64

	MaxSpeed     float64
	Acceleration float64

	RotationSpeed gmath.Rad

	MainWeapon      *WeaponDesign
	SecondaryWeapon *WeaponDesign
}
