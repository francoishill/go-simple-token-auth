package JWTBackend

import (
	"github.com/peterbourgon/diskv"
)

type DefaultKeyValueStoreProvider struct {
	storeInstance *DefaultKeyValueStore
}

func ConstructDefaultKeyValueStoreProvider() IKeyValueStoreProvider {
	flatTransform := func(s string) []string { return []string{} }

	return &DefaultKeyValueStoreProvider{
		storeInstance: &DefaultKeyValueStore{
			diskv.New(diskv.Options{
				BasePath:     "tmp-key-vals",
				Transform:    flatTransform,
				CacheSizeMax: 10 * 1024 * 1024,
			}),
		},
	}
}

func (this *DefaultKeyValueStoreProvider) GetNewConnection() IKeyValueStore {
	return this.storeInstance
}
