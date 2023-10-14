package battle

import (
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/ge"
	"github.com/quasilyte/vcgj7-game/assets"
)

func playSound(scene *ge.Scene, id resource.AudioID) {
	numSamples := assets.NumSamples(id)
	if numSamples == 1 {
		scene.Audio().PlaySound(id)
	} else {
		soundIndex := scene.Rand().IntRange(0, numSamples-1)
		sound := resource.AudioID(int(id) + soundIndex)
		scene.Audio().PlaySound(sound)
	}
}
