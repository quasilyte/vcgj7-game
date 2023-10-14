package gamedata

type Faction int

const (
	FactionNone Faction = iota
	FactionA
	FactionB
	FactionC
	NumFactions
)

func (f Faction) Name() string {
	switch f {
	case FactionA:
		return "Alpha"
	case FactionB:
		return "Beta"
	case FactionC:
		return "Gamma"
	default:
		return "Unknown"
	}
}
