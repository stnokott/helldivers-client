package transform

func ptr[T any](x T) *T {
	return &x
}
