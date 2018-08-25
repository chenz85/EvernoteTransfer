package erro

import "github.com/czsilence/go/erro"

var (
	E_API_UnknownAPI            = erro.New(201001, "unknown api")
	E_API_MarshalResponseFailed = erro.New(201002, "marshal response data failed")
	E_API_WriteResponseFailed   = erro.New(201003, "write response failed")
	E_API_Redirect              = erro.New(201004, "already redirect, ignore write response")
)
