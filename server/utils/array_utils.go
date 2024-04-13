package utils

import (
	"errors"

	"schej.it/server/models"
)

// Inserts the given value at the specified index in the slice. Returns the updated slice
func Insert[T any](arr []T, index int, value T) ([]T, error) {
	if index < 0 {
		return nil, errors.New("index cannot be less than 0")
	}

	if index >= len(arr) {
		return append(arr, value), nil
	}

	arr = append(arr[:index+1], arr[index:]...)
	arr[index] = value

	return arr, nil
}

// Returns whether the given slice contains the given value
func Contains[T comparable](arr []T, value T) bool {
	for _, v := range arr {
		if v == value {
			return true
		}
	}
	return false
}

// Returns the index of the element found using the equals function, -1 if not found
func Find[T any](arr []T, equals func(T) bool) int {
	for i, v := range arr {
		if equals(v) {
			return i
		}
	}
	return -1
}

func ArrayToSet[T comparable](arr []T) models.Set[T] {
	set := make(models.Set[T])

	for _, v := range arr {
		set[v] = struct{}{}
	}

	return set
}

func Map[T any, U any](arr []T, mapFunc func(T) U) []U {
	newArr := make([]U, 0)
	for _, v := range arr {
		newArr = append(newArr, mapFunc(v))
	}
	return newArr
}

type ElementWithIndex[T any] struct {
	Index int
	Value T
}

// Returns an array of elements that were added, removed, and kept based on the arrays that were passed in
func FindAddedRemovedKept[T comparable](arr []T, origArr []T) ([]ElementWithIndex[T], []ElementWithIndex[T], []ElementWithIndex[T]) {
	added := make([]ElementWithIndex[T], 0)
	removed := make([]ElementWithIndex[T], 0)
	kept := make([]ElementWithIndex[T], 0)

	// Find elements that were removed / kept
	for i, origVal := range origArr {
		// Check if original element is present in updated array
		found := false
		for _, val := range arr {
			if val == origVal {
				found = true
				break
			}
		}

		if found {
			// Element was kept
			kept = append(kept, ElementWithIndex[T]{Index: i, Value: origVal})
		} else {
			// Element was removed
			removed = append(removed, ElementWithIndex[T]{Index: i, Value: origVal})
		}
	}

	// Find elements that were added
	for i, val := range arr {
		found := false
		for _, origVal := range kept {
			if val == origVal.Value {
				found = true
				break
			}
		}

		if !found {
			// Element was added
			added = append(added, ElementWithIndex[T]{Index: i, Value: val})
		}
	}

	return added, removed, kept
}
