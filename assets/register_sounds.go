package assets

import (
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/ge"
)

func registerSoundResources(ctx *ge.Context) {
	soundResources := map[resource.AudioID]resource.AudioInfo{
		AudioIonCannon1:      {Path: "audio/ion_cannon1.wav"},
		AudioIonCannonImpact: {Path: "audio/ion_cannon_impact.wav"},

		AudioPulseLaser1: {Path: "audio/pulse_laser1.wav"},
		AudioPulseLaser2: {Path: "audio/pulse_laser2.wav"},
		AudioPulseLaser3: {Path: "audio/pulse_laser3.wav"},
	}

	for id, res := range soundResources {
		ctx.Loader.AudioRegistry.Set(id, res)
		ctx.Loader.LoadAudio(id)
	}
}

func NumSamples(a resource.AudioID) int {
	switch a {
	case AudioPulseLaser1:
		return 3
	default:
		return 1
	}
}

const (
	AudioNone resource.AudioID = iota

	AudioIonCannon1
	AudioIonCannonImpact

	AudioPulseLaser1
	AudioPulseLaser2
	AudioPulseLaser3
)
