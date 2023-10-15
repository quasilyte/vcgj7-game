package assets

import (
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/ge"
)

const (
	SoundGroupEffect uint = iota
	SoundGroupMusic
)

func registerSoundResources(ctx *ge.Context) {
	soundResources := map[resource.AudioID]resource.AudioInfo{
		AudioMusicGlobal: {Path: "audio/music/global.ogg", Group: SoundGroupMusic, Volume: -0.2},
		AudioMusicCombat: {Path: "audio/music/combat.ogg", Group: SoundGroupMusic, Volume: -0.2},

		AudioIonCannon1:      {Path: "audio/ion_cannon1.wav"},
		AudioIonCannonImpact: {Path: "audio/ion_cannon_impact.wav"},

		AudioPhotonCannon1: {Path: "audio/photon_cannon1.wav"},
		AudioPhotonCannon2: {Path: "audio/photon_cannon2.wav"},
		AudioPhotonCannon3: {Path: "audio/photon_cannon3.wav"},

		AudioPulseLaser1: {Path: "audio/pulse_laser1.wav"},
		AudioPulseLaser2: {Path: "audio/pulse_laser2.wav"},
		AudioPulseLaser3: {Path: "audio/pulse_laser3.wav"},

		AudioAssaultLaser1: {Path: "audio/assault_laser1.wav"},
		AudioAssaultLaser2: {Path: "audio/assault_laser2.wav"},

		AudioTrident1: {Path: "audio/trident1.wav"},

		AudioScatterGun1: {Path: "audio/scatter_gun1.wav"},

		AudioLance1: {Path: "audio/lance1.wav"},

		AudioMissile1: {Path: "audio/missile1.wav"},
		AudioMissile2: {Path: "audio/missile2.wav"},
		AudioMissile3: {Path: "audio/missile3.wav"},

		AudioExplosion1: {Path: "audio/explosion1.wav"},
		AudioExplosion2: {Path: "audio/explosion2.wav"},
		AudioExplosion3: {Path: "audio/explosion3.wav"},

		AudioBigExplosion1: {Path: "audio/big_explosion1.wav"},
		AudioBigExplosion2: {Path: "audio/big_explosion2.wav"},

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
	case AudioPhotonCannon1:
		return 3
	case AudioMissile1:
		return 3
	case AudioAssaultLaser1:
		return 2
	case AudioExplosion1:
		return 3
	case AudioBigExplosion1:
		return 2
	default:
		return 1
	}
}

const (
	AudioNone resource.AudioID = iota

	AudioMusicGlobal
	AudioMusicCombat

	AudioIonCannon1
	AudioIonCannonImpact

	AudioScatterGun1

	AudioTrident1

	AudioPhotonCannon1
	AudioPhotonCannon2
	AudioPhotonCannon3

	AudioPulseLaser1
	AudioPulseLaser2
	AudioPulseLaser3

	AudioAssaultLaser1
	AudioAssaultLaser2

	AudioLance1

	AudioMissile1
	AudioMissile2
	AudioMissile3

	AudioExplosion1
	AudioExplosion2
	AudioExplosion3

	AudioBigExplosion1
	AudioBigExplosion2

	AudioShieldAbsorb
)
