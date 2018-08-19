package storage

import (
	"github.com/czsilence/EvernoteTransfer/erro"
	"github.com/czsilence/go/log"
	"github.com/czsilence/go/typo"
	"github.com/gogo/protobuf/proto"
)

func Get(db DB, key string, val typo.Any) bool {
	if data, ex := db.Get(key); !ex {
		return ex
	} else if proto_val, ok := val.(proto.Message); ok {
		if err := proto.Unmarshal(data, proto_val); err != nil {
			log.W2("[storage] parse from db failed, type: %T, err: %v", val, err)
			return false
		}
		return true
	} else {
		log.E2("[storage] invalid type to get from db, type: %T", val)
		return false
	}
}

func Put(db DB, key string, val proto.Message) erro.Error {
	if data, err := proto.Marshal(val); err != nil {
		return erro.E_Storage_DataMarshalFailed.F("err: %v", err)
	} else {
		return db.Put(key, data)
	}
}
