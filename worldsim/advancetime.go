package worldsim

import (
	"fmt"

	"github.com/quasilyte/gmath"
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
	case gamedata.ModeAttack:
		encounterChance = 1.0
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
			enemy := gamedata.CreateVesselDesign(r.scene.Rand(), r.world, enemyFaction)
			r.eventInfo = eventInfo{
				kind:  eventBattleInterrupt,
				enemy: enemy,
			}
			return true
		}
	}

	return false
}

func (r *Runner) updateFaction(delta float64, state *gamedata.FactionState) {
	if state.Tag == gamedata.FactionNone {
		return
	}

	state.AttackDelay = gmath.ClampMin(state.AttackDelay-delta, 0)
	state.CaptureDelay = gmath.ClampMin(state.CaptureDelay-delta, 0)

	if state.AttackDelay == 0 {
		if r.tryFactionAttack(state) {
			state.AttackDelay = r.scene.Rand().FloatRange(50, 150)
			return
		}
		state.AttackDelay = r.scene.Rand().FloatRange(15, 40)
	}
}

func (r *Runner) tryFactionAttack(state *gamedata.FactionState) bool {
	planet := randIterate(r.scene.Rand(), r.world.Planets, func(p *gamedata.Planet) bool {
		return p.Faction == state.Tag && p.VesselsByFaction[state.Tag] > r.scene.Rand().IntRange(5, 15)
	})
	if planet == nil {
		return false
	}

	largeSquad := false
	attackVessels := r.scene.Rand().IntRange(3, 6)
	if r.scene.Rand().Chance(0.4) && r.world.GameTime > 5*24 {
		largeSquad = true
		attackVessels *= 2
	}
	if attackVessels > planet.VesselsByFaction[state.Tag] {
		attackVessels = planet.VesselsByFaction[state.Tag] - r.scene.Rand().IntRange(2, 4)
	}

	targetPlanet := randIterate(r.scene.Rand(), r.world.Planets, func(p *gamedata.Planet) bool {
		if p.Faction == planet.Faction || p.Faction == gamedata.FactionNone {
			return false
		}
		dist := p.Info.MapOffset.DistanceTo(planet.Info.MapOffset)
		if dist > r.scene.Rand().FloatRange(50, 100) {
			return false
		}
		return true
	})
	if targetPlanet == nil {
		return false
	}

	if state.Tag != r.world.Player.Faction {
		if largeSquad {
			r.world.PushEvent(fmt.Sprintf("%s (controlled by %s) dispatched a large group of vessels", planet.Faction.Name(), planet.Info.Name))
		}
	} else {
		r.world.PushEvent(fmt.Sprintf("Allies start an attack operation on %s", targetPlanet.Info.Name))
	}

	speed := r.scene.Rand().FloatRange(5, 9)
	if largeSquad {
		speed *= 0.5
	}

	squad := &gamedata.Squad{
		NumVessels: attackVessels,
		Faction:    state.Tag,
		Speed:      speed,
		Dist:       planet.Info.MapOffset.DistanceTo(targetPlanet.Info.MapOffset),
		Dst:        targetPlanet,
	}
	r.world.Squads = append(r.world.Squads, squad)
	planet.VesselsByFaction[state.Tag] -= attackVessels

	return true
}

func (r *Runner) updateWorld(delta float64) bool {
	r.world.UpgradeRerollDelay = gmath.ClampMin(r.world.UpgradeRerollDelay-delta, 0)
	r.world.NextUpgradeDelay = gmath.ClampMin(r.world.NextUpgradeDelay-delta, 0)
	if r.world.UpgradeRerollDelay == 0 {
		r.world.UpgradeRerollDelay = float64(r.scene.Rand().IntRange(5, 15))
		r.world.UpgradeAvailable = gamedata.UpgradeKind(r.scene.Rand().IntRange(int(gamedata.FirstUpgrade), int(gamedata.LastUpgrade)))
	}

	for _, fs := range r.world.StateByFaction {
		r.updateFaction(delta, fs)
	}

	squads := r.world.Squads[:0]
	for _, squad := range r.world.Squads {
		squad.Dist -= delta * squad.Speed
		if squad.Dist <= 0 {
			squad.Dst.VesselsByFaction[squad.Faction] += squad.NumVessels
			continue
		}
		squads = append(squads, squad)
	}
	r.world.Squads = squads

	for _, p := range r.world.Planets {
		p.MineralsDelay = gmath.ClampMin(p.MineralsDelay-delta, 0)
		p.WeaponsRerollDelay = gmath.ClampMin(p.WeaponsRerollDelay-delta, 0)
		p.ShopSwapDelay = gmath.ClampMin(p.ShopSwapDelay-delta, 0)
		p.ResourceGenDelay = gmath.ClampMin(p.ResourceGenDelay-delta, 0)

		if p.WeaponsRerollDelay == 0 {
			p.WeaponsRerollDelay = r.scene.Rand().FloatRange(28, 40)
			r.rerollWeaponsSelection(p)
		}

		if p.ShopSwapDelay == 0 {
			p.ShopSwapDelay = r.scene.Rand().FloatRange(10, 15)
			p.ShopModeWeapons = r.scene.Rand().Bool()
		}

		if p.ResourceGenDelay == 0 {
			p.ResourceGenDelay = r.scene.Rand().FloatRange(30, 50)
			if p.Info.GasGiant {
				p.ResourceGenDelay *= 2
			}
			if p.Faction != gamedata.FactionNone {
				generated := r.scene.Rand().IntRange(1, 4)
				switch p.Faction {
				case gamedata.FactionB:
					generated *= 2
				case gamedata.FactionC:
					generated *= 3
				}
				if p.Faction != r.world.Player.Faction && r.scene.Rand().Chance(0.3) {
					generated += 5
				}
				p.MineralDeposit += generated
			}
		}

		if p.Faction == gamedata.FactionNone {
			continue
		}

		if p.VesselProduction {
			p.VesselProductionTime = gmath.ClampMin(p.VesselProductionTime-delta, 0)
			if p.VesselProductionTime == 0 {
				p.VesselsByFaction[p.Faction]++
				p.VesselProduction = false
			}
		} else {
			if p.MineralDeposit >= 50 && p.VesselsByFaction[p.Faction] < p.GarrisonLimit {
				cost := r.scene.Rand().IntRange(20, 50)
				p.MineralDeposit -= cost
				p.VesselProductionTime = float64(r.scene.Rand().IntRange(40, 100))
				p.VesselProduction = true
			}
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
