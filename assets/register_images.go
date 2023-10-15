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

		ImageSystemMap:     {Path: "image/map.png"},
		ImageMapLocation:   {Path: "image/map_location.png"},
		ImageAlliedPlanet:  {Path: "image/allied_planet_sector.png"},
		ImageHostilePlanet: {Path: "image/hostile_planet_sector.png"},

		ImageMenuBg: {Path: "image/menu_bg.png"},

		ImageBattleHUD:       {Path: "image/battle_hud.png"},
		ImageBattleBarHP:     {Path: "image/hp_bar.png"},
		ImageBattleBarEnergy: {Path: "image/energy_bar.png"},
		ImageBattleBg:        {Path: "image/combat_bg.png"},

		ImageEnergyShield: {Path: "image/energy_shield.png"},

		ImageProjectilePhotonCannon:  {Path: "image/projectile/photon_cannon.png"},
		ImageProjectileIonCannon:     {Path: "image/projectile/ion_cannon.png"},
		ImageProjectilePulseLaser:    {Path: "image/projectile/pulse_laser.png"},
		ImageProjectileAssaultLaser:  {Path: "image/projectile/assault_laser.png"},
		ImageProjectileScatterGun:    {Path: "image/projectile/scatter_gun.png"},
		ImageProjectileTrident:       {Path: "image/projectile/trident.png"},
		ImageProjectileLance:         {Path: "image/projectile/lance.png"},
		ImageProjectileMissile:       {Path: "image/projectile/missile.png"},
		ImageProjectileHomingMissile: {Path: "image/projectile/homing_missile.png"},
		ImageProjectileTorpedo:       {Path: "image/projectile/torpedo.png"},

		ImagePhotonCannonImpact: {Path: "image/effect/photon_cannon_impact.png", FrameWidth: 14},
		ImageIonCannonImpact:    {Path: "image/effect/ion_cannon_impact.png", FrameWidth: 10},
		ImageAssaultLaserImpact: {Path: "image/effect/assault_laser_impact.png", FrameWidth: 14},
		ImageScatterGunImpact:   {Path: "image/effect/scatter_gun_impact.png", FrameWidth: 11},
		ImageTridentImpact:      {Path: "image/effect/trident_impact.png", FrameWidth: 11},
		ImageLanceImpact:        {Path: "image/effect/lance_impact.png", FrameWidth: 32},
		ImageMissileImpact:      {Path: "image/effect/missile_impact.png", FrameWidth: 24},

		ImageBigExplosion: {Path: "image/effect/big_explosion.png", FrameWidth: 32},

		ImageVesselPlayer:     {Path: "image/vessel/player.png", FrameWidth: 48},
		ImageVesselBetaSmall:  {Path: "image/vessel/beta_small.png", FrameWidth: 48},
		ImageVesselBetaBig:    {Path: "image/vessel/beta_big.png", FrameWidth: 48},
		ImageVesselGammaSmall: {Path: "image/vessel/gamma_small.png", FrameWidth: 48},
		ImageVesselGammaBig:   {Path: "image/vessel/gamma_big.png", FrameWidth: 48},
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

	ImageMenuBg

	ImageSystemMap
	ImageMapLocation
	ImageAlliedPlanet
	ImageHostilePlanet

	ImageBattleHUD
	ImageBattleBarHP
	ImageBattleBarEnergy
	ImageBattleBg

	ImageEnergyShield

	ImageProjectilePhotonCannon
	ImageProjectileIonCannon
	ImageProjectilePulseLaser
	ImageProjectileAssaultLaser
	ImageProjectileScatterGun
	ImageProjectileTrident
	ImageProjectileLance
	ImageProjectileMissile
	ImageProjectileHomingMissile
	ImageProjectileTorpedo

	ImagePhotonCannonImpact
	ImageIonCannonImpact
	ImageAssaultLaserImpact
	ImageScatterGunImpact
	ImageTridentImpact
	ImageLanceImpact
	ImageMissileImpact

	ImageBigExplosion

	ImageVesselPlayer
	ImageVesselBetaSmall
	ImageVesselBetaBig
	ImageVesselGammaSmall
	ImageVesselGammaBig
)
