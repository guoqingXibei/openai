package logic

import (
	"log"
	"openai/internal/constant"
	"openai/internal/service/gptredis"
)

const (
	triggerTimes = 20
)

func incUsedTimes(user string) int {
	times, err := gptredis.IncUsedTimes(user)
	if err != nil {
		log.Println("gptredis.IncUsedTimes failed", err)
		return 0
	}
	return times
}

func ShouldAppend(user string) bool {
	times := incUsedTimes(user)
	return times%triggerTimes == 0
}

func appendHelpDesc(answer string) string {
	return answer + "\n\n" + constant.DonateReminder
}

func AppendIfPossible(user string, answer string) string {
	if ShouldAppend(user) {
		return appendHelpDesc(answer)
	}
	return answer
}
