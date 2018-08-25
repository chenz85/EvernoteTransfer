package erro

import "github.com/czsilence/go/erro"

var (
	E_OAUTH_FAILED            = erro.New(300001, "oauth failed")
	E_OAuth_NoRequestSecret   = erro.New(300002, "no request secret found")
	E_OAuth_OAuthInfoNotFound = erro.New(300003, "oauth info not found")
)
