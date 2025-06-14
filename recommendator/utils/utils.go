package utils

type ComparableSet[T comparable] struct {
	engine map[T]bool
}

func NewSet[T comparable]() *ComparableSet[T] {
	return &ComparableSet[T]{
		engine: make(map[T]bool, 0),
	}
}

func (set *ComparableSet[T]) Add(val T) {
	set.engine[val] = true
}

func (set *ComparableSet[T]) AsArr() []T {
	lst := make([]T, 0, len(set.engine))

	for k, _ := range set.engine {
		lst = append(lst, k)
	}

	return lst
}
