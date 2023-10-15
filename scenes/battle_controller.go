package scenes

import (
	"github.com/quasilyte/ge"
	"github.com/quasilyte/vcgj7-game/assets"
	"github.com/quasilyte/vcgj7-game/battle"
	"github.com/quasilyte/vcgj7-game/gamedata"
	"github.com/quasilyte/vcgj7-game/session"
)

type BattleController struct {
	state *session.State

	challenge int
	enemy     *gamedata.VesselDesign
	runner    *battle.Runner
}

func NewBattleController(state *session.State, enemy *gamedata.VesselDesign) *BattleController {
	return &BattleController{
		state:     state,
		challenge: enemy.Challenge,
		enemy:     enemy,
	}
}

func (c *BattleController) Init(scene *ge.Scene) {
	scene.Audio().PauseCurrentMusic()
	scene.Audio().PlayMusic(assets.AudioMusicCombat)

	c.runner = battle.NewRunner(battle.RunnerConfig{
		Input:  c.state.Input,
		Enemy:  c.enemy,
		Player: c.state.World.Player,
	})
	scene.AddObject(c.runner)

	c.runner.EventBattleOver.Connect(nil, func(results battle.Results) {
		scene.DelayedCall(2, func() {
			player := c.state.World.Player

			var minExp int
			var maxExp int
			var minCredits int
			var maxCredits int
			var minCargo int
			var maxCargo int
			creditsChance := 0.0
			cargoChance := 0.0
			switch c.challenge {
			case 0:
				minExp = 5
				maxExp = 10
				cargoChance = 0.3
				minCargo = 2
				maxCargo = 6
			case 1:
				minExp = 15
				maxExp = 25
				creditsChance = 0.2
				minCredits = 5
				maxCredits = 10
				cargoChance = 0.4
				minCargo = 2
				maxCargo = 10
			case 2:
				minExp = 40
				maxExp = 60
				creditsChance = 0.4
				minCredits = 15
				maxCredits = 30
				cargoChance = 0.5
				minCargo = 2
				maxCargo = 16
			case 3:
				minExp = 100
				maxExp = 150
				creditsChance = 0.7
				minCredits = 25
				maxCredits = 50
				cargoChance = 0.6
				minCargo = 3
				maxCargo = 25
			}
			player.BattleRewards = gamedata.BattleRewards{
				Victory:    results.Victory,
				Experience: scene.Rand().IntRange(minExp, maxExp),
			}
			if creditsChance > 0 && scene.Rand().Chance(creditsChance) {
				player.BattleRewards.Credits = scene.Rand().IntRange(minCredits, maxCredits)
			}
			if cargoChance > 0 && scene.Rand().Chance(cargoChance) {
				player.BattleRewards.Cargo = scene.Rand().IntRange(minCargo, maxCargo)
			}
			if c.enemy.Elite {
				player.BattleRewards.Experience *= 2
			}

			if player.BattleRewards.Cargo == 0 && player.BattleRewards.Credits == 0 {
				if player.Fuel < 70 && scene.Rand().Chance(0.6) {
					player.BattleRewards.Fuel = scene.Rand().IntRange(2, 10)
				}
			}

			player.BattleRewards.SystemLiberated = c.enemy.LastDefender

			player.VesselHP = results.HP
			player.Mode = gamedata.ModeAfterCombat
			player.Battles++
			scene.Context().ChangeScene(NewChoiceController(c.state))
		})
	})
}

func (c *BattleController) Update(delta float64) {}
