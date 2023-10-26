package store

import (
	"encoding/json"
	"fmt"
	"openai/internal/model"
)

func buildErrorsKey(day string) string {
	return fmt.Sprintf("day:%s:errors", day)
}

func AppendError(day string, myErr model.MyError) error {
	errBytes, _ := json.Marshal(myErr)
	err := client.RPush(ctx, buildErrorsKey(day), string(errBytes)).Err()
	if err != nil {
		return err
	}
	err = client.Expire(ctx, buildErrorsKey(day), WEEK).Err()
	return err
}

func GetErrors(day string) ([]model.MyError, error) {
	var myErrors []model.MyError
	errStrs, err := client.LRange(ctx, buildErrorsKey(day), 0, -1).Result()
	if err != nil {
		return myErrors, err
	}

	for _, errStr := range errStrs {
		var chatApiError model.MyError
		_ = json.Unmarshal([]byte(errStr), &chatApiError)
		myErrors = append(myErrors, chatApiError)
	}
	return myErrors, err
}

func GetErrorsLen(day string) (int64, error) {
	return client.LLen(ctx, buildErrorsKey(day)).Result()
}
