package main

import (
	"encoding/json"
	"flag"
	// "fmt"
	// "time"
	"github.com/260by/qywx/token"
	"github.com/sirupsen/logrus"
	"github.com/tidwall/gjson"
	"github.com/kataras/iris"
	"github.com/iris-contrib/middleware/cors"
	"github.com/kataras/iris/middleware/logger"
	"github.com/kataras/iris/middleware/recover"
	"github.com/260by/qywx/msg"
)

type Message interface {}

var accessToken token.AccessToken

func main() {
	corpID := flag.String("corpID", "", "Weixin corp id")
	agentID := flag.String("agentID", "", "Weixin corp app id")
	appSecret := flag.String("appSecret", "", "Weixin corp app secret")
	flag.Parse()

	go token.Loop(corpID, appSecret, &accessToken)

	app := iris.New()
	app.Logger().SetLevel("debug")
	app.Use(recover.New())
	app.Use(logger.New())
	app.Use(cors.Default())

	app.Post("/api/send", func(ctx iris.Context)  {
		logrus.Infof("Token:", accessToken.Ticket)
		var message Message
		err := ctx.ReadJSON(&message)
		if err != nil {
			ctx.StatusCode(iris.StatusBadRequest)
			ctx.WriteString(err.Error())
			return
		}

		postJSON, err := json.Marshal(message)
		if err != nil {
			logrus.Errorln(err)
		}
		
		msgType := gjson.GetBytes(postJSON, "type").String()
		if msgType == "text" {
			 text := msg.Text{
				ToUser: "@all",
				MsgType: "text",
				AgentID: agentID,
				Safe: 0,
			}
			text.Text.Content = gjson.GetBytes(postJSON, "content").String()
			err := text.Send(accessToken.Ticket)
			if err != nil {
				logrus.Errorln(err)
			} else {
				logrus.Infoln("Send message is OK.")
			}
		}
		if msgType == "textcard" {
			var textCard msg.TextCard
			err = json.Unmarshal(postJSON, &textCard)
			if err != nil {
				logrus.Errorln(err)
			}
		}
		

	})

	app.Run(iris.Addr(":9000"), iris.WithoutServerError(iris.ErrServerClosed))
}
