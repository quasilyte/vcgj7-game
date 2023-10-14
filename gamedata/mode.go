package gamedata

type Mode int

const (
	ModeUnknown Mode = iota
	ModeOrbiting
	ModeScavenging
	ModeJustEntered
)
