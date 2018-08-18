package service

import (
	"errors"

	"github.com/czsilence/go/typo"
)

const (
	E_EvernoteConsumerKey    string = "opt:service:eck"
	E_EvernoteConsumerSecret string = "opt:service:ecs"
	E_YinxiangConsumerKey    string = "opt:service:yck"
	E_YinxiangConsumerSecret string = "opt:service:ycs"
)

var (
	E_MissingOption = errors.New("missing option")
	E_InvalidOption = errors.New("invalid option")
)

type ConsumerInfo struct {
	Key    string
	Secret string
}

type ServiceOption struct {
	// evernote的api key信息
	Evernote ConsumerInfo
	// 印象笔记的api key信息
	Yinxiang ConsumerInfo
}

var (
	opt ServiceOption
)

func config(_opt typo.Map) (err error) {
	var ok bool
	if val, ex := _opt[E_EvernoteConsumerKey]; !ex {
		return E_MissingOption
	} else if opt.Evernote.Key, ok = val.(string); !ok || len(opt.Evernote.Key) == 0 {
		return E_InvalidOption
	}
	if val, ex := _opt[E_EvernoteConsumerSecret]; !ex {
		return E_MissingOption
	} else if opt.Evernote.Secret, ok = val.(string); !ok || len(opt.Evernote.Secret) == 0 {
		return E_InvalidOption
	}

	if val, ex := _opt[E_YinxiangConsumerKey]; !ex {
		return E_MissingOption
	} else if opt.Yinxiang.Key, ok = val.(string); !ok || len(opt.Yinxiang.Key) == 0 {
		return E_InvalidOption
	}
	if val, ex := _opt[E_YinxiangConsumerSecret]; !ex {
		return E_MissingOption
	} else if opt.Yinxiang.Secret, ok = val.(string); !ok || len(opt.Yinxiang.Secret) == 0 {
		return E_InvalidOption
	}

	return
}
