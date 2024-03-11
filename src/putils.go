package shelf

func GetPtr[T any](x T) *T {
	return &x
}
