package handler

import (
	"errors"
	"fmt"
	"github.com/silenceper/wechat/v2/officialaccount/message"
	"openai/internal/config"
	"openai/internal/constant"
	"openai/internal/logic"
	"openai/internal/service/errorx"
	"openai/internal/service/wechat"
	"openai/internal/store"
	"openai/internal/util"
	"runtime/debug"
	"strconv"
	"strings"
	"time"
	"unicode"
)

const (
	maxLengthOfReply     = 4000
	maxRuneLengthOfReply = 200
	maxLengthOfQuestion  = 3000 // ~ 1000 Chinese characters
)

func onReceiveText(msg *message.MixMessage) (reply *message.Reply) {
	if msg.MsgType == message.MsgTypeVoice {
		if msg.Recognition == "" {
			reply = util.BuildTextReply("抱歉，未识别到有效内容。")
			return
		}
		msg.Content = msg.Recognition
	}

	if len(msg.Content) > maxLengthOfQuestion {
		reply = util.BuildTextReply("哎呀，输入太长了~")
		return
	}

	hit, reply := hitKeyword(msg)
	if hit {
		return
	}

	// when WeChat server retries
	msgID := msg.MsgID
	times, _ := store.IncRequestTimesForMsg(msgID)
	if times > 1 {
		mode, _ := store.GetMode(string(msg.FromUserName))
		reply = util.BuildTextReply(buildLateReply(msgID, mode))
		return
	}

	reply = util.BuildTextReply(genReply4Text(msg))
	return
}

func genReply4Text(msg *message.MixMessage) (reply string) {
	msgId := msg.MsgID
	user := string(msg.FromUserName)
	question := strings.TrimSpace(msg.Content)
	mode, _ := store.GetMode(user)
	ok, balanceTip := logic.DecreaseBalance(user, mode)
	if !ok {
		reply = balanceTip
		return
	}

	drawReplyIsLate := false
	replyChan := make(chan string, 1)
	go func() {
		defer func() {
			if r := recover(); r != nil {
				panicMsg := fmt.Sprintf("%v\n%s", r, debug.Stack())
				errorx.RecordError("failed due to a panic", errors.New(panicMsg))
			}
		}()

		isVoice := msg.MsgType == message.MsgTypeVoice
		if mode == constant.Draw {
			drawReply := logic.SubmitDrawTask(question, user, mode, isVoice)
			replyChan <- drawReply
			if drawReplyIsLate {
				err := wechat.GetAccount().GetCustomerMessageManager().
					Send(message.NewCustomerTextMessage(user, drawReply))
				if err != nil {
					errorx.RecordError("GetCustomerMessageManager().Send() failed", err)
				}
			}
			return
		}

		logic.CreateChatStreamEx(user, msgId, question, isVoice, mode)
		replyChan <- buildReplyForChat(msgId)
	}()
	select {
	case reply = <-replyChan:
	case <-time.After(time.Millisecond * 2500):
		if mode == constant.Draw {
			drawReplyIsLate = true
		}
		reply = buildLateReply(msgId, mode)
	}
	return
}

func buildLateReply(msgId int64, mode string) (reply string) {
	if mode == constant.Draw {
		reply = "正在提交绘画任务，静候佳音..."
	} else {
		reply = buildReplyForChat(msgId)
	}
	return
}

func buildReplyForChat(msgId int64) string {
	reply, reachEnd := logic.FetchReply(msgId)
	if len(reply) > maxLengthOfReply {
		reply = buildReplyWithShowMore(string([]rune(reply)[:maxRuneLengthOfReply]), msgId)
	} else {
		if reachEnd {
			// Intent to display internal images via web
			if strings.Contains(reply, "![](./images/") {
				runes := []rune(reply)
				length := len(runes)
				if length > 30 {
					length = 30
				}
				reply = buildReplyWithShowMore(string(runes[:length]), msgId)
			}
		} else {
			if reply == "" {
				reply = buildReplyURL(msgId, "查看回复")
			} else {
				reply = buildReplyWithShowMore(reply, msgId)
			}
		}
	}
	return strings.TrimSpace(reply)
}

func buildReplyWithShowMore(answer string, msgId int64) string {
	return trimTailingPuncts(answer) + "..." + buildReplyURL(msgId, "更多")
}

func trimTailingPuncts(answer string) string {
	runeAnswer := []rune(answer)
	if len(runeAnswer) <= 0 {
		return ""
	}
	tailIdx := -1
	for i := len(runeAnswer) - 1; i >= 0; i-- {
		if !unicode.IsPunct(runeAnswer[i]) {
			tailIdx = i
			break
		}
	}
	return string(runeAnswer[:tailIdx+1])
}

func buildReplyURL(msgId int64, desc string) string {
	url := config.C.Wechat.MessageUrlPrefix + "/answer/#/?msgId=" + strconv.FormatInt(msgId, 10)
	return fmt.Sprintf("<a href=\"%s\">%s</a>", url, desc)
}
