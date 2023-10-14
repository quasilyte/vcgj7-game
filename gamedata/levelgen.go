package gamedata

func NewWorld() *World {
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
	}

	return w
}
