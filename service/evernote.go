package service

import (
	"net/http"

	"github.com/czsilence/EvernoteTransfer/erro"
	"github.com/czsilence/EvernoteTransfer/storage"
	"github.com/czsilence/EvernoteTransfer/storage/badger"
	"github.com/czsilence/EvernoteTransfer/web"
	"github.com/gogo/protobuf/proto"
)

func init() {
	web.RegisterAPIFunc("en/oauth", http.MethodPost, _api_oauth)
	web.RegisterAPIFunc("en/oauth/callback", http.MethodGet, _api_oauth_callback)
}

func _api_oauth(ctx *web.APIContext) (resp_msg proto.Message, err erro.Error) {
	if auth_url, request_secret, ae := OAuth_Auth(); ae != nil {
		err = ae
	} else {
		resp_msg = &ApiRespOauth{
			AuthorizationUrl: auth_url,
		}

		badger.DB().Put("oauth_req_secret:"+ctx.Sid(), []byte(request_secret))
	}
	return
}

func _api_oauth_callback(ctx *web.APIContext) (resp_msg proto.Message, err erro.Error) {
	if tok, verifier, pe := OAuth_ParseCallback(ctx.Req()); pe != nil {
		err = pe
	} else {
		var sid = ctx.Sid()
		if sec, ex := badger.DB().Get("oauth_req_secret:" + sid); !ex {
			err = erro.E_OAuth_NoRequestSecret
		} else {
			var access_token, access_secret string
			if access_token, access_secret, err = OAuth_AccessToken(tok, verifier, string(sec)); err == nil {
				var oauth_info = &OauthInfo{
					Sid:          sid,
					AccessToken:  access_token,
					AccessSecret: access_secret,
				}
				err = storage.Put(badger.DB(), "oauth:"+sid, oauth_info)
			}
		}

	}

	if err == nil {
		ctx.Gin().Redirect(http.StatusMovedPermanently, "/")
		err = erro.E_API_Redirect
	}
	return
}
