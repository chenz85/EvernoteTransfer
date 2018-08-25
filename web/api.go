package web

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/czsilence/EvernoteTransfer/erro"
	"github.com/czsilence/go/log"
	"github.com/czsilence/go/typo"
	"github.com/gogo/protobuf/jsonpb"
	"github.com/gogo/protobuf/proto"
)

type APIContext struct {
	msg    proto.Message
	header typo.Any
	req    *http.Request
}

func (ctx *APIContext) Req() *http.Request {
	return ctx.req
}

type APIProvider interface {
	Process(ctx *APIContext) (resp_msg proto.Message, err erro.Error)
}

type APIFunc = func(ctx *APIContext) (resp_msg proto.Message, err erro.Error)

type APIItem struct {
	n string
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
			n: name,
			o: privider,
		}
	}
}

func RegisterAPIFunc(name string, f APIFunc) {
	if _, ex := api_map[name]; ex {
		log.E("[web] duplicated api:", name)
	} else {
		api_map[name] = &APIItem{
			n: name,
			f: f,
		}
	}
}

func map_api(r *gin.Engine) {
	for name, item := range api_map {
		r.POST("/api/"+name, item.Handle)
	}
}

func (item *APIItem) Handle(c *gin.Context) {
	log.I("[local] handle api req:", item.n)
	var req = c.Request
	var ctx = &APIContext{
		req: req,
	}
	if resp, re := item.Process(ctx); re != nil {
		log.D("[local] handle api failed:", re)
		// 返回错误信息
		var err_msg = &ErrorMessage{
			ErrCode: re.Code(),
			ErrMsg:  re.Msg(),
		}
		write_response_msg(c, err_msg)
	} else {
		log.D("[local] handle api done:", resp)
		write_response_msg(c, resp)
	}
	return
}

func write_response_msg(c *gin.Context, resp proto.Message) (err erro.Error) {
	if resp == nil {
		resp = default_response_success
	}

	marshaler := &jsonpb.Marshaler{
		OrigName: true,
	}
	if data, me := marshaler.MarshalToString(resp); me != nil {
		err = erro.E_API_MarshalResponseFailed.F("err: %v", me)
	} else {
		c.JSON(http.StatusOK, data)
	}
	return
}
