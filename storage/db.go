package storage

import "github.com/czsilence/EvernoteTransfer/erro"

type DB interface {
	Close()

	Put(key string, data []byte) erro.Error
	Get(key string) (data []byte, exist bool)
	Del(key string) erro.Error

	// TODO: 事务接口
}
