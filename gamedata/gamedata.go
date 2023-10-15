package gamedata

import (
	"github.com/quasilyte/gmath"
)

type BattleRewards struct {
	Victory bool

	Experience int
	Cargo      int
	Credits    int
	Fuel       int
}

type World struct {
	Player *Player

	Planets []*Planet

	GameTime int // In hours

	NextUpgradeDelay   float64
	UpgradeRerollDelay float64
	UpgradeAvailable   UpgradeKind
}

type UpgradeKind int

const (
	UpgradeUnknown UpgradeKind = iota
	UpgradeMaxFuel
	UpgradeMaxCargo
	UpgradeJumpSpeed
	UpgradeJumpMaxDistance
	_numUpgrades
)

const (
	FirstUpgrade = UpgradeUnknown + 1
	LastUpgrade  = _numUpgrades - 1
)

type Player struct {
	Planet *Planet

	Faction Faction

	BattleRewards BattleRewards

	Mode Mode

	SpeedLevel        int
	AccelerationLevel int
	RotationLevel     int
	EnergyLevel       int
	ArmorLevel        int

	VesselDesign *VesselDesign
	VesselHP     float64 // percentage

	JumpSpeed   float64
	MaxJumpDist float64
	FuelUsage   float64

	Experience int
	Credits    int
	Fuel       int
	MaxFuel    int
	Cargo      int
	MaxCargo   int
}

func (p *Player) HasWeapon(w *WeaponDesign) bool {
	if w.Primary {
		return p.VesselDesign.MainWeapon == w
	}
	return p.VesselDesign.SecondaryWeapon == w
}

func (p *Player) FreeCargoSpace() int {
	return p.MaxCargo - p.Cargo
}

func (p *Player) LoadCargo(amount int) int {
	freeSpace := p.FreeCargoSpace()
	if amount > freeSpace {
		amount = freeSpace
	}
	p.Cargo += amount
	return amount
}

type Planet struct {
	Faction Faction

	Info *PlanetInfo

	VesselProduction     bool
	VesselProductionTime float64

	GarrisonLimit int

	MineralsDelay  float64
	MineralDeposit int

	VesselsByFaction [NumFactions]int

	ShopModeWeapons bool
	ShopSwapDelay   float64

	WeaponsRerollDelay float64
	WeaponsAvailable   []string

	AreasVisited PlanetVisitStatus
}

type PlanetVisitStatus struct {
	VisitedMineralsMarket bool
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
