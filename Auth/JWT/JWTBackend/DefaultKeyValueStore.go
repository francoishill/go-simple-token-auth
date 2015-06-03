package JWTBackend

import (
	"github.com/garyburd/redigo/redis"
)

type DefaultKeyValueStore struct {
	conn redis.Conn
}

func (this *DefaultKeyValueStore) SetValue(key string, value string, expiration ...interface{}) error {
	_, err := this.conn.Do("SET", key, value)

	if err == nil && expiration != nil {
		this.conn.Do("EXPIRE", key, expiration[0])
	}

	return err
}

func (this *DefaultKeyValueStore) GetValue(key string) (interface{}, error) {
	return this.conn.Do("GET", key)
}

func (this *DefaultKeyValueStore) CloseConnection() error {
	return this.conn.Close()
}
