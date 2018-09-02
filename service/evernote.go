package service

import (
	"context"
	"net/http"

	"git.apache.org/thrift.git/lib/go/thrift"
	"github.com/czsilence/EvernoteTransfer/erro"
	"github.com/czsilence/EvernoteTransfer/storage"
	"github.com/czsilence/EvernoteTransfer/storage/badger"
	"github.com/czsilence/EvernoteTransfer/web"
	"github.com/czsilence/evernote-sdk-go/evernote"
	"github.com/czsilence/go/log"
	"github.com/gogo/protobuf/proto"
)

func init() {
	web.RegisterAPIFunc("en/oauth", http.MethodPost, _api_oauth)
	web.RegisterAPIFunc("en/oauth/callback", http.MethodGet, _api_oauth_callback)
	web.RegisterAPIFunc("en/user", http.MethodPost, _api_en_user)
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

//////////////////////////
func _api_en_user(ctx *web.APIContext) (resp_msg proto.Message, err erro.Error) {
	var sid = ctx.Sid()
	var oauth_info = new(OauthInfo)
	if !storage.Get(badger.DB(), "oauth:"+sid, oauth_info) {
		err = erro.E_OAuth_OAuthInfoNotFound
	} else {
		resp_msg = oauth_info
		_user_info(oauth_info)
	}
	return
}

func _user_info(o *OauthInfo) (user_info proto.Message, err erro.Error) {

	if trans, te := thrift.NewTHttpClient(o.EdamNoteStoreUrl); te != nil {

	} else if pf := thrift.NewTBinaryProtocolFactory(true, true); pf == nil {

	} else if clt := evernote.NewNoteStoreClientFactory(trans, pf); clt == nil {

	} else {
		ctx := context.Background()
		if tags, te := clt.ListTags(ctx, o.AccessToken); te != nil {

		} else {
			log.W("tags:", len(tags))
			for i, tag := range tags {
				log.W2("tag [#%d]: %s", i, tag.GetName())
			}
		}

		if nbs, nbse := clt.ListNotebooks(ctx, o.AccessToken); nbse != nil {

		} else {
			log.W("notebooks:", len(nbs))
			for i, notebook := range nbs {
				log.W2("notebook [#%d]: %s", i, notebook.GetName())
				// log.W2("notebook: %+v", notebook)
				log.W2("notebook USN: %d", notebook.GetUpdateSequenceNum())

				filter := evernote.NewSyncChunkFilter()
				include_notes := true
				filter.IncludeNotes = &include_notes
				if chunk, ce := clt.GetFilteredSyncChunk(ctx, o.AccessToken, 0, 5, filter); ce != nil {
				} else {
					log.W2("sync chunk: %+v", chunk)
				}
			}
		}

		if state, se := clt.GetSyncState(ctx, o.AccessToken); se != nil {
		} else {
			log.W("sync state:", state)
		}

	}
	return
}
