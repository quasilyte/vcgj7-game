package gamedata

func GetSalary(exp int) int {
	return (GetRank(exp) * 3) + 4
}

func GetRank(exp int) int {
	switch {
	case exp == 0:
		return 0
	case exp <= 10:
		return 1
	case exp <= 30:
		return 2
	case exp <= 70:
		return 3
	case exp <= 150:
		return 4
	case exp <= 300:
		return 5
	case exp <= 600:
		return 6
	case exp <= 1200:
		return 7
	case exp <= 2000:
		return 8
	case exp <= 4000:
		return 9
	case exp <= 9000:
		return 10
	case exp <= 20000:
		return 11
	default:
		return 12
	}
}
