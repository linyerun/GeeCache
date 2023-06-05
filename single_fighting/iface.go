package single_fighting

type ISingleFighting interface {
	Do(uniqueMark string, fn func() (any, error)) (any, error)
}
