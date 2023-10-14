package worldsim

import (
	"fmt"
	"math"
	"strings"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/gmath"
	"github.com/quasilyte/gsignal"
	"github.com/quasilyte/vcgj7-game/gamedata"
)

type Runner struct {
	world       *gamedata.World
	scene       *ge.Scene
	choices     []Choice
	jumpOptions []jumpOption
	textLines   []string

	eventInfo eventInfo

	EventChoiceSelected gsignal.Event[gamedata.Mode]
}

type eventInfo struct {
	kind eventKind
}

type jumpOption struct {
	planet   *gamedata.Planet
	fuelCost int
	time     int
}

type GeneratedChoices struct {
	Text    string
	Choices []Choice
}

func NewRunner(w *gamedata.World) *Runner {
	return &Runner{
		world:       w,
		choices:     make([]Choice, 0, 8),
		jumpOptions: make([]jumpOption, 0, 8),
		textLines:   make([]string, 0, 20),
	}
}

func (r *Runner) Init(scene *ge.Scene) {
	r.scene = scene
}

func (r *Runner) generateEventChoices(event eventInfo) string {
	player := r.world.Player

	switch event.kind {
	case eventFuelScavenge:
		fuelScavenged := r.scene.Rand().IntRange(3, 12)
		r.choices = append(r.choices, Choice{
			Text: "Done",
			OnSelected: func() {
				player.Fuel = gmath.ClampMax(player.Fuel+fuelScavenged, player.MaxFuel)
				r.commitChoice(gamedata.ModeOrbiting)
			},
		})
		if r.scene.Rand().Bool() {
			return fmt.Sprintf("%d fuel units acquired.", fuelScavenged)
		}
		return fmt.Sprintf("Scavenged %d fuel units.", fuelScavenged)

	case eventBuyFuel:
		fuelPrice := 3
		if r.scene.Rand().Chance(0.3) {
			fuelPrice = 2
		}
		maxSpent := 90
		if maxSpent > player.Credits {
			maxSpent = player.Credits
		}
		bought := maxSpent / fuelPrice
		if player.Fuel+bought > player.MaxFuel {
			bought = player.MaxFuel - player.Fuel
		}
		spent := bought * fuelPrice
		r.choices = append(r.choices, Choice{
			Text: "Done",
			OnSelected: func() {
				player.Credits -= spent
				player.Fuel += bought
				r.commitChoice(gamedata.ModeDocked)
			},
		})
		return fmt.Sprintf("Bought %d fuel units for %d credits.", bought, spent)

	default:
		panic(fmt.Sprintf("unexpected event kind: %d", event.kind))
	}
}

func (r *Runner) GenerateChoices() GeneratedChoices {
	r.textLines = r.textLines[:0]
	r.choices = r.choices[:0]
	r.jumpOptions = r.jumpOptions[:0]

	event := r.eventInfo
	r.eventInfo = eventInfo{}
	if event.kind != eventUnknown {
		s := r.generateEventChoices(event)
		return GeneratedChoices{
			Choices: r.choices,
			Text:    s,
		}
	}

	player := r.world.Player
	planet := r.world.Player.Planet

	r.textLines = append(r.textLines, genModeText(r.scene, r.world))

	canJump := true

	if len(r.choices) < MaxChoices && planet.Faction == player.Faction {
		switch player.Mode {
		case gamedata.ModeJustEntered, gamedata.ModeOrbiting:
			s := "Enter the planetary docks"
			if planet.Info.GasGiant {
				s = "Dock the station"
			}
			r.choices = append(r.choices, Choice{
				Time: 1,
				Text: s,
				OnSelected: func() {
					r.commitChoiceExtra(gamedata.ModeOrbiting, gamedata.ModeDocked)
				},
			})
		}
	}

	if len(r.choices) < MaxChoices && player.Mode == gamedata.ModeDocked {
		canJump = false
		if player.Credits > 0 && player.Fuel < player.MaxFuel {
			r.choices = append(r.choices, Choice{
				Time: 2,
				Text: "Buy fuel",
				OnSelected: func() {
					r.eventInfo = eventInfo{kind: eventBuyFuel}
					r.commitChoice(gamedata.ModeDocked)
				},
			})
		}
	}

	hasFuel := false
	if canJump {
		// Find all possible routes first.
		for _, p := range r.world.Planets {
			if p == player.Planet {
				continue
			}
			dist := player.Planet.Info.MapOffset.DistanceTo(p.Info.MapOffset)
			fuelNeeded := gmath.ClampMin(int(dist*player.FuelUsage), 1)
			if player.Fuel < fuelNeeded {
				continue
			}
			hasFuel = true
			if dist > player.MaxJumpDist {
				continue
			}
			hours := int(math.Ceil(dist / player.JumpSpeed))
			r.jumpOptions = append(r.jumpOptions, jumpOption{
				planet:   p,
				fuelCost: fuelNeeded,
				time:     hours,
			})
		}
		gmath.Shuffle(r.scene.Rand(), r.jumpOptions)
		// Add as many travel options as possible.
		for len(r.jumpOptions) > 0 && len(r.choices) < MaxChoices {
			j := r.jumpOptions[len(r.jumpOptions)-1]
			r.jumpOptions = r.jumpOptions[:len(r.jumpOptions)-1]
			r.choices = append(r.choices, Choice{
				Time: j.time,
				Text: fmt.Sprintf("Jump to %s [%d fuel]", j.planet.Info.Name, j.fuelCost),
				OnSelected: func() {
					player.Planet = j.planet
					player.Fuel -= j.fuelCost
					r.commitChoiceExtra(gamedata.ModeJump, gamedata.ModeJustEntered)
				},
			})
		}
	}

	if len(r.choices) < MaxChoices && player.Mode == gamedata.ModeDocked {
		r.choices = append(r.choices, Choice{
			Time: 2,
			Text: "Take off",
			OnSelected: func() {
				r.commitChoice(gamedata.ModeOrbiting)
			},
		})
	}

	if len(r.choices) < MaxChoices {
		switch r.world.Player.Mode {
		case gamedata.ModeJustEntered, gamedata.ModeOrbiting:
			canScavenge := !hasFuel || (player.Fuel < 100 && r.scene.Rand().Chance(0.4))
			if canScavenge {
				r.choices = append(r.choices, Choice{
					Time: 8,
					Text: "Scavenge for fuel",
					OnSelected: func() {
						r.eventInfo = eventInfo{kind: eventFuelScavenge}
						r.commitChoice(gamedata.ModeScavenging)
					},
				})
			} else {
				r.choices = append(r.choices, Choice{
					Time: 5,
					Text: "Wait.",
					OnSelected: func() {
						r.commitChoice(gamedata.ModeOrbiting)
					},
				})
			}
		}
	}

	return GeneratedChoices{
		Text:    strings.Join(r.textLines, "\n"),
		Choices: r.choices,
	}
}

func (r *Runner) commitChoice(m gamedata.Mode) {
	r.world.Player.Mode = m
	r.EventChoiceSelected.Emit(m)
}

func (r *Runner) commitChoiceExtra(m, postMode gamedata.Mode) {
	r.world.Player.Mode = m
	r.EventChoiceSelected.Emit(postMode)
}

func (r *Runner) AdvanceTime(hours int) {
	for i := 0; i < hours; i++ {
		r.world.GameTime++
		// One in-game hour is simulated during 1 second in delta time term.
		for j := 0; j < 5; j++ {
			r.updateWorld(0.2)
		}
	}
}

func (r *Runner) updateWorld(delta float64) {

}
