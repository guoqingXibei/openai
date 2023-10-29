module openai

go 1.16

require (
	github.com/bradfitz/gomemcache v0.0.0-20230905024940-24af94b03874 // indirect
	github.com/bsm/redislock v0.9.4
	github.com/disintegration/imaging v1.6.2
	github.com/dlclark/regexp2 v1.10.0 // indirect
	github.com/felixge/httpsnoop v1.0.3
	github.com/go-http-utils/headers v0.0.0-20181008091004-fed159eddc2a
	github.com/go-pay/gopay v1.5.97
	github.com/google/uuid v1.4.0
	github.com/redis/go-redis/v9 v9.2.1
	github.com/robfig/cron v1.2.0
	github.com/sashabaranov/go-openai v1.16.0
	github.com/silenceper/wechat/v2 v2.1.5
	github.com/sirupsen/logrus v1.9.3
	github.com/spf13/cast v1.5.1 // indirect
	github.com/stretchr/testify v1.8.2 // indirect
	github.com/tidwall/gjson v1.17.0 // indirect
	github.com/tidwall/pretty v1.2.1 // indirect
	github.com/tiktoken-go/tokenizer v0.1.0
	golang.org/x/image v0.13.0 // indirect
	golang.org/x/sync v0.4.0
)

replace github.com/silenceper/wechat/v2 => github.com/guoqingxibei/wechat/v2 v2.2.7
