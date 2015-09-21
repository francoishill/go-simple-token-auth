package JWTBackend

import (
	"github.com/garyburd/redigo/redis"
)

type RedisKeyValueStore struct {
	conn redis.Conn
}

func (this *RedisKeyValueStore) SetValue(key string, value string, expiration ...interface{}) error {
	_, err := this.conn.Do("SET", key, value)

	if err == nil && expiration != nil && len(expiration) > 0 {
		this.conn.Do("EXPIRE", key, expiration[0])
	}

	return err
}

func (this *RedisKeyValueStore) GetValue(key string) (interface{}, error) {
	return this.conn.Do("GET", key)
}

func (this *RedisKeyValueStore) CloseConnection() error {
	return this.conn.Close()
}
