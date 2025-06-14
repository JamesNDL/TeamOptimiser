package math_functions

import (
	"errors"
	"math/rand"
)

func RandomSampleWithoutReplacement[T any](list []T, numberOfSamples int) ([]T, error) {
	if numberOfSamples <= 0 {
		return nil, errors.New("number of samples must be > 0")
	}

	shuffled := append([]T(nil), list...) // make a copy
	rand.Shuffle(len(shuffled), func(i, j int) {
		shuffled[i], shuffled[j] = shuffled[j], shuffled[i]
	})

	if numberOfSamples > len(shuffled) {
		numberOfSamples = len(shuffled)
	}

	return shuffled[:numberOfSamples], nil
}
