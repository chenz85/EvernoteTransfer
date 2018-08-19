package service

import (
	"github.com/czsilence/EvernoteTransfer/erro"
	"github.com/czsilence/EvernoteTransfer/web"
	"github.com/czsilence/go/typo"
	"github.com/gogo/protobuf/proto"
)

func init() {
	web.RegisterAPIProvider("en/oauth", new(_api_en_oauth))
}

type _api_en_oauth struct {
}

func (api *_api_en_oauth) Process(msg proto.Message, header typo.Any) (resp_msg proto.Message, err erro.Error) {
	var oauth_ctx = &OAuthContext{
		Key:    opt.Evernote.Key,
		Secret: opt.Evernote.Secret,
	}
	if auth_url, ae := oauth_ctx.Auth(); ae != nil {
		err = ae
	} else {
		resp_msg = &ApiRespOauth{
			AuthorizationUrl: auth_url,
		}
	}
	return
}
