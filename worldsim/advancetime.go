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
		for j := 0; j < 5; j++ {
			if r.updateWorld(0.2) {
				return false
			}
		}
	}
	return true
}

func (r *Runner) updateWorld(delta float64) bool {
	r.world.UpgradeRerollDelay = gmath.ClampMin(r.world.UpgradeRerollDelay-delta, 0)
	r.world.NextUpgradeDelay = gmath.ClampMin(r.world.NextUpgradeDelay-delta, 0)
	if r.world.UpgradeRerollDelay == 0 {
		r.world.UpgradeRerollDelay = float64(r.scene.Rand().IntRange(5, 15))
		r.world.UpgradeAvailable = gamedata.UpgradeKind(r.scene.Rand().IntRange(int(gamedata.FirstUpgrade), int(gamedata.LastUpgrade)))
	}

	player := r.world.Player
	planet := player.Planet

	encounterChance := 0.0
	switch r.world.Player.Mode {
	case gamedata.ModeJustEntered:
		encounterChance = 0.002
	case gamedata.ModeOrbiting:
		encounterChance = 0.04
	case gamedata.ModeScavenging:
		encounterChance = 0.07
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

	for _, p := range r.world.Planets {
		p.MineralsDelay = gmath.ClampMin(p.MineralsDelay-delta, 0)
	}
	return false
}
