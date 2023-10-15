package gamedata

import (
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/vcgj7-game/assets"
)

func CreateVesselDesign(rand *gmath.Rand, world *World, faction Faction) *VesselDesign {
	challenge := chooseBattleChallenge(rand, world)
	eliteVessel := challenge >= 1 && rand.Chance(0.2)
	design := &VesselDesign{
		Faction:     faction,
		Challenge:   challenge,
		Elite:       eliteVessel,
		MaxEnergy:   75 + float64(rand.IntRange(5, 50)) + float64(challenge*15),
		EnergyRegen: 1.25 + rand.FloatRange(0.1, 0.3) + (float64(challenge) * 0.2),
	}

	switch challenge {
	case 0:
		assignChallenge0weapons(rand, design)
	case 1:
		if rand.Chance(0.2) {
			assignChallenge0weapons(rand, design)
		} else {
			assignChallenge1weapons(rand, design)
		}
	case 2:
		roll := rand.Float()
		if roll <= 0.1 {
			assignChallenge0weapons(rand, design)
		} else if roll < 0.3 {
			assignChallenge1weapons(rand, design)
		} else {
			assignChallenge2weapons(rand, design)
		}
	case 3:
		roll := rand.Float()
		if roll <= 0.1 {
			assignChallenge1weapons(rand, design)
		} else if roll < 0.3 {
			assignChallenge2weapons(rand, design)
		} else {
			assignChallenge3weapons(rand, design)
		}
	}

	switch faction {
	default:
		panic("unexpected faction")
	case FactionB: // Beta
		design.MaxHP = float64(rand.IntRange(70, 110)) + float64(challenge*25)
		design.MaxSpeed = float64(rand.IntRange(180, 240))
		design.Acceleration = float64(rand.IntRange(40, 50))
		design.RotationSpeed = gmath.Rad(rand.FloatRange(1.4, 2.0))
		if eliteVessel {
			design.RotationSpeed -= gmath.Rad(rand.FloatRange(0.2, 0.6))
			design.MaxSpeed -= float64(rand.IntRange(20, 40))
			design.MaxHP += float64(rand.IntRange(40, 80))
			design.Image = assets.ImageVesselBetaBig
		} else {
			design.Image = assets.ImageVesselBetaSmall
		}
	case FactionC: // Gamma
		design.MaxHP = float64(rand.IntRange(110, 140)) + float64(challenge*35)
		design.MaxSpeed = float64(rand.IntRange(90, 120))
		design.Acceleration = float64(rand.IntRange(100, 140))
		design.RotationSpeed = gmath.Rad(rand.FloatRange(2.2, 2.8))
		if eliteVessel {
			design.RotationSpeed -= gmath.Rad(rand.FloatRange(0.2, 0.6))
			design.MaxSpeed -= float64(rand.IntRange(20, 40))
			design.MaxHP += float64(rand.IntRange(70, 120))
			design.Image = assets.ImageVesselGammaBig
		} else {
			design.Image = assets.ImageVesselGammaSmall
		}
	}

	return design
}

func chooseBattleChallenge(rand *gmath.Rand, world *World) int {
	// Challenges are in 0-3 range.
	if world.Player.Battles < 3 {
		return 0
	}
	if world.Player.Battles < 5 {
		if rand.Chance(0.6) {
			return 1
		}
		return 0
	}
	if world.Player.Battles < 9 {
		if rand.Chance(0.6) {
			return 2
		}
		if rand.Chance(0.6) {
			return 1
		}
		return 0
	}
	if world.Player.Battles < 15 {
		if rand.Chance(0.4) {
			return 3
		}
		if rand.Chance(0.6) {
			return 2
		}
		if rand.Chance(0.6) {
			return 1
		}
		return 0
	}
	if rand.Chance(0.7) {
		return 3
	}
	if rand.Chance(0.6) {
		return 2
	}
	if rand.Chance(0.8) {
		return 1
	}
	return 0
}

func assignChallenge0weapons(rand *gmath.Rand, design *VesselDesign) {
	if rand.Chance(0.7) {
		// A single primary weapon.
		roll := rand.Float()
		if roll <= 0.6 {
			design.MainWeapon = FindWeaponDesign("Ion Cannon")
		} else if roll <= 0.9 {
			design.MainWeapon = FindWeaponDesign("Pulse Laser")
		} else {
			design.MainWeapon = FindWeaponDesign("Photon Cannon")
		}
	} else {
		// A single secondary weapon.
		if rand.Chance(0.6) {
			design.SecondaryWeapon = FindWeaponDesign("Missile Launcher")
		} else {
			design.SecondaryWeapon = FindWeaponDesign("Homing Missile Launcher")
		}
	}
}

func assignChallenge1weapons(rand *gmath.Rand, design *VesselDesign) {
	roll := rand.Float()
	if roll <= 0.4 {
		design.MainWeapon = FindWeaponDesign("Ion Cannon")
	} else if roll <= 0.7 {
		design.MainWeapon = FindWeaponDesign("Pulse Laser")
	} else if roll <= 0.85 {
		design.MainWeapon = FindWeaponDesign("Scatter Gun")
	} else {
		design.MainWeapon = FindWeaponDesign("Assault Laser")
	}

	if rand.Chance(0.6) {
		design.SecondaryWeapon = FindWeaponDesign("Homing Missile Launcher")
	} else {
		design.SecondaryWeapon = FindWeaponDesign("Missile Launcher")
	}
}

func assignChallenge2weapons(rand *gmath.Rand, design *VesselDesign) {
	roll := rand.Float()
	if roll <= 0.4 {
		design.MainWeapon = FindWeaponDesign("Assault Laser")
	} else if roll <= 0.7 {
		design.MainWeapon = FindWeaponDesign("Scatter Gun")
	} else {
		design.MainWeapon = FindWeaponDesign("Trident")
	}

	roll = rand.Float()
	if roll <= 0.4 {
		design.SecondaryWeapon = FindWeaponDesign("Homing Missile Launcher")
	} else if roll <= 0.6 {
		design.SecondaryWeapon = FindWeaponDesign("Missile Launcher")
	} else if roll <= 0.95 {
		design.SecondaryWeapon = FindWeaponDesign("Torpedo Launcher")
	} else {
		design.SecondaryWeapon = FindWeaponDesign("Firestorm")
	}
}

func assignChallenge3weapons(rand *gmath.Rand, design *VesselDesign) {
	roll := rand.Float()
	if roll <= 0.3 {
		design.MainWeapon = FindWeaponDesign("Assault Laser")
	} else if roll <= 0.5 {
		design.MainWeapon = FindWeaponDesign("Scatter Gun")
	} else if roll <= 0.9 {
		design.MainWeapon = FindWeaponDesign("Trident")
	} else {
		design.MainWeapon = FindWeaponDesign("Lance")
	}

	roll = rand.Float()
	if roll <= 0.4 {
		design.SecondaryWeapon = FindWeaponDesign("Torpedo Launcher")
	} else if roll <= 0.7 {
		design.SecondaryWeapon = FindWeaponDesign("Firestorm")
	} else {
		design.SecondaryWeapon = FindWeaponDesign("Homing Missile Launcher")
	}
}
