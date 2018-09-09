package oauth

import (
	"net/http"

	"github.com/czsilence/EvernoteTransfer/erro"
	"github.com/czsilence/oauth1"
)

// oauth for evernote
// ref: https://dev.evernote.com/doc/articles/authentication.php

type _OAuthConfigInternal struct {
	c *oauth1.Config
}

func NewConfig(key, secret, host, callback string) OAuthConfig {
	var oauth_config = &oauth1.Config{
		ConsumerKey:    key,
		ConsumerSecret: secret,
		CallbackURL:    callback,
		Endpoint: oauth1.Endpoint{
			RequestTokenURL: host + "oauth",
			AuthorizeURL:    host + "OAuth.action",
			AccessTokenURL:  host + "oauth",
		},
	}

	return &_OAuthConfigInternal{
		c: oauth_config,
	}
}

func (o *_OAuthConfigInternal) OAuth_Auth() (auth_url string, request_secret string, err erro.Error) {
	if tok, sec, re := o.c.RequestToken(); re != nil {
		err = erro.E_OAUTH_FAILED.With(re)
	} else if authorizationURL, ae := o.c.AuthorizationURL(tok); ae != nil {
		err = erro.E_OAUTH_FAILED.With(ae)
	} else {
		auth_url = authorizationURL.String()
		request_secret = sec
	}
	return
}

func (o *_OAuthConfigInternal) OAuth_ParseCallback(req *http.Request) (tok, verifier string, err erro.Error) {
	if _tok, _verifier, pe := oauth1.ParseAuthorizationCallback(req); pe != nil {
		err = erro.E_OAUTH_FAILED.With(pe)
	} else {
		tok, verifier = _tok, _verifier
	}
	return
}

func (o *_OAuthConfigInternal) OAuth_AccessToken(tok, verifier, request_secret string) (response_values map[string]string, err erro.Error) {
	values, ae := o.c.RetrieveAccessToken(tok, request_secret, verifier)
	if ae != nil {
		err = erro.E_OAUTH_FAILED.With(ae)
	} else {
		response_values = make(map[string]string)
		for k, _ := range values {
			response_values[k] = values.Get(k)
		}
	}
	return
}
