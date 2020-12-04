package datasource

//go:generate mockgen -source=datasource.go -destination=../../gen/mocks/mock_datasource/datasource.go -self_package=../pkg/datasource

type Datasource interface {
	GetKey(key string) (string, error)
	SetKey(key string, val string, timeoutSeconds int) error
	DelKey(key string) error
}

type KeyNotFoundError string

func (e KeyNotFoundError) Error() string {
	return string(e)
}

const KeyNotFound = KeyNotFoundError("redis: nil")
