package service

import (
	"fmt"

	"github.com/czsilence/EvernoteTransfer/service/oauth"
	"github.com/czsilence/go/typo"
)

var (
	oauth_mapper = make(map[string]oauth.OAuthConfig)
)

func Start(_opt typo.Map) {
	config(_opt)

	oauth_mapper["en:from"] = oauth.NewConfig(opt.Evernote.Key, opt.Evernote.Secret, "http://127.0.0.1:8001/api/en/oauth/callback?svr=en&side=from")
	oauth_mapper["en:to"] = oauth.NewConfig(opt.Evernote.Key, opt.Evernote.Secret, "http://127.0.0.1:8001/api/en/oauth/callback?svr=en&side=to")
	oauth_mapper["yx:from"] = oauth.NewConfig(opt.Yinxiang.Key, opt.Yinxiang.Secret, "http://127.0.0.1:8001/api/en/oauth/callback?svr=yx&side=from")
	oauth_mapper["yx:to"] = oauth.NewConfig(opt.Yinxiang.Key, opt.Yinxiang.Secret, "http://127.0.0.1:8001/api/en/oauth/callback?svr=yx&side=to")
}

func get_oauth_config(svr, side string) (conf oauth.OAuthConfig, exist bool) {
	conf, exist = oauth_mapper[fmt.Sprintf("%s:%s", svr, side)]
	return
}
