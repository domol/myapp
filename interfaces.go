package main

type Repo[T any] interface {
	list() ([]T, error)
	get(int64) (T, error)
	create(T) (T, error)
	delete(int64) error
	update(int64, T) error
}
