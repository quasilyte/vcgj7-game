package worldsim

import (
	"fmt"
	"math"

	"github.com/quasilyte/gmath"
	"github.com/quasilyte/vcgj7-game/assets"
	"github.com/quasilyte/vcgj7-game/gamedata"
)

func (r *Runner) checkVictory() {
	others := false
	for _, p := range r.world.Planets {
		if p.Faction == gamedata.FactionNone {
			continue
		}
		if p.Faction != r.world.Player.Faction {
			others = true
			break
		}
	}
	if !others {
		r.EventGameOver.Emit(true)
	}
}

func (r *Runner) AdvanceTime(hours int) bool {
	player := r.world.Player

	canRegen := false
	switch player.Mode {
	case gamedata.ModeJustEntered, gamedata.ModeOrbiting, gamedata.ModeScavenging, gamedata.ModeSneaking:
		canRegen = true
	}

	for i := 0; i < hours; i++ {
		r.world.GameTime++
		r.checkVictory()

		if r.world.GameTime%24 == 0 {
			salary := gamedata.GetSalary(player.Experience) + player.ExtraSalary
			player.Credits += salary
		}

		if player.HasArtifact("Fuel Generator") && canRegen {
			if r.scene.Rand().Chance(0.6) {
				player.Fuel = gmath.ClampMax(player.Fuel+1, player.MaxFuel)
			}
		}
		if player.HasArtifact("Repair Bots") && canRegen {
			if r.scene.Rand().Chance(0.8) {
				player.VesselHP = gmath.ClampMax(player.VesselHP+0.02, 1.0)
			}
		}

		// One in-game hour is simulated during 1 second in delta time terms.
		if r.processEncounters() {
			return false
		}
		for j := 0; j < 5; j++ {
			if r.updateWorld(0.2) {
				return false
			}
		}
		for _, p := range r.world.Planets {
			r.processPlanetBattles(p)
		}
	}
	return true
}

func (r *Runner) makePirate() *gamedata.VesselDesign {
	pirate := &gamedata.VesselDesign{
		Image:         assets.ImageVesselPirate,
		MaxHP:         float64(r.scene.Rand().IntRange(120, 150) + (r.world.PirateSeq * 50)),
		MaxEnergy:     float64(r.scene.Rand().IntRange(200, 300) + (r.world.PirateSeq * 30)),
		EnergyRegen:   2.0,
		MaxSpeed:      200,
		Acceleration:  80,
		Challenge:     2,
		RotationSpeed: 0.5,
	}
	if r.scene.Rand().Chance(0.8) {
		pirate.MainWeapon = gamedata.FindWeaponDesign("Scatter Gun")
	} else {
		pirate.MainWeapon = gamedata.FindWeaponDesign("Trident")
	}
	return pirate
}

