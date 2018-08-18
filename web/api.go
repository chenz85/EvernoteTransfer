package web

import (
	"net/http"

	"github.com/czsilence/EvernoteTransfer/erro"
	"github.com/czsilence/go/log"
	"github.com/czsilence/go/typo"
	"github.com/gogo/protobuf/proto"
)

type APIProvider interface {
	Process(msg proto.Message, header typo.Any) (resp_msg proto.Message, err erro.Error)
}

var (
	api_map map[string]APIProvider
)

func init() {
	api_map = make(map[string]APIProvider)
}

func RegisterAPIProvider(name string, privider APIProvider) {
	if _, ex := api_map[name]; ex {
		log.E("[web] duplicated api:", name)
	} else {
		api_map[name] = privider
	}
}

func map_api(name string, w http.ResponseWriter, req *http.Request) (err erro.Error) {
	if provider, ex := api_map[name]; ex {
		log.I("[local] handle api req:", name)
		if resp, re := provider.Process(nil, nil); re != nil {
			log.D("[local] handle api failed:", re)
			// 返回错误信息
			var err_msg = &ErrorMessage{
				ErrCode: re.Code(),
				ErrMsg:  re.Msg(),
			}
			err = write_response_msg(w, err_msg)
		} else {
			log.D("[local] handle api failed:", resp)
			err = write_response_msg(w, resp)
		}
	} else {
		err = erro.E_API_UnknownAPI.F("name")
	}
	return
}

func write_response_msg(w http.ResponseWriter, msg proto.Message) (err erro.Error) {
	if data, me := proto.Marshal(msg); me != nil {
		err = erro.E_API_MarshalResponseFailed.F("err: %v", me)
	} else {
		var cnt = 0
		for cnt < len(data) {
			if _cnt, we := w.Write(data[cnt:]); we != nil {
				err = erro.E_API_WriteResponseFailed.F("err: %v", we)
				break
			} else {
				cnt += _cnt
			}
		}
	}
	return
}
