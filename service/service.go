package service

import (
	"github.com/czsilence/EvernoteTransfer/service/oauth"
	"github.com/czsilence/go/typo"
)

var (
	oauth_en, oauth_yx oauth.OAuthConfig
)

func Start(_opt typo.Map) {
	config(_opt)
	oauth_en = oauth.NewConfig(opt.Evernote.Key, opt.Evernote.Secret, "http://127.0.0.1:8001/api/en/oauth/callback?t=en")
	oauth_yx = oauth.NewConfig(opt.Evernote.Key, opt.Evernote.Secret, "http://127.0.0.1:8001/api/en/oauth/callback?t=yx")
}
