package wechat

import (
	"context"
	"errors"
	"github.com/go-pay/gopay"
	"github.com/go-pay/gopay/wechat/v3"
	"log"
	"time"
)

var client *wechat.ClientV3
var ctx = context.Background()

func init() {
	// NewClientV3 初始化微信客户端 v3
	// mchid：商户ID 或者服务商模式的 sp_mchid
	// serialNo：商户证书的证书序列号
	// apiV3Key：apiV3Key，商户平台获取
	// privateKey：私钥 apiclient_key.pem 读取后的内容
	var err error
	client, err = wechat.NewClientV3(wechatConfig.MchId, wechatConfig.SerialNo, wechatConfig.APIv3Key, wechatConfig.PrivateKey)
	if err != nil {
		log.Println(err)
		return
	}

	// 设置微信平台API证书和序列号（推荐开启自动验签，无需手动设置证书公钥等信息）
	//client.SetPlatformCert([]byte(""), "")

	// 启用自动同步返回验签，并定时更新微信平台API证书（开启自动验签时，无需单独设置微信平台API证书和序列号）
	err = client.AutoVerifySign()
	if err != nil {
		log.Println(err)
		return
	}

	// 自定义配置http请求接收返回结果body大小，默认 10MB
	//client.SetBodySize() // 没有特殊需求，可忽略此配置

	// 打开Debug开关，输出日志，默认是关闭的
	client.DebugSwitch = gopay.DebugOn
}

func InitiateTransaction(openid string, tradeNo string, total int, description string) (string, error) {
	expire := time.Now().Add(10 * time.Minute).Format(time.RFC3339)
	bm := make(gopay.BodyMap)
	bm.Set("appid", wechatConfig.AppId).
		Set("description", description).
		Set("out_trade_no", tradeNo).
		Set("time_expire", expire).
		Set("notify_url", wechatConfig.NotifyUrl).
		SetBodyMap("amount", func(bm gopay.BodyMap) {
			bm.Set("total", total).
				Set("currency", "CNY")
		}).
		SetBodyMap("payer", func(bm gopay.BodyMap) {
			bm.Set("openid", openid)
		})

	wxRsp, err := client.V3TransactionJsapi(ctx, bm)
	if err != nil {
		log.Println(err)
		return "", err
	}
	if wxRsp.Code != wechat.Success {
		log.Printf("wxRsp:%s", wxRsp.Error)
		return "", errors.New("wxRsp error")
	}
	log.Printf("wxRsp: %+v", wxRsp.Response)
	return wxRsp.Response.PrepayId, nil
}

func VerifySignAndDecrypt(notifyReq *wechat.V3NotifyReq) (*wechat.V3DecryptResult, error) {
	// 获取微信平台证书
	certMap := client.WxPublicKeyMap()
	// 验证异步通知的签名
	err := notifyReq.VerifySignByPKMap(certMap)
	if err != nil {
		return nil, err
	}

	return notifyReq.DecryptCipherText(wechatConfig.APIv3Key)
}

func GeneratePaySignParams(prepayid string) (*wechat.JSAPIPayParams, error) {
	return client.PaySignOfJSAPI(wechatConfig.AppId, prepayid)
}