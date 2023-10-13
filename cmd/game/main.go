package main

import (
	"time"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/vcgj7-game/scenes"
)

func main() {
	ctx := ge.NewContext(ge.ContextConfig{})
	ctx.Rand.SetSeed(time.Now().Unix())
	ctx.GameName = "planet_eaters"
	ctx.WindowTitle = "Planet Eaters"
	ctx.WindowWidth = 1920 / 2
	ctx.WindowHeight = 1080 / 2
	ctx.FullScreen = true

	if err := ge.RunGame(ctx, scenes.NewMainMenuController()); err != nil {
		panic(err)
	}
}
