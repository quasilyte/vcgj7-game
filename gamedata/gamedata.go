package gamedata

import "github.com/quasilyte/gmath"

type World struct {
	Player *Player

	Planets []*Planet
}

type Player struct {
	Planet *Planet

	VesselDesign *VesselDesign
	VesselHP     float64 // percentage

	Experience int
	Fuel       int
}

type Planet struct {
	Info *PlanetInfo
}

type PlanetInfo struct {
	Name      string
	RealName  string
	GasGiant  bool
	MapOffset gmath.Vec
}

var Planets = []*PlanetInfo{
	{
		Name:      "Planet VIII",
		RealName:  "Neptune",
		GasGiant:  true,
		MapOffset: gmath.Vec{X: 14, Y: 37},
	},

	{
		Name:      "Planet VII",
		RealName:  "Uranus",
		GasGiant:  true,
		MapOffset: gmath.Vec{X: 34, Y: 19},
	},

	{
		Name:      "Planet VI",
		RealName:  "Saturn",
		GasGiant:  true,
		MapOffset: gmath.Vec{X: 39, Y: 125},
	},

	{
		Name:      "Planet V",
		RealName:  "Jupiter",
		GasGiant:  true,
		MapOffset: gmath.Vec{X: 38, Y: 80},
	},

	{
		Name:      "Planet IV",
		RealName:  "Mars",
		MapOffset: gmath.Vec{X: 95, Y: 30},
	},

	{
		Name:      "Planet III",
		RealName:  "Earth",
		MapOffset: gmath.Vec{X: 80, Y: 67},
	},

	{
		Name:      "Planet II",
		RealName:  "Venus",
		MapOffset: gmath.Vec{X: 97, Y: 94},
	},

	{
		Name:      "Planet I",
		RealName:  "Mercury",
		MapOffset: gmath.Vec{X: 125, Y: 102},
	},
}
