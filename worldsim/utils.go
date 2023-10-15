package worldsim

import "github.com/quasilyte/gmath"

func randIterate[T any](rand *gmath.Rand, slice []T, f func(x T) bool) T {
	var result T
	if len(slice) == 0 {
		return result
	}
	if len(slice) == 1 {
		// Don't use rand() if there is only 1 element.
		x := slice[0]
		if f(x) {
			result = x
		}
		return result
	}

	var slider gmath.Slider
	slider.SetBounds(0, len(slice)-1)
	slider.TrySetValue(rand.IntRange(0, len(slice)-1))
	inc := rand.Bool()
	for i := 0; i < len(slice); i++ {
		x := slice[slider.Value()]
		if inc {
			slider.Inc()
		} else {
			slider.Dec()
		}
		if f(x) {
			result = x
			break
		}
	}
	return result
}
