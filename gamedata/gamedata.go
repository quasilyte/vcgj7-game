package gamedata

type World struct {
	Player *Player

	Planets []*Planet
}

type Player struct {
	Planet *Planet

	VesselDesign *VesselDesign

	Experience int
}

type Planet struct {
	Info *PlanetInfo
}

type PlanetInfo struct {
	Name     string
	RealName string
	GasGiant bool
}

var Planets = []*PlanetInfo{
	{
		Name:     "Planet VIII",
		RealName: "Neptune",
		GasGiant: true,
	},

	{
		Name:     "Planet VII",
		RealName: "Uranus",
		GasGiant: true,
	},

	{
		Name:     "Planet VI",
		RealName: "Saturn",
		GasGiant: true,
	},

	{
		Name:     "Planet V",
		RealName: "Jupiter",
		GasGiant: true,
	},

	{
		Name:     "Planet IV",
		RealName: "Mars",
	},

	{
		Name:     "Planet III",
		RealName: "Earth",
	},

	{
		Name:     "Planet II",
		RealName: "Venus",
	},

	{
		Name:     "Planet I",
		RealName: "Mercury",
	},
}
