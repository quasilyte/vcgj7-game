package worldsim

import (
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/vcgj7-game/assets"
	"github.com/quasilyte/vcgj7-game/gamedata"
)

func (r *Runner) AdvanceTime(hours int) bool {
	for i := 0; i < hours; i++ {
		r.world.GameTime++
		// One in-game hour is simulated during 1 second in delta time terms.
		if r.processEncounters() {
			return false
		}
		for j := 0; j < 5; j++ {
			if r.updateWorld(0.2) {
				return false
			}
		}
	}
	return true
}

func (r *Runner) processEncounters() bool {
	player := r.world.Player
	planet := player.Planet

	encounterChance := 0.0
	switch r.world.Player.Mode {
	case gamedata.ModeJustEntered:
		encounterChance = 0.05
	case gamedata.ModeOrbiting:
		encounterChance = 0.1
	case gamedata.ModeScavenging:
		encounterChance = 0.2
	}
	switch {
	case planet.Faction == player.Faction:
		encounterChance *= 0.25
	case planet.Faction == gamedata.FactionNone:
		encounterChance *= 0.65
	}
	if encounterChance > 0 && r.scene.Rand().Chance(encounterChance) {
		// If there is any hostile vessels around here, the battle will start.
		r.encounterOptions = r.encounterOptions[:0]
		for i, num := range planet.VesselsByFaction {
			if num == 0 {
				continue
			}
			f := gamedata.Faction(i)
			if f == player.Faction {
				continue
			}
			r.encounterOptions = append(r.encounterOptions, f)
		}
		if len(r.encounterOptions) != 0 {
			enemyFaction := gmath.RandElem(r.scene.Rand(), r.encounterOptions)
			r.eventInfo = eventInfo{
				kind: eventBattleInterrupt,
				enemy: &gamedata.VesselDesign{
					Faction:         enemyFaction,
					Image:           assets.ImageVesselMarauder,
					MaxHP:           150,
					MaxEnergy:       120,
					EnergyRegen:     3.0,
					MaxSpeed:        180,
					Acceleration:    90,
					RotationSpeed:   2.5,
					MainWeapon:      gamedata.FindWeaponDesign("Pulse Laser"),
					SecondaryWeapon: gamedata.FindWeaponDesign("Homing Missile Launcher"),
				},
			}
			return true
		}
	}

	return false
}

func (r *Runner) updateWorld(delta float64) bool {
	r.world.UpgradeRerollDelay = gmath.ClampMin(r.world.UpgradeRerollDelay-delta, 0)
	r.world.NextUpgradeDelay = gmath.ClampMin(r.world.NextUpgradeDelay-delta, 0)
	if r.world.UpgradeRerollDelay == 0 {
		r.world.UpgradeRerollDelay = float64(r.scene.Rand().IntRange(5, 15))
		r.world.UpgradeAvailable = gamedata.UpgradeKind(r.scene.Rand().IntRange(int(gamedata.FirstUpgrade), int(gamedata.LastUpgrade)))
	}

	for _, p := range r.world.Planets {
		p.MineralsDelay = gmath.ClampMin(p.MineralsDelay-delta, 0)
		p.WeaponsRerollDelay = gmath.ClampMin(p.WeaponsRerollDelay-delta, 0)
		p.ShopSwapDelay = gmath.ClampMin(p.ShopSwapDelay-delta, 0)
		if p.WeaponsRerollDelay == 0 {
			p.WeaponsRerollDelay = r.scene.Rand().FloatRange(28, 40)
			r.rerollWeaponsSelection(p)
		}
		if p.ShopSwapDelay == 0 {
			p.ShopSwapDelay = r.scene.Rand().FloatRange(10, 15)
			p.ShopModeWeapons = r.scene.Rand().Bool()
		}
	}
	return false
}

func (r *Runner) rerollWeaponsSelection(p *gamedata.Planet) {
	p.WeaponsAvailable = p.WeaponsAvailable[:0]
	if r.scene.Rand().Chance(0.05) {
		return // No weapons available
	}
	numWeapons := r.scene.Rand().IntRange(2, 3)
	weapons := make([]*gamedata.WeaponDesign, len(gamedata.Weapons))
	copy(weapons, gamedata.Weapons)
	gmath.Shuffle(r.scene.Rand(), weapons)
	for _, w := range weapons[:numWeapons] {
		p.WeaponsAvailable = append(p.WeaponsAvailable, w.Name)
	}
}
