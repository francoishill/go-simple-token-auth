package JWTBackend

import (
	"github.com/garyburd/redigo/redis"
)

type RedisKeyValueStoreProvider struct {
	redisNetworkHostAndPort string
	redisAuthPassword       string
	instanceKeyValueStore   *RedisKeyValueStore
}

func ConstructRedisKeyValueStoreProvider(redisNetworkHostAndPort, redisAuthPassword string) IKeyValueStoreProvider {
	return &RedisKeyValueStoreProvider{
		redisAuthPassword:       redisAuthPassword,
		redisNetworkHostAndPort: redisNetworkHostAndPort,
		instanceKeyValueStore:   nil,
	}
}

func (this *RedisKeyValueStoreProvider) GetNewConnection() IKeyValueStore {
	if this.instanceKeyValueStore == nil {
		this.instanceKeyValueStore = new(RedisKeyValueStore)
		var err error

		this.instanceKeyValueStore.conn, err = redis.Dial("tcp", this.redisNetworkHostAndPort)

		if err != nil {
			panic(err)
		}

		if this.redisAuthPassword != "" {
			if _, err := this.instanceKeyValueStore.conn.Do("AUTH", this.redisAuthPassword); err != nil {
				this.instanceKeyValueStore.conn.Close()
				panic(err)
			}
		}
	}

	return this.instanceKeyValueStore
}
