package battle

const (
	collisionPlayer1 = 0b0001
	collisionPlayer2 = 0b0010
)

func enemyCollisionMask(m uint16) uint16 {
	if m == collisionPlayer1 {
		return collisionPlayer2
	}
	return collisionPlayer1
}
