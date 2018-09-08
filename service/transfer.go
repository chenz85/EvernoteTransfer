package service

import (
	"context"

	"git.apache.org/thrift.git/lib/go/thrift"

	"github.com/czsilence/EvernoteTransfer/erro"
	"github.com/czsilence/evernote-sdk-go/evernote"
	"github.com/czsilence/go/log"
)

////////////////////////////////////////////////////////

type transport_info struct {
	oauth *OauthInfo
	t     thrift.TTransport
	clt   *evernote.NoteStoreClient
}

func init_transport(o *OauthInfo) (info *transport_info, err erro.Error) {

	if trans, te := thrift.NewTHttpClient(o.EdamNoteStoreUrl); te != nil {
		err = erro.E_Transfer_InitTransportFailed.F("err: %v", te)
	} else if pf := thrift.NewTBinaryProtocolFactory(true, true); pf == nil {
		err = erro.E_Transfer_InitProtocolFailed
	} else if clt := evernote.NewNoteStoreClientFactory(trans, pf); clt == nil {
		err = erro.E_Transfer_InitClientFailed
	} else {
		info = &transport_info{
			oauth: o,
			t:     trans,
			clt:   clt,
		}
	}
	return
}

////////////////////////////////////////////////////////

type Transfer struct {
	from, to *transport_info
	ctx      context.Context

	tag_map      map[evernote.GUID]evernote.GUID
	notebook_map map[evernote.GUID]evernote.GUID

	last_err erro.Error
}

func (t *Transfer) Start(from, to *OauthInfo) (err erro.Error) {
	if t.from, err = init_transport(from); err != nil {
	} else if t.to, err = init_transport(to); err != nil {
	} else {
		t.ctx = context.Background()
		go t.transfer()
	}
	return
}

func (t *Transfer) transfer() {
	for {
		if t.t_tags(); t.last_err != nil {
			break
		}
		if t.t_notebooks(); t.last_err != nil {
			break
		}
		if t.t_notes(); t.last_err != nil {
			break
		}
	}
}

func (t *Transfer) t_tags() {
	if tags, te := t.from.clt.ListTags(t.ctx, t.from.oauth.AccessToken); te != nil {
		t.last_err = erro.E_Transfer_Tag_ListTagsFailed.F("err: %v", te)
	} else {
		log.D("tags:", len(tags))
		for _, tag := range tags {
			log.D2("tag [%s]: %s", tag.GetGUID(), tag.GetName())
			if new_tag, ne := t.to.clt.CreateTag(t.ctx, t.to.oauth.AccessToken, tag); ne != nil {
				t.last_err = erro.E_Transfer_Tag_CreateTagFailed.F("err: %v", ne)
				break
			} else {
				t.tag_map[tag.GetGUID()] = new_tag.GetGUID()
				log.D2("-> new tag [%s]: %s", new_tag.GetGUID(), new_tag.GetName())
			}
		}
	}
}

func (t *Transfer) t_notebooks() {
	if nbs, nbse := t.from.clt.ListNotebooks(t.ctx, t.from.oauth.AccessToken); nbse != nil {
		t.last_err = erro.E_Transfer_NB_ListNotebookFailed.F("err: %v", nbse)
	} else {
		log.D("notebooks:", len(nbs))
		for _, notebook := range nbs {
			log.D2("notebook [%s]: %s", notebook.GetGUID(), notebook.GetName())
			if new_nb, ne := t.to.clt.CreateNotebook(t.ctx, t.to.oauth.AccessToken, notebook); ne != nil {
				t.last_err = erro.E_Transfer_NB_CreateNotebookFailed.F("err: %v", ne)
			} else {
				t.notebook_map[notebook.GetGUID()] = new_nb.GetGUID()
				log.D2("-> new tag [%s]: %s", new_nb.GetGUID(), new_nb.GetName())
			}
		}
	}
}

func (t *Transfer) t_notes() {
	filter := evernote.NewSyncChunkFilter()
	include_notes := true
	filter.IncludeNotes = &include_notes
	if chunk, ce := t.from.clt.GetFilteredSyncChunk(t.ctx, t.from.oauth.AccessToken, 0, 5, filter); ce != nil {
		t.last_err = erro.E_Transfer_Note_GetSyncChunkFailed.F("err: %v", ce)
	} else {
		log.D2("sync chunk: %d", chunk.GetChunkHighUSN())
		for _, note := range chunk.Notes {
			log.D2("note [%s]: %s", note.GetGUID(), note.GetTitle())

			// get note content
			// if note_data, ne := t.from.clt.GetNote(t.ctx, t.from.oauth.AccessToken, note.GetGUID(), true, false, false, true); ne == nil {
			// 	log.W2("note, content: %s", note_data.GetContent())
			// }
		}
	}

}
