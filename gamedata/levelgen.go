package gamedata

func NewWorld() *World {
	w := &World{}

	planets := make([]*Planet, len(Planets))
	for i := range planets {
		p := &Planet{
			Info: Planets[i],
		}
		planets[i] = p
	}

	w.Planets = planets

	w.Player = &Player{
		Planet: planets[0],
	}

	return w
}
