package service

import (
	"net/http"

	"github.com/czsilence/EvernoteTransfer/erro"
	"github.com/czsilence/oauth1"
)

// oauth for evernote
// ref: https://dev.evernote.com/doc/articles/authentication.php

var oauth_config *oauth1.Config

func oauth_init() {
	oauth_config = &oauth1.Config{
		ConsumerKey:    opt.Evernote.Key,
		ConsumerSecret: opt.Evernote.Secret,
		CallbackURL:    "http://127.0.0.1:8001/api/en/oauth/callback",
		Endpoint: oauth1.Endpoint{
			RequestTokenURL: "https://sandbox.evernote.com/oauth",
			AuthorizeURL:    "https://sandbox.evernote.com/OAuth.action",
			AccessTokenURL:  "https://sandbox.evernote.com/oauth",
		},
	}
}

func OAuth_Auth() (auth_url string, request_secret string, err erro.Error) {
	if tok, sec, re := oauth_config.RequestToken(); re != nil {
		err = erro.E_OAUTH_FAILED.F("err: %v", re)
	} else if authorizationURL, ae := oauth_config.AuthorizationURL(tok); ae != nil {
		err = erro.E_OAUTH_FAILED.F("err: %v", ae)
	} else {
		auth_url = authorizationURL.String()
		request_secret = sec
	}
	return
}

func OAuth_ParseCallback(req *http.Request) (tok, verifier string, err erro.Error) {
	if _tok, _verifier, pe := oauth1.ParseAuthorizationCallback(req); pe != nil {
		err = erro.E_OAUTH_FAILED.F("err: %v", pe)
	} else {
		tok, verifier = _tok, _verifier
	}
	return
}

func OAuth_AccessToken(tok, verifier, request_secret string) (access_token, access_secret string, err erro.Error) {
	at, as, ae := oauth_config.AccessToken(tok, request_secret, verifier)
	if ae != nil {
		err = erro.E_OAUTH_FAILED.F("err: %v", ae)
	} else {
		access_token, access_secret = at, as
	}
	return
}
