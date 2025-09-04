package bg_test

func unwrap[T any](val T, err error) T {
	return val
}
