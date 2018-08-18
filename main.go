package main

import (
	"github.com/czsilence/EvernoteTransfer/service"
	"github.com/czsilence/EvernoteTransfer/typo"
)

const (
	// Evernote 的 api key 信息，从 https://dev.evernote.com/ 获取
	E_CONSUMER_KEY    string = "<put consumer key here>"
	E_CONSUMER_SECRET string = "<put consumer secret here>"

	// 印象笔记的 api key 信息，从 https://dev.yinxiang.com/ 获取
	Y_CONSUMER_KEY    string = "<put consumer key here>"
	Y_CONSUMER_SECRET string = "<put consumer secret here>"
)

func main() {
	service.Start(typo.Map{
		service.E_EvernoteConsumerKey:    E_CONSUMER_KEY,
		service.E_EvernoteConsumerSecret: E_CONSUMER_SECRET,
		service.E_YinxiangConsumerKey:    Y_CONSUMER_KEY,
		service.E_YinxiangConsumerSecret: Y_CONSUMER_SECRET,
	})
}
