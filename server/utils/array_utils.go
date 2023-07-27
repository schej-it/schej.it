package utils

import "errors"

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
