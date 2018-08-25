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
			if values, ate := OAuth_AccessToken(tok, verifier, string(sec)); ate != nil {
				err = ate
			} else {
				var oauth_info = &OauthInfo{
					Sid: sid,
				}

				oauth_info.AccessToken, _ = values["oauth_token"]
				oauth_info.AccessSecret, _ = values["oauth_token_secret"]
				oauth_info.EdamUserId, _ = values["edam_userId"]
				oauth_info.EdamShard, _ = values["edam_shard"]
				oauth_info.EdamNoteStoreUrl, _ = values["edam_noteStoreUrl"]
				oauth_info.EdamExpires, _ = values["edam_expires"]
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
