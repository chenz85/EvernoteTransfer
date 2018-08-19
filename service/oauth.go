package service

import (
	"net/http"

	"github.com/czsilence/EvernoteTransfer/erro"
	"github.com/dghubble/oauth1"
)

// oauth for evernote
// ref: https://dev.evernote.com/doc/articles/authentication.php

func OAuth_Auth(consumer_key, consumer_secret string) (auth_url string, err erro.Error) {
	config := oauth1.Config{
		ConsumerKey:    consumer_key,
		ConsumerSecret: consumer_secret,
		CallbackURL:    "http://127.0.0.1:8001/api/en/oauth/callback",
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

func OAuth_ParseAccessToken(req *http.Request) (tok, verifier string, err erro.Error) {
	if _tok, _verifier, pe := oauth1.ParseAuthorizationCallback(req); pe != nil {
		err = erro.E_OAUTH_FAILED.F("err: %v", pe)
	} else {
		tok, verifier = _tok, _verifier
	}
	return
}
