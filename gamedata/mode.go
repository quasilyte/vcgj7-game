package gamedata

//go:generate stringer -type=Mode -trimprefix=Mode
type Mode int

const (
	ModeUnknown Mode = iota
	ModeJump
	ModeOrbiting
	ModeScavenging
	ModeJustEntered
	ModeDocked
)
