package web

import (
	"flag"
	"path/filepath"

	"github.com/gin-gonic/gin"
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
	r := gin.Default()

	map_api(r)
	var index_file string = filepath.Join(webroot, "index.html")
	r.StaticFile("/", index_file)
	r.StaticFile("/index.html", index_file)
	r.Static("/js", filepath.Join(webroot, "js"))
	r.Run("127.0.0.1:8001")
}
