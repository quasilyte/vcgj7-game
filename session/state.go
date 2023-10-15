package session

import (
	"github.com/quasilyte/ge/input"
	"github.com/quasilyte/vcgj7-game/eui"
	"github.com/quasilyte/vcgj7-game/gamedata"
)

type State struct {
	UIResources *eui.Resources

	Settings Settings

	Input *input.Handler

	World *gamedata.World
}

type Settings struct {
	SoundLevel int
	MusicLevel int
	Difficulty int
}
