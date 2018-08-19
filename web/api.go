package web

import (
	"fmt"
	"net/http"
	"net/url"

	"github.com/czsilence/EvernoteTransfer/erro"
	"github.com/czsilence/go/log"
	"github.com/czsilence/go/typo"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
)

type APIContext struct {
	msg    proto.Message
	header typo.Any
	Querys url.Values
}

type APIProvider interface {
	Process(ctx *APIContext) (resp_msg proto.Message, err erro.Error)
}

type APIFunc = func(ctx *APIContext) (resp_msg proto.Message, err erro.Error)

type APIItem struct {
	o APIProvider
	f APIFunc
}

func (i *APIItem) Process(ctx *APIContext) (resp_msg proto.Message, err erro.Error) {
	if i.f != nil {
		return i.f(ctx)
	} else {
		return i.o.Process(ctx)
	}
}

var (
	api_map map[string]*APIItem

	// 默认处理成功的返回值
	default_response_success = &ErrorMessage{
		ErrCode: 0,
		ErrMsg:  "sucess",
	}
)

func init() {
	api_map = make(map[string]*APIItem)
}

func RegisterAPIProvider(name string, privider APIProvider) {
	if _, ex := api_map[name]; ex {
		log.E("[web] duplicated api:", name)
	} else {
		api_map[name] = &APIItem{
			o: privider,
		}
	}
}

func RegisterAPIFunc(name string, f APIFunc) {
	if _, ex := api_map[name]; ex {
		log.E("[web] duplicated api:", name)
	} else {
		api_map[name] = &APIItem{
			f: f,
		}
	}
}

func map_api(name string, w http.ResponseWriter, req *http.Request) (err erro.Error) {
	if item, ex := api_map[name]; ex {
		log.I("[local] handle api req:", name)
		var ctx = &APIContext{
			Querys: req.URL.Query(),
		}
		if resp, re := item.Process(ctx); re != nil {
			log.D("[local] handle api failed:", re)
			// 返回错误信息
			var err_msg = &ErrorMessage{
				ErrCode: re.Code(),
				ErrMsg:  re.Msg(),
			}
			err = write_response_msg(w, err_msg)
		} else {
			log.D("[local] handle api done:", resp)
			err = write_response_msg(w, resp)
		}
	} else {
		log.I("[local] unknown api:", name)
		err = erro.E_API_UnknownAPI.F("name")
	}
	return
}

func write_response_msg(w http.ResponseWriter, resp proto.Message) (err erro.Error) {
	if resp == nil {
		resp = default_response_success
	}

	marshaler := &jsonpb.Marshaler{
		OrigName: true,
	}
	if data, me := marshaler.MarshalToString(resp); me != nil {
		err = erro.E_API_MarshalResponseFailed.F("err: %v", me)
	} else {
		fmt.Fprint(w, data)
		// var cnt = 0
		// for cnt < len(data) {
		// 	if _cnt, we := w.Write(data[cnt:]); we != nil {
		// 		err = erro.E_API_WriteResponseFailed.F("err: %v", we)
		// 		break
		// 	} else {
		// 		cnt += _cnt
		// 	}
		// }
	}
	return
}
