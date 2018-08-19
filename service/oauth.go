package service

import (
	"github.com/czsilence/EvernoteTransfer/erro"
	"github.com/dghubble/oauth1"
)

// oauth for evernote
// ref: https://dev.evernote.com/doc/articles/authentication.php

type OAuthContext struct {
	Key    string
	Secret string
}

func (oauth *OAuthContext) Auth() (auth_url string, err erro.Error) {
	config := oauth1.Config{
		ConsumerKey:    oauth.Key,
		ConsumerSecret: oauth.Secret,
		CallbackURL:    "http://127.0.0.1:8001/api/oauth/callback",
		Endpoint: oauth1.Endpoint{
			RequestTokenURL: "https://sandbox.evernote.com/oauth",
			AuthorizeURL:    "https://sandbox.evernote.com/OAuth.action",
			AccessTokenURL:  "https://sandbox.evernote.com/oauth",
		},
	}

	if tok, _, re := config.RequestToken(); re != nil {
		err = erro.E_OAUTH_FAILED.F("err: %v", re)
	} else if authorizationURL, ae := config.AuthorizationURL(tok); ae != nil {
		err = erro.E_OAUTH_FAILED.F("err: %v", ae)
	} else {
		auth_url = authorizationURL.String()
	}
	return
}
