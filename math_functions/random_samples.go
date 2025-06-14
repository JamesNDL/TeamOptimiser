package math_functions

import (
	"errors"
	"math/rand"
)

func RandomSampleWithoutReplacement[T any](list *[]T, numberOfSamples int) ([]T, error) {
	if numberOfSamples <= 0 {
		return []T{}, errors.New("number of samples in RandomSampleWithoutReplacement needs to be > 0.")
	}

	myList := *list
	rand.Shuffle(len(myList), func(i, j int) {
		myList[i], myList[j] = myList[j], myList[i]
	})

	if numberOfSamples > len(myList) {
		numberOfSamples = len(myList)
	}

	samples := myList[:numberOfSamples]
	*list = myList[numberOfSamples:]
	return samples, nil
}
