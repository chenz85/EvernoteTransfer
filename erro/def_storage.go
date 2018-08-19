package erro

import "github.com/czsilence/go/erro"

var (
	E_Storage_DataMarshalFailed = erro.New(400001, "marshal data from db failed")
	E_Storage_PutFailed         = erro.New(400002, "put data to db failed")
	E_Storage_DelFailed         = erro.New(400003, "del dta from db failed")
)
