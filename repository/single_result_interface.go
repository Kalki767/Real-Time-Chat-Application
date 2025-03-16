package repository

type SingleResultInterface interface {
	Decode(v interface{}) error
}
