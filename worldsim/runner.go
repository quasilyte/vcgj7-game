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

	EventChoiceSelected gsignal.Event[gsignal.Void]
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

	r.textLines = append(r.textLines, genModeText(r.scene, r.world))

	canJump := true
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
					r.commitChoice(gamedata.ModeJustEntered)
				},
			})
		}
	}

	if len(r.choices) < MaxChoices {
		canScavenge := !hasFuel || (player.Fuel < 100 && r.scene.Rand().Chance(0.4))
		if canScavenge {
			r.choices = append(r.choices, Choice{
				Time: 8,
				Text: "Scavenge for fuel",
				OnSelected: func() {
					r.eventInfo = eventInfo{kind: eventFuelScavenge}
					// fuelFound := r.scene.Rand().IntRange(3, 15)
					// player.Fuel = gmath.ClampMax(player.Fuel+fuelFound, player.MaxFuel)
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

	return GeneratedChoices{
		Text:    strings.Join(r.textLines, "\n"),
		Choices: r.choices,
	}
}

func (r *Runner) commitChoice(m gamedata.Mode) {
	r.world.Player.Mode = m
	r.EventChoiceSelected.Emit(gsignal.Void{})
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
