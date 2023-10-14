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

		AudioMissile1: {Path: "audio/missile1.wav"},
		AudioMissile2: {Path: "audio/missile2.wav"},
		AudioMissile3: {Path: "audio/missile3.wav"},

		AudioExplosion1: {Path: "audio/explosion1.wav"},
		AudioExplosion2: {Path: "audio/explosion2.wav"},
		AudioExplosion3: {Path: "audio/explosion3.wav"},
		AudioExplosion4: {Path: "audio/explosion4.wav"},
		AudioExplosion5: {Path: "audio/explosion5.wav"},

		AudioShieldAbsorb: {Path: "audio/shield_absorb.wav"},
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
	case AudioMissile1:
		return 3
	case AudioExplosion1:
		return 5
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

	AudioMissile1
	AudioMissile2
	AudioMissile3

	AudioExplosion1
	AudioExplosion2
	AudioExplosion3
	AudioExplosion4
	AudioExplosion5

	AudioShieldAbsorb
)
