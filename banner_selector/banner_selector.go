package banner_selector

import (
	"errors"
	"math"
)

var errIncorrectInput = errors.New("incorrect input")

const invalidIndex = -1

func ucb1(sum_p, p int, q float64) float64 {
	return q + math.Sqrt(2*math.Log(float64(sum_p))/float64(p))
}

func isCorrectInput(click, display int) bool {
	return click >= 0 && display >= 0 && display >= click
}

func SelectBannerIndex(displays []int, clicks []int) (int, error) {
	if len(displays) != len(clicks) || len(displays) == 0 {
		return invalidIndex, errIncorrectInput
	}

	N := int(len(displays)) // количество баннеров в выборке
	var sum_displays int
	for i := 0; i < N; i++ {
		sum_displays += displays[i]
	}

	// Определяем баннер для показа
	var max_ucb float64
	var b_index int
	for i := 0; i < N; i++ {
		if !isCorrectInput(clicks[i], displays[i]) {
			return invalidIndex, errIncorrectInput
		}
		if displays[i] == 0 {
			return i, nil
		}

		q := float64(clicks[i]) / float64(displays[i])
		ucb := ucb1(sum_displays, displays[i], q)
		if max_ucb < ucb {
			max_ucb = ucb
			b_index = i
		}
	}

	return b_index, nil
}