func (r *Runner) processEncounters() bool {
	player := r.world.Player

	if r.world.NextPirateDelay == 0 && r.world.PirateSeq < 3 {
		if player.VesselHP >= 0.8 && player.Battles >= 2 {
			r.world.NextPirateDelay = r.scene.Rand().FloatRange(600, 1200)
			r.eventInfo = eventInfo{
				kind:  eventBattleInterrupt,
				enemy: r.makePirate(),
			}
			return true
		}

		r.world.NextPirateDelay = r.scene.Rand().FloatRange(20, 40)
	}

	planet := player.Planet

	encounterChance := 0.0
	switch r.world.Player.Mode {
	case gamedata.ModeSneaking:
		encounterChance = 0.01
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

func (r *Runner) processPlanetBattles(p *gamedata.Planet) {
	r.planetFactions = r.planetFactions[:0]
	for i, num := range p.VesselsByFaction {
		if num == 0 {
			continue
		}
		f := gamedata.Faction(i)
		r.planetFactions = append(r.planetFactions, f)
	}
	if len(r.planetFactions) < 2 {
		return
	}
	battleChance := 0.45
	if !r.scene.Rand().Chance(battleChance) {
		return
	}
	gmath.Shuffle(r.scene.Rand(), r.planetFactions)
	faction1 := r.planetFactions[0]
	faction2 := r.planetFactions[1]
	loser := faction1
	winner := faction2
	if r.scene.Rand().Bool() {
		loser, winner = winner, loser
	}
	p.VesselsByFaction[loser]--

	if p.Faction == loser && p.VesselsByFaction[loser] == 0 {
		if p.Faction == r.world.Player.Faction {
			r.world.PushEvent(fmt.Sprintf("We lost control over %s", p.Info.Name))
		} else {
			if winner == r.world.Player.Faction {
				r.world.PushEvent(fmt.Sprintf("%s is liberated from the enemy forces", p.Info.Name))
			} else {
				r.world.PushEvent(fmt.Sprintf("%s lost %s to %s", winner.Name(), p.Info.Name, loser.Name()))
			}
		}
		p.Faction = gamedata.FactionNone
		p.VesselProduction = false
		p.VesselProductionTime = 0
	}
}

func (r *Runner) processPlanetActions(p *gamedata.Planet) {
	if p.AttackDelay == 0 {
		numVessels := p.VesselsByFaction[p.Faction]
		if numVessels < 10 && r.scene.Rand().Chance(0.8) {
			p.AttackDelay = r.scene.Rand().FloatRange(60, 100)
			return
		}
		if numVessels < 20 && r.scene.Rand().Chance(0.5) {
			p.AttackDelay = r.scene.Rand().FloatRange(20, 150)
			return
		}
		if r.tryFactionAttack(p) {
			p.AttackDelay = r.scene.Rand().FloatRange(70, 300)
			return
		}
		p.AttackDelay = r.scene.Rand().FloatRange(20, 40)
		return
	}

	if p.CaptureDelay == 0 {
		numVessels := p.VesselsByFaction[p.Faction]
		if numVessels < 10 && r.scene.Rand().Chance(0.9) {
			p.CaptureDelay = r.scene.Rand().FloatRange(60, 100)
			return
		}
		if r.tryFactionCapture(p) {
			p.CaptureDelay = r.scene.Rand().FloatRange(150, 400)
			return
		}
		p.CaptureDelay = r.scene.Rand().FloatRange(40, 80)
		return
	}
}

func (r *Runner) tryFactionAttack(planet *gamedata.Planet) bool {
	if planet.VesselsByFaction[planet.Faction] <= r.scene.Rand().IntRange(5, 15) {
		return false
	}

	largeSquad := false
	attackVessels := r.scene.Rand().IntRange(3, 6)
	if r.scene.Rand().Chance(0.4) && r.world.GameTime > 5*24 {
		largeSquad = true
		attackVessels *= 2
	}
	if attackVessels > planet.VesselsByFaction[planet.Faction] {
		attackVessels = planet.VesselsByFaction[planet.Faction] - r.scene.Rand().IntRange(2, 4)
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

	if planet.Faction != r.world.Player.Faction {
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
		Faction:    planet.Faction,
		Speed:      speed,
		Dist:       planet.Info.MapOffset.DistanceTo(targetPlanet.Info.MapOffset),
		Dst:        targetPlanet,
	}
	r.world.Squads = append(r.world.Squads, squad)
	planet.VesselsByFaction[planet.Faction] -= attackVessels
	return true
}

func (r *Runner) tryFactionCapture(planet *gamedata.Planet) bool {
	if planet.VesselsByFaction[planet.Faction] <= r.scene.Rand().IntRange(5, 10) {
		return false
	}

	attackVessels := r.scene.Rand().IntRange(1, 3)

	targetPlanet := randIterate(r.scene.Rand(), r.world.Planets, func(p *gamedata.Planet) bool {
		if p.Faction != gamedata.FactionNone {
			return false
		}
		dist := p.Info.MapOffset.DistanceTo(planet.Info.MapOffset)
		if dist > r.scene.Rand().FloatRange(70, 110) {
			return false
		}
		return true
	})
	if targetPlanet == nil {
		return false
	}

	speed := r.scene.Rand().FloatRange(6, 11)
	squad := &gamedata.Squad{
		NumVessels: attackVessels,
		Faction:    planet.Faction,
		Speed:      speed,
		Dist:       planet.Info.MapOffset.DistanceTo(targetPlanet.Info.MapOffset),
		Dst:        targetPlanet,
	}
	r.world.Squads = append(r.world.Squads, squad)
	planet.VesselsByFaction[planet.Faction] -= attackVessels
	return true
}

func (r *Runner) maybeRollQuest() {
	r.alliedPlanets = r.alliedPlanets[:0]
	for _, p := range r.world.Planets {
		if p.Faction != r.world.Player.Faction {
			continue
		}
		r.alliedPlanets = append(r.alliedPlanets, p)
	}
	if len(r.alliedPlanets) < 2 {
		return
	}
	gmath.Shuffle(r.scene.Rand(), r.alliedPlanets)
	r.world.CurrentQuest = &gamedata.Quest{
		Active:        false,
		Giver:         r.alliedPlanets[0],
		Receiver:      r.alliedPlanets[1],
		CreditsReward: r.scene.Rand().IntRange(20, 200),
		ExpReward:     r.scene.Rand().IntRange(10, 60),
	}
	fmt.Println("quest rolled", r.alliedPlanets[0].Info.Name, "=>", r.alliedPlanets[1].Info.Name)
}

func (r *Runner) updateWorld(delta float64) bool {
	r.world.NextPirateDelay = gmath.ClampMin(r.world.NextPirateDelay-delta, 0)
	r.world.QuestRerollDelay = gmath.ClampMin(r.world.QuestRerollDelay-delta, 0)
	r.world.UpgradeRerollDelay = gmath.ClampMin(r.world.UpgradeRerollDelay-delta, 0)
	r.world.NextUpgradeDelay = gmath.ClampMin(r.world.NextUpgradeDelay-delta, 0)
	if r.world.UpgradeRerollDelay == 0 {
		r.world.UpgradeRerollDelay = float64(r.scene.Rand().IntRange(5, 15))
		r.world.UpgradeAvailable = gamedata.UpgradeKind(r.scene.Rand().IntRange(int(gamedata.FirstUpgrade), int(gamedata.LastUpgrade)))
	}

	if r.world.CurrentQuest != nil && r.world.CurrentQuest.Active {
		q := r.world.CurrentQuest
		p := r.world.Player
		if q.Giver.Faction != p.Faction || q.Receiver.Faction != p.Faction {
			// Quest failed.
			r.world.CurrentQuest = nil
		}
	}

	if r.world.QuestRerollDelay == 0 {
		r.world.QuestRerollDelay = float64(r.scene.Rand().IntRange(60, 130))
		if r.world.CurrentQuest != nil && !r.world.CurrentQuest.Active {
			r.world.CurrentQuest = nil
		}
		if r.world.CurrentQuest == nil {
			r.maybeRollQuest()
		}
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
		p.AttackDelay = gmath.ClampMin(p.AttackDelay-delta, 0)
		p.CaptureDelay = gmath.ClampMin(p.CaptureDelay-delta, 0)

		if p.Faction == gamedata.FactionNone {
			for i := range p.InfluenceByFaction {
				p.InfluenceByFaction[i] = gmath.ClampMin(p.InfluenceByFaction[i]-delta, 0)
			}
			faction := gamedata.FactionNone
			numFactions := 0
			for i, num := range p.VesselsByFaction {
				if num == 0 {
					continue
				}
				numFactions++
				faction = gamedata.Faction(i)
			}
			if numFactions == 1 {
				numVessels := p.VesselsByFaction[faction]
				v := math.Log(float64(numVessels)) + 1.0
				p.InfluenceByFaction[faction] += v * delta
				// 1 vessels (v=1.000) capture in 30.000 days
				// 2 vessels (v=1.693) capture in 17.718 days
				// 3 vessels (v=2.099) capture in 14.295 days
				// 4 vessels (v=2.386) capture in 12.572 days
				// 5 vessels (v=2.609) capture in 11.497 days
				// 10 vessels (v=3.303) capture in 9.084 days
				// 20 vessels (v=3.996) capture in 7.508 days
				// 50 vessels (v=4.912) capture in 6.107 days
				if p.InfluenceByFaction[faction] > 30.0 {
					p.InfluenceByFaction = [4]float64{}
					p.Faction = faction
					p.AttackDelay = r.scene.Rand().FloatRange(100, 500)
					p.CaptureDelay = r.scene.Rand().FloatRange(400, 600)
					r.world.PushEvent(fmt.Sprintf("%s established control over %s", faction.Name(), p.Info.Name))
				}
			}
		}

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
					generated += 10
				}
				p.MineralDeposit += generated
			}
		}

		if p.Faction == gamedata.FactionNone {
			continue
		}

		r.processPlanetActions(p)

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
