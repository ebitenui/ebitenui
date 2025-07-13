package constantutil

func ConstantToPointer[T any](input T) *T {
	return &input
}
