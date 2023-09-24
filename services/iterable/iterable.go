package iterable

type Iterable[T any] interface {
	Next() bool
	Value() T
}
