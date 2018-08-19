package service

import (
	"github.com/czsilence/EvernoteTransfer/erro"
	"github.com/czsilence/EvernoteTransfer/web"
	"github.com/czsilence/go/log"
	"github.com/gogo/protobuf/proto"
)

func init() {
	web.RegisterAPIFunc("en/oauth", _api_oauth)
	web.RegisterAPIFunc("en/oauth/callback", _api_oauth_callback)
}

func _api_oauth(ctx *web.APIContext) (resp_msg proto.Message, err erro.Error) {
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

func _api_oauth_callback(ctx *web.APIContext) (resp_msg proto.Message, err erro.Error) {
	var oauth_token = ctx.Querys["oauth_token"]
	var oauth_verifier = ctx.Querys["oauth_verifier"]
	var sandbox_lnb = ctx.Querys["sandbox_lnb"]
	log.W("TODO:", oauth_token, oauth_verifier, sandbox_lnb)
	return
}
