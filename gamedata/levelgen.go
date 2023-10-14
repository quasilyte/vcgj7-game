package gamedata

import (
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/vcgj7-game/assets"
)

func NewWorld(rand *gmath.Rand) *World {
	w := &World{}

	planets := make([]*Planet, len(Planets))
	for i := range planets {
		p := &Planet{
			Info: Planets[i],
		}
		planets[i] = p
	}

	planets[0].Faction = FactionA
	planets[2].Faction = FactionB
	planets[7].Faction = FactionC

	for _, p := range planets {
		if p.Faction == FactionNone {
			continue
		}
		p.MineralDeposit = rand.IntRange(5, 200)
	}

	w.Planets = planets

	w.Player = &Player{
		Faction:  FactionA,
		Planet:   planets[0],
		VesselHP: 1.0,

		Mode: ModeOrbiting,

		MaxJumpDist: 60,
		JumpSpeed:   10,
		FuelUsage:   1.0,

		Credits: 150,
		Fuel:    75,
		MaxFuel: 100,

		Cargo:    0,
		MaxCargo: 40,

		VesselDesign: &VesselDesign{
			Image:           assets.ImageVesselRaider,
			MaxHP:           100,
			MaxEnergy:       100,
			EnergyRegen:     2.5,
			MaxSpeed:        200,
			Acceleration:    140,
			RotationSpeed:   4,
			MainWeapon:      FindWeaponDesign("Assault Laser"),
			SecondaryWeapon: FindWeaponDesign("Missile Launcher"),
		},
	}

	return w
}
