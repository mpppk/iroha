package lib

type Storage interface {
	Set(curs []int) error
	Get(curs []int) ([][]*Word, bool, error)
}

type NoStorage struct{}

func (e *NoStorage) Get(indices []int) ([][]*Word, bool, error) {
	return nil, false, nil
}

func (e *NoStorage) Set(indices []int) error {
	return nil
}
