package worldsim

type Choice struct {
	Time       int // In hours
	Text       string
	OnSelected func()
}

const MaxChoices = 6
