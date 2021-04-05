package database

type DisconectDB func()

type Database interface {
	Disconnect() error
	Ping() error
}

//go:generate mockgen -destination mongo_mock.go -package database . MongoDB
type MongoDB interface {
	Database
	Create(namespace string, data interface{}) (string, error)
	Find(namespace, idStr string, dst interface{}) error
	Query(namespace string, filters, sort interface{}, offset, limit int, dst interface{}) error
	Total(namespace string, filters interface{}) (int, error)
	Update(namespace, idStr string, data interface{}) error
	Delete(namespace, idStr string) error
}

//go:generate mockgen -destination redis_mock.go -package database . RedisDB
type RedisDB interface {
	Database
	Set(idStr string, data interface{}) (string, error)
	Get(idStr string, dst interface{}) error
	Delete(idStr string) error
}
