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
		JumpSpeed:   8,
		FuelUsage:   1.0,

		SpeedLevel:        1,
		AccelerationLevel: 1,
		RotationLevel:     1,
		EnergyLevel:       1,
		ArmorLevel:        1,

		Credits: 75,
		Fuel:    100,
		MaxFuel: 130,

		Cargo:    0,
		MaxCargo: 40,

		VesselDesign: &VesselDesign{
			Faction:       FactionA,
			Image:         assets.ImageVesselPlayer,
			MaxHP:         120,
			MaxEnergy:     90,
			EnergyRegen:   1.5,
			MaxSpeed:      150,
			Acceleration:  90,
			RotationSpeed: 2.5,
			MainWeapon:    FindWeaponDesign("Photon Cannon"),
			// MainWeapon:      FindWeaponDesign("Lance"),
			// SecondaryWeapon: FindWeaponDesign("Firestorm"),
		},
	}

	planets := make([]*Planet, len(Planets))
	for i := range planets {
		p := &Planet{
			Info:          Planets[i],
			GarrisonLimit: rand.IntRange(25, 40),
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
		p.AttackDelay = rand.FloatRange(75, 100)
		p.CaptureDelay = rand.FloatRange(160, 200)
		numVessels := rand.IntRange(4, 8)
		if p.Faction != w.Player.Faction {
			numVessels += 10
		}
		p.VesselsByFaction[p.Faction] = numVessels
	}

	w.Player.Planet = planets[0]
	w.Planets = planets

	w.PushEvent("All three major factions declare war to each other")

	w.Artifacts = []string{
		// Passively generates fuel over time.
		"Fuel Generator",
		// Repairs vessel hull over time.
		"Repair Bots",
		// Faster scanning time.
		"Scantide",
		// More rewards in some situations.
		"Lucky Charm",
	}

	return w
}
