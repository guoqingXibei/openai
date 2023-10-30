package logic

import (
	_openai "github.com/sashabaranov/go-openai"
	"openai/internal/constant"
	"openai/internal/service/baidu"
	"openai/internal/service/openai"
	"openai/internal/store"
	"openai/internal/util"
	"strings"
	"unicode"
)

const (
	StartMark     = "[START]"
	EndMark       = "[END]"
	censorWarning = "【温馨提醒】很抱歉识别出了可能不宜讨论的内容。如有误判，可回复contact联系作者，作者将继续进行优化调整。\n\n为了公众号能持续向大家提供服务，请大家不要在这里讨论色情、政治、暴恐、VPN等相关内容，谢谢配合❤️"
)

func ChatCompletionStream(
	aiVendor string, userName string, msgId int64,
	question string, isVoice bool, mode string,
) error {
	_ = store.AppendReplyChunk(msgId, StartMark)
	messages, err := store.GetMessages(userName)
	if err != nil {
		return err
	}

	messages = append(messages, _openai.ChatCompletionMessage{
		Role:    _openai.ChatMessageRoleUser,
		Content: question,
	})
	messages, err = util.RotateMessages(messages, openai.CurrentModel)
	if err != nil {
		return err
	}

	if isVoice {
		_ = store.AppendReplyChunk(msgId, "「"+question+"」\n\n")
	}
	var chunk, answer string
	chunkLen := 60
	isFirstChunk := true
	passedCensor := true
	openai.ChatCompletionsStream(
		aiVendor,
		mode,
		messages,
		func(word string) bool {
			chunk += word
			answer += word
			if len(chunk) >= chunkLen && endsWithPunct(word) || len(chunk) >= 0 {
				chunkLen = 300
				passedCensor, chunk = censorChunk(chunk, isFirstChunk)
				isFirstChunk = false
				_ = store.AppendReplyChunk(msgId, chunk)
				chunk = ""
				if !passedCensor {
					_ = store.AppendReplyChunk(msgId, EndMark)
					return false
				}
			}
			return true
		},
		func() {
			_, chunk = censorChunk(chunk, isFirstChunk)
			_ = store.AppendReplyChunk(msgId, chunk)
			if ShouldAppend(userName) {
				_ = store.AppendReplyChunk(msgId, "\n\n"+selectAppending())
			}
			_ = store.AppendReplyChunk(msgId, EndMark)
			messages = util.AppendAssistantMessage(messages, answer)
			_ = store.SetMessages(userName, messages)
		},
		func(_err error) {
			_ = store.AppendReplyChunk(msgId, constant.TryAgain)
			_ = store.AppendReplyChunk(msgId, EndMark)
			err = _err
		},
	)
	return err
}

func censorChunk(chunk string, isFirstChunk bool) (bool, string) {
	passedCensor := true || baidu.Censor(chunk)
	if !passedCensor {
		if isFirstChunk {
			chunk = censorWarning
		} else {
			chunk = "\n\n" + censorWarning
		}
	}
	return passedCensor, chunk
}

func endsWithPunct(word string) bool {
	runeWord := []rune(word)
	if len(runeWord) <= 0 {
		return false
	}
	return unicode.IsPunct(runeWord[len(runeWord)-1])
}

func FetchReply(msgId int64) (string, bool) {
	chunks, _ := store.GetReplyChunks(msgId, 1, -1)
	if len(chunks) <= 0 {
		return "", false
	}

	reachEnd := chunks[len(chunks)-1] == EndMark
	if reachEnd {
		chunks = chunks[:len(chunks)-1]
	}
	reply := strings.Join(chunks, "")
	return reply, reachEnd
}
