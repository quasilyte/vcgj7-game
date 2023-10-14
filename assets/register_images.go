package assets

import (
	resource "github.com/quasilyte/ebitengine-resource"
	"github.com/quasilyte/ge"

	_ "image/png"
)

func registerImageResources(ctx *ge.Context) {
	imageResources := map[resource.ImageID]resource.ImageInfo{
		ImageUIButtonDisabled: {Path: "image/ebitenui/button-disabled.png"},
		ImageUIButtonIdle:     {Path: "image/ebitenui/button-idle.png"},
		ImageUIButtonHover:    {Path: "image/ebitenui/button-hover.png"},
		ImageUIButtonPressed:  {Path: "image/ebitenui/button-pressed.png"},
		ImageUIPanelIdle:      {Path: "image/ebitenui/panel-idle.png"},

		ImageSystemMap:   {Path: "image/map.png"},
		ImageMapLocation: {Path: "image/map_location.png"},

		ImageBattleHUD:       {Path: "image/battle_hud.png"},
		ImageBattleBarHP:     {Path: "image/hp_bar.png"},
		ImageBattleBarEnergy: {Path: "image/energy_bar.png"},
		ImageBattleBg:        {Path: "image/combat_bg.png"},

		ImageEnergyShield: {Path: "image/energy_shield.png"},

		ImageProjectileIonCannon:     {Path: "image/projectile/ion_cannon.png"},
		ImageProjectilePulseLaser:    {Path: "image/projectile/pulse_laser.png"},
		ImageProjectileMissile:       {Path: "image/projectile/missile.png"},
		ImageProjectileHomingMissile: {Path: "image/projectile/homing_missile.png"},

		ImageIonCannonImpact: {Path: "image/effect/ion_cannon_impact.png", FrameWidth: 10},
		ImageMissileImpact:   {Path: "image/effect/missile_impact.png", FrameWidth: 24},

		ImageVesselRaider:   {Path: "image/vessel/raider.png", FrameWidth: 48},
		ImageVesselMarauder: {Path: "image/vessel/marauder.png", FrameWidth: 48},
	}

	for id, res := range imageResources {
		ctx.Loader.ImageRegistry.Set(id, res)
		ctx.Loader.LoadImage(id)
	}
}

const (
	ImageNone resource.ImageID = iota

	ImageUIButtonDisabled
	ImageUIButtonIdle
	ImageUIButtonHover
	ImageUIButtonPressed
	ImageUIPanelIdle

	ImageSystemMap
	ImageMapLocation

	ImageBattleHUD
	ImageBattleBarHP
	ImageBattleBarEnergy
	ImageBattleBg

	ImageEnergyShield

	ImageProjectileIonCannon
	ImageProjectilePulseLaser
	ImageProjectileMissile
	ImageProjectileHomingMissile

	ImageIonCannonImpact
	ImageMissileImpact

	ImageVesselRaider
	ImageVesselMarauder
)
