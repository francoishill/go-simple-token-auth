package JWTBackend

import (
	"github.com/peterbourgon/diskv"
)

type DefaultKeyValueStore struct {
	diskvInstance *diskv.Diskv
}

func (this *DefaultKeyValueStore) SetValue(key string, value string, expiration ...interface{}) error {
	return this.diskvInstance.Write(key, []byte(value))
}

func (this *DefaultKeyValueStore) GetValue(key string) (interface{}, error) {
	val, err := this.diskvInstance.Read(key)
	if err != nil {
		//This is required, otherwise we return an empty byte array
		return nil, err
	}
	return val, err
}

func (this *DefaultKeyValueStore) CloseConnection() error {
	//Do nothing
	return nil
}
