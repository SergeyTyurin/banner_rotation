package bannerselector

import (
	"errors"
	"math"
)

var errIncorrectInput = errors.New("incorrect input")

const invalidIndex = -1

func ucb1(sump, p int, q float64) float64 {
	return q + math.Sqrt(2*math.Log(float64(sump))/float64(p))
}

func isCorrectInput(click, display int) bool {
	return click >= 0 && display >= 0 && display >= click
}

func SelectBannerIndex(displays []int, clicks []int) (int, error) {
	if len(displays) != len(clicks) || len(displays) == 0 {
		return invalidIndex, errIncorrectInput
	}

	N := len(displays) // количество баннеров в выборке
	var sumDisplays int
	for i := 0; i < N; i++ {
		sumDisplays += displays[i]
	}

	// Определяем баннер для показа
	var maxUcb float64
	var bIndex int
	for i := 0; i < N; i++ {
		if !isCorrectInput(clicks[i], displays[i]) {
			return invalidIndex, errIncorrectInput
		}
		if displays[i] == 0 {
			return i, nil
		}

		q := float64(clicks[i]) / float64(displays[i])
		ucb := ucb1(sumDisplays, displays[i], q)
		if maxUcb < ucb {
			maxUcb = ucb
			bIndex = i
		}
	}

	return bIndex, nil
}
