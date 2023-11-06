package optional

import "encoding/json"

type Of[T any] struct {
	value T
	valid bool
}

func New[T any](value T, valid bool) Of[T] {
	return Of[T]{
		value: value,
		valid: valid,
	}
}

func NewWithValue[T any](value T) Of[T] {
	return Of[T]{
		value: value,
		valid: true,
	}
}

func (o Of[T]) ValueOrZero() T {
	if o.valid {
		return o.value
	}

	return defaultValue[T]()
}

func (o Of[T]) HasValue() bool {
	return o.valid
}

func defaultValue[T any]() T {
	v := new(T)
	return *v
}

func (o Of[T]) MarshalJSON() ([]byte, error) {
	return json.Marshal(o.ValueOrZero())
}
