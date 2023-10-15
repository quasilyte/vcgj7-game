package gamedata

import (
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/vcgj7-game/assets"
)

func NewWorld(rand *gmath.Rand) *World {
	w := &World{}

	w.Player = &Player{
		Faction:  FactionA,
		VesselHP: 1.0,

		Mode: ModeOrbiting,

		MaxJumpDist: 60,
		JumpSpeed:   10,
		FuelUsage:   1.0,

		SpeedLevel:        1,
		AccelerationLevel: 1,
		RotationLevel:     1,
		EnergyLevel:       1,
		ArmorLevel:        1,

		Credits: 75,
		Fuel:    75,
		MaxFuel: 100,

		Cargo:    0,
		MaxCargo: 40,

		VesselDesign: &VesselDesign{
			Faction:       FactionA,
			Image:         assets.ImageVesselRaider,
			MaxHP:         100,
			MaxEnergy:     100,
			EnergyRegen:   2.5,
			MaxSpeed:      200,
			Acceleration:  140,
			RotationSpeed: 4,
			MainWeapon:    FindWeaponDesign("Ion Cannon"),
		},
	}

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

	planets[1].VesselsByFaction[FactionB] = 2

	for _, p := range planets {
		if p.Faction == FactionNone {
			continue
		}
		p.MineralDeposit = rand.IntRange(5, 200)
		numVessels := rand.IntRange(10, 20)
		if p.Faction != w.Player.Faction {
			numVessels += 10
		}
		p.VesselsByFaction[p.Faction] = numVessels
	}

	w.Player.Planet = planets[0]
	w.Planets = planets

	return w
}
