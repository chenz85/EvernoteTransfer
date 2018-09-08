package oauth

import (
	"net/http"

	"github.com/czsilence/EvernoteTransfer/erro"
)

/////////////////////////////////////////////////////////
// oauth 配置

type OAuthConfig interface {
	OAuth_Auth() (auth_url string, request_secret string, err erro.Error)
	OAuth_ParseCallback(req *http.Request) (tok, verifier string, err erro.Error)
	OAuth_AccessToken(tok, verifier, request_secret string) (response_values map[string]string, err erro.Error)
}
