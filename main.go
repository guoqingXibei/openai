package main

import (
	"log"
	"net/http"
	"openai/bootstrap"
	"openai/internal/config"
	"openai/internal/constant"
	"openai/internal/handler"
	wechat2 "openai/internal/service/wechat"
	"openai/internal/util"
	"os"
)

var (
	env = os.Getenv("GO_ENV")
)

func main() {
	ConfigLog()

	engine := bootstrap.New()
	// 公众号消息处理
	engine.POST("/talk", serveWechat)
	// 用于公众号自动验证
	engine.GET("/talk", serveWechat)
	// Provide reply content for the webpage
	engine.GET("/reply-stream", handler.GetReplyStream)
	engine.GET("/openid", handler.GetOpenId)
	engine.POST("/transactions", handler.Transaction)
	engine.POST("/notify-transaction-result", handler.NotifyTransactionResult)
	engine.GET("/trade-result", handler.GetTradeResult)
	handlerWithRequestLog := bootstrap.LogRequestHandler(engine)

	http.Handle("/answer/", http.StripPrefix("/answer", http.FileServer(http.Dir("./public"))))
	http.Handle("/images/", http.FileServer(http.Dir("./public")))
	http.Handle("/", handlerWithRequestLog)

	log.Println("Server started in port " + config.C.Http.Port)
	err := http.ListenAndServe("127.0.0.1:"+config.C.Http.Port, nil)
	if err != nil {
		panic(err)
	}
}

func serveWechat(rw http.ResponseWriter, req *http.Request) {
	officialAccount := wechat2.GetAccount()

	// 传入request和responseWriter
	server := officialAccount.GetServer(req, rw)
	server.SetParseXmlToMsgFn(util.ParseXmlToMsg)

	//设置接收消息的处理方法
	server.SetMessageHandler(handler.Talk)

	//处理消息接收以及回复
	err := server.Serve()
	if err != nil {
		log.Println("server.Serve() failed", err)
		err = server.BuildResponse(util.BuildTextReply(constant.TryAgain))
		if err != nil {
			log.Println("server.BuildResponse() failed", err)
			return
		}
	}

	//发送回复的消息
	err = server.Send()
	if err != nil {
		log.Println("server.Send() failed", err)
	}
}

func ConfigLog() {
	if env == "dev" {
		log.SetOutput(os.Stdout)
	} else {
		dir := "./log"
		path := dir + "/data.log"
		_, err := os.Stat(dir)
		if err != nil && os.IsNotExist(err) {
			if err := os.Mkdir(dir, 0755); err != nil {
				panic(err)
			}
		}
		file, err := os.OpenFile(path, os.O_CREATE|os.O_APPEND|os.O_RDWR, 0755)
		if err != nil {
			panic(err)
		}
		log.SetOutput(file)
	}
}
