package service

import (
	"net/http"

	"github.com/czsilence/EvernoteTransfer/erro"
	"github.com/czsilence/EvernoteTransfer/storage"
	"github.com/czsilence/EvernoteTransfer/storage/badger"
	"github.com/czsilence/EvernoteTransfer/web"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
)

func init() {
	web.RegisterAPIFunc("en/oauth", http.MethodPost, _api_oauth)
	web.RegisterAPIFunc("en/oauth/callback", http.MethodGet, _api_oauth_callback)
	web.RegisterAPIFunc("en/xfer", http.MethodPost, _api_en_transfer)
}

func _api_oauth(ctx *web.APIContext) (resp_msg proto.Message, err erro.Error) {
	var query = new(ApiReqOauth)
	if de := jsonpb.UnmarshalString(string(ctx.Data()), query); de != nil {
		err = erro.E_API_InvalidRequest
		return
	}

	if oc, ex := get_oauth_config(query.Svr, query.Side); !ex {
		err = erro.E_OAuth_OAuthConfigNotFound
	} else if auth_url, request_secret, ae := oc.OAuth_Auth(); ae != nil {
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
	if oc, ex := get_oauth_config(ctx.Gin().Query("svr"), ctx.Gin().Query("side")); !ex {
		err = erro.E_OAuth_OAuthConfigNotFound
	} else if tok, verifier, pe := oc.OAuth_ParseCallback(ctx.Req()); pe != nil {
		err = pe
	} else {
		var sid = ctx.Sid()
		if sec, ex := badger.DB().Get("oauth_req_secret:" + sid); !ex {
			err = erro.E_OAuth_NoRequestSecret
		} else {
			if values, ate := oc.OAuth_AccessToken(tok, verifier, string(sec)); ate != nil {
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

//////////////////////////
func _api_en_transfer(ctx *web.APIContext) (resp_msg proto.Message, err erro.Error) {
	var sid = ctx.Sid()
	var oauth_info = new(OauthInfo)
	if !storage.Get(badger.DB(), "oauth:"+sid, oauth_info) {
		err = erro.E_OAuth_OAuthInfoNotFound
	} else {
		resp_msg = oauth_info
	}
	return
}
