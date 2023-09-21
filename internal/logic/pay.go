package logic

import (
	"github.com/go-pay/gopay/pkg/util"
	"log"
	"openai/internal/model"
	"openai/internal/service/gptredis"
	"openai/internal/service/wechat"
	"time"
)

func InitiateTransaction(
	openid string,
	uncleOpenid string,
	priceInFen int,
	times int,
	description string) (string, string, error) {
	tradeNo := util.RandomString(32)
	log.Println("tradeNo:", tradeNo)
	prepayId, err := wechat.InitiateTransaction(openid, tradeNo, priceInFen, description)
	if err != nil {
		return "", "", err
	}

	now := time.Now().Unix()
	transaction := model.Transaction{
		OutTradeNo:  tradeNo,
		OpenId:      openid,
		UncleOpenId: uncleOpenid,
		PrepayId:    prepayId,
		TradeState:  "",
		Redeemed:    false,
		PriceInFen:  priceInFen,
		Times:       times,
		Description: description,
		CreatedTime: now,
		UpdatedTime: now,
	}
	_ = gptredis.SetTransaction(tradeNo, transaction)
	return prepayId, tradeNo, err
}