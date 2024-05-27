package logic

import (
	"fmt"
	"github.com/robfig/cron"
	"log"
	"openai/internal/constant"
	"openai/internal/model"
	"openai/internal/service/email"
	"openai/internal/service/errorx"
	"openai/internal/service/ohmygpt"
	"openai/internal/service/sb"
	"openai/internal/store"
	"openai/internal/util"
	"slices"
)

func init() {
	c1 := cron.New()
	// Execute once every day at 00:00
	err := c1.AddFunc("0 0 0 * * ?", func() {
		sendYesterdayReportEmail()
	})
	if err != nil {
		errorx.RecordError("AddFunc() failed", err)
		return
	}
	c1.Start()

	if util.AccountIsBrother() && util.EnvIsProd() {
		c2 := cron.New()
		// Execute once every hour
		err = c2.AddFunc("0 0 * * * *", func() {
			checkVendorBalance()
		})
		if err != nil {
			errorx.RecordError("AddFunc() failed", err)
			return
		}
		c2.Start()
	}
}

func checkVendorBalance() {
	log.Println("Checking balance of vendors...")
	alarm := false
	ohmygptBalance, _ := ohmygpt.GetOhmygptBalance()
	if ohmygptBalance < 30 {
		alarm = true
	}
	sbBalance, _ := sb.GetSbBalance()
	if sbBalance < 0.1 {
		alarm = true
	}
	if alarm {
		log.Println("Balance is insufficient, sending email...")
		body := buildBalanceSection(ohmygptBalance, sbBalance)
		email.SendEmail("Insufficient Balance", body)
	}
	log.Println("Check finished")
}

func sendYesterdayReportEmail() {
	yesterday := util.Yesterday()
	subject := fmt.Sprintf("[%s/%s] Summary on %s", util.GetAccount(), util.GetEnv(), yesterday)

	body := ""
	var purchasedUsers []string
	if util.AccountIsBrother() {
		ohmygptBalance, _ := ohmygpt.GetOhmygptBalance()
		sbBalance, _ := sb.GetSbBalance()
		balanceSect := buildBalanceSection(ohmygptBalance, sbBalance)
		body += balanceSect

		tradeNos, _ := store.GetSuccessOutTradeNos(yesterday)
		cnt := len(tradeNos)
		txnTitle := fmt.Sprintf("\n# %d purchases\n", cnt)
		txnContent := ""
		if cnt > 0 {
			txnContent = `
Time | User | Amount | Account
-----|------|--------|--------
`
			colTmpl := "%s | %s | %.1f | %s\n"
			for _, tradeNo := range tradeNos {
				transaction, _ := store.GetTransaction(tradeNo)
				paidAccount := constant.Brother
				openId := transaction.OpenId
				if transaction.UncleOpenId != "" {
					paidAccount = constant.Uncle
					openId = transaction.UncleOpenId
				}
				purchasedUsers = append(purchasedUsers, openId)
				txnContent += fmt.Sprintf(colTmpl,
					util.FormatTime(transaction.UpdatedAt),
					openId,
					transaction.PriceInFen/100,
					paidAccount,
				)
			}
		}
		body += txnTitle + txnContent
	}

	errCnt, errorContent := errorx.GetErrorsDesc(yesterday)
	errorTitle := fmt.Sprintf("\n# %d errors\n", errCnt)
	body += errorTitle + errorContent

	users, _ := store.GetActiveUsers(yesterday)
	user2Convs := make(map[string][]model.Conversation)
	for _, user := range users {
		convs, _ := store.GetConversations(user, yesterday)
		user2Convs[user] = convs
	}
	compareUserFunc := func(a, b string) int {
		aPurchased := slices.Contains(purchasedUsers, a)
		bPurchased := slices.Contains(purchasedUsers, b)
		if !aPurchased && !bPurchased || aPurchased && bPurchased {
			return user2Convs[a][0].Time.Compare(user2Convs[b][0].Time)
		}

		if aPurchased {
			return -1
		}

		return 1
	}
	slices.SortFunc(users, compareUserFunc)

	userCnt := len(users)
	convCnt := 0
	convContent := ""
	for idx, user := range users {
		convs := user2Convs[user]
		convsCnt := len(convs)
		mark := ""
		if slices.Contains(purchasedUsers, user) {
			mark = " ★"
		}
		convContent += fmt.Sprintf("## %d/%d %s %d %s\n", idx+1, userCnt, user, convsCnt, mark)
		for convIdx, conv := range convs {
			convTmpl := `
### %d/%d %s %s %d
**Q**: %s
**A**: %s
`
			convContent += fmt.Sprintf(convTmpl,
				convIdx+1,
				convsCnt,
				util.FormatTime(conv.Time),
				conv.Mode,
				conv.PaidBalance,
				truncateAndEscape(conv.Question),
				truncateAndEscape(conv.Answer),
			)
		}
		convCnt += convsCnt
	}
	convTitle := fmt.Sprintf("\n# %d users | %d convs\n", userCnt, convCnt)
	body += convTitle + convContent

	email.SendEmail(subject, body)
}

func buildBalanceSection(ohmygptBalance float64, sbBalance float64) string {
	balanceTmpl := `
# Balance
Vendor | Balance
-------|------
%s | %.2f
%s | %.2f
`
	balanceSect := fmt.Sprintf(balanceTmpl,
		"Ohmygpt", ohmygptBalance,
		"SB", sbBalance,
	)
	return balanceSect
}

func truncateAndEscape(origin string) string {
	maxLen := 100
	if util.IsEnglishSentence(origin) {
		maxLen *= 5
	}
	return util.EscapeNewline(util.EscapeHtmlTags(util.TruncateString(origin, maxLen)))
}
