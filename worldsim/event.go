package worldsim

type eventKind int

const (
	eventUnknown eventKind = iota
	eventFuelScavenge
	eventBuyFuel
)
