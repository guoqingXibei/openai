package handler

import (
	"encoding/json"
	"fmt"
	"github.com/google/uuid"
	"github.com/redis/go-redis/v9"
	"log"
	"net/http"
	"openai/internal/constant"
	"openai/internal/logic"
	"openai/internal/service/gptredis"
	"openai/internal/service/wechat"
	"strconv"
	"strings"
)

const (
	donate  = "donate"
	help    = "help"
	contact = "contact"
	report  = "report"
)

const (
	generateCode = "generate-code"
	code         = "code:"
)

type CodeDetail struct {
	Code   string `json:"code"`
	Times  int    `json:"times"`
	Status string `json:"status"`
}

const (
	created = "created"
	used    = "used"
)

var keywords = []string{donate, help, contact, report}
var keywordPrefixes = []string{generateCode, code}

func hitKeyword(inMsg *wechat.Msg, writer http.ResponseWriter) bool {
	question := inMsg.Content
	question = strings.TrimSpace(question)
	question = strings.ToLower(question)
	var keyword string
	for _, word := range keywords {
		if question == word {
			keyword = word
			break
		}
	}
	for _, word := range keywordPrefixes {
		if strings.HasPrefix(question, word) {
			keyword = word
			break
		}
	}
	if keyword == "" {
		return false
	}

	switch keyword {
	case contact:
		showContactInfo(inMsg, writer)
	case donate:
		showDonateQr(inMsg, writer)
	case help:
		showUsage(inMsg, writer)
	case report:
		showReport(inMsg, writer)
	case generateCode:
		doGenerateCode(question, inMsg, writer)
	case code:
		useCode(question, inMsg, writer)
	}
	return true
}

func useCode(question string, inMsg *wechat.Msg, writer http.ResponseWriter) {
	code := strings.Replace(question, code, "", 1)
	codeDetailStr, err := gptredis.FetchCodeDetail(code)
	if err == redis.Nil {
		echoWechatTextMsg(writer, inMsg, "无效的code")
		return
	}

	var codeDetail CodeDetail
	_ = json.Unmarshal([]byte(codeDetailStr), &codeDetail)
	if codeDetail.Status == used {
		echoWechatTextMsg(writer, inMsg, "此code之前已被激活，无需重复激活。")
		return
	}

	userName := inMsg.FromUserName
	balance, _ := gptredis.FetchPaidBalance(userName)
	_ = gptredis.SetPaidBalance(userName, codeDetail.Times+balance)
	codeDetail.Status = used
	codeDetailBytes, _ := json.Marshal(codeDetail)
	_ = gptredis.SetCodeDetail(code, string(codeDetailBytes))
	echoWechatTextMsg(writer, inMsg, fmt.Sprintf("此code已被激活，额度为%d。", codeDetail.Times))
}

func doGenerateCode(question string, inMsg *wechat.Msg, writer http.ResponseWriter) {
	fields := strings.Fields(question)
	if len(fields) <= 1 {
		echoWechatTextMsg(writer, inMsg, "Invalid generate-code usage")
		return
	}

	timesStr := fields[1]
	times, err := strconv.Atoi(timesStr)
	if err != nil {
		log.Printf("timesStr is %s, strconv.Atoi error is %v", timesStr, err)
		echoWechatTextMsg(writer, inMsg, "Invalid generate-code usage")
		return
	}

	code := uuid.New().String()
	codeDetail := CodeDetail{
		Code:   code,
		Times:  times,
		Status: created,
	}
	codeDetailBytes, _ := json.Marshal(codeDetail)
	_ = gptredis.SetCodeDetail(code, string(codeDetailBytes))
	echoWechatTextMsg(writer, inMsg, "code:"+code)
}

func showContactInfo(inMsg *wechat.Msg, writer http.ResponseWriter) {
	echoWechatTextMsg(writer, inMsg, constant.ContactInfo)
}

func showReport(inMsg *wechat.Msg, writer http.ResponseWriter) {
	echoWechatTextMsg(writer, inMsg, constant.ReportInfo)
}

func showDonateQr(inMsg *wechat.Msg, writer http.ResponseWriter) {
	QrMediaId, err := wechat.GetMediaIdOfDonateQr()
	if err != nil {
		log.Println("wechat.GetMediaIdOfDonateQr failed", err)
		echoWechatTextMsg(writer, inMsg, constant.TryAgain)
		return
	}
	echoWechatImageMsg(writer, inMsg, QrMediaId)
}

func showUsage(inMsg *wechat.Msg, writer http.ResponseWriter) {
	userName := inMsg.FromUserName
	usage := logic.BuildChatUsage(userName)
	balance, err := gptredis.FetchPaidBalance(userName)
	if err == nil {
		usage += fmt.Sprintf("付费剩余次数为%d。", balance)
	}
	usage += "\n\n" + constant.ContactDesc + "\n" + constant.DonateDesc
	echoWechatTextMsg(writer, inMsg, usage)
}

func switchMode(mode string, inMsg *wechat.Msg, writer http.ResponseWriter) {
	userName := inMsg.FromUserName
	err := gptredis.SetModeForUser(userName, mode)
	if err != nil {
		log.Println("gptredis.SetModeForUser failed", err)
		echoWechatTextMsg(writer, inMsg, constant.TryAgain)
	} else {
		echoWechatTextMsg(writer, inMsg, buildReplyWhenSwitchMode(userName, mode))
	}
}

func buildReplyWhenSwitchMode(userName string, mode string) string {
	reply := "已切换到" + mode + "模式，"
	if mode == constant.Image {
		reply += logic.BuildImageUsage(userName)
	} else {
		reply += logic.BuildChatUsage(userName)
	}
	return reply
}
