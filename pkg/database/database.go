package database

type IDatabase interface {
	Connect() error
	Close() error
}
