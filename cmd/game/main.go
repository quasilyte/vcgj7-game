package main

import (
	"time"

	"github.com/quasilyte/ge"
	"github.com/quasilyte/ge/input"
	"github.com/quasilyte/vcgj7-game/assets"
	"github.com/quasilyte/vcgj7-game/eui"
	"github.com/quasilyte/vcgj7-game/scenes"
	"github.com/quasilyte/vcgj7-game/session"
)

func main() {
	ctx := ge.NewContext(ge.ContextConfig{})
	ctx.Rand.SetSeed(time.Now().Unix())
	ctx.GameName = "planet_eaters"
	ctx.WindowTitle = "Planet Eaters"
	ctx.WindowWidth = 1920 / 2
	ctx.WindowHeight = 1080 / 2
	ctx.FullScreen = true

	ctx.Loader.OpenAssetFunc = assets.MakeOpenAssetFunc(ctx)
	assets.RegisterResources(ctx)

	state := &session.State{
		UIResources: eui.PrepareResources(ctx.Loader),
	}

	keymap := input.Keymap{}
	state.Input = ctx.Input.NewHandler(0, keymap)

	if err := ge.RunGame(ctx, scenes.NewMainMenuController(state)); err != nil {
		panic(err)
	}
}