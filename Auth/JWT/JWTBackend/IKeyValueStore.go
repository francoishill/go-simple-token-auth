package JWTBackend

type IKeyValueStore interface {
	SetValue(key string, value string, expiration ...interface{}) error
	GetValue(key string) (interface{}, error)
	CloseConnection() error
}
