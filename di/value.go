package di

func AddValue[T any](value T) {
	AddSingleton[T](func() T {
		return value
	})
}
