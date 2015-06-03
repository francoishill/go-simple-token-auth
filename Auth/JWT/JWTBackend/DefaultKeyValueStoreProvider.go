package JWTBackend

import (
	"github.com/garyburd/redigo/redis"
)

type DefaultKeyValueStoreProvider struct {
	redisNetworkHostAndPort string
	redisAuthPassword       string
	instanceKeyValueStore   *DefaultKeyValueStore
}

func ConstructDefaultKeyValueStoreProvider(redisNetworkHostAndPort, redisAuthPassword string) IKeyValueStoreProvider {
	return &DefaultKeyValueStoreProvider{
		redisAuthPassword:       redisAuthPassword,
		redisNetworkHostAndPort: redisNetworkHostAndPort,
		instanceKeyValueStore:   nil,
	}
}

func (this *DefaultKeyValueStoreProvider) GetNewConnection() IKeyValueStore {
	if this.instanceKeyValueStore == nil {
		this.instanceKeyValueStore = new(DefaultKeyValueStore)
		var err error

		this.instanceKeyValueStore.conn, err = redis.Dial("tcp", this.redisNetworkHostAndPort)

		if err != nil {
			panic(err)
		}

		if _, err := this.instanceKeyValueStore.conn.Do("AUTH", this.redisAuthPassword); err != nil {
			this.instanceKeyValueStore.conn.Close()
			panic(err)
		}
	}

	return this.instanceKeyValueStore
}
