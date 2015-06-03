package JWTBackend

type IKeyValueStoreProvider interface {
	GetNewConnection() IKeyValueStore
}
