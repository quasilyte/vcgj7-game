package session

import (
	"github.com/quasilyte/ge/input"
	"github.com/quasilyte/vcgj7-game/eui"
)

type State struct {
	UIResources *eui.Resources

	Settings Settings

	Input *input.Handler
}

type Settings struct {
	SoundLevel int
	MusicLevel int
	Difficulty int
	Tutorial   int
}
