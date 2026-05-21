package ptr

func FromPtr[T any](ptr *T) T {
	if ptr == nil {
		var zero T
		return zero
	}
	return *ptr
}

func ToPtr[T any](value T) *T {
	return &value
}

func ToPtrNonEmpty[T comparable](v T) *T {
	var zero T
	if v == zero {
		return nil
	}
	return &v
}
