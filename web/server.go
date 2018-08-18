package web

import (
	"flag"
	"net/http"
	"path/filepath"
	"strings"

	"github.com/czsilence/go/log"
)

var (
	// 默认使用相对路径定位
	webroot string = "webroot"
)

func init() {
	flag.StringVar(&webroot, "webroot", webroot, "specify web app root.")
}

func StartLocalHttpServer() {
	go start_internal()
}

func start_internal() {
	http.ListenAndServe("127.0.0.1:8001", new(_Handler))
}

type _Handler struct {
}

func (h *_Handler) ServeHTTP(w http.ResponseWriter, req *http.Request) {
	if req.URL.Path == "/" {
		index(w, req)
	} else if strings.Index(req.URL.Path, "/api/") == 0 {
		// TODO: handle api request
		log.I("[local] handle api req:", req.URL)
	} else if strings.HasPrefix(req.URL.Path, "/") {
		var f string = filepath.Join(webroot, req.URL.Path)
		http.ServeFile(w, req, f)
	}
}

func index(w http.ResponseWriter, req *http.Request) {
	var index_file string = filepath.Join(webroot, "index.html")
	log.D2("[local] index file: %s", index_file)
	http.ServeFile(w, req, index_file)
}
