package utils

import "os"

// Miscellaneous utility functions //

// Returns a pointer to the value.
//
// Can be useful in situations in which we deal
// with literal values and we don't really want
// create a separate variable for pointer value e.g:
//
//	var pVal *float64 = Ptr(3.14)
func Ptr[T any](val T) *T {
	return &val
}

// Attempts to unwrap the most common return tuple:
// value and error.
// If provided error value is not nil, then this function panics
//
// Can be used in situations when we don't really care
// about the error and instead only want the value
func Unwrap[T any](val T, err error) T {
	if err != nil {
		panic(err)
	}
	return val
}

// Gets value of the environment variable env and if value does not exist,
// fallback to specified value
func Getenv(env string, fallback string) string {
	val := os.Getenv(env)
	if val != "" {
		return val
	}
	return fallback
}
