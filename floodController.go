package FloodController

import (
	"context"
	"sync"
	"time"
)

type FloodControl interface {
	// Check возвращает false если достигнут лимит максимально разрешенного
	// кол-ва запросов согласно заданным правилам флуд контроля.
	Check(ctx context.Context, userID int64) (bool, error)
}

type FloodController struct {
	callNumber         int
	callExpireDuration time.Duration
	callHistory        sync.Map
}

func (fc *FloodController) Check(ctx context.Context, userID int64) (bool, error) {
	history, _ := fc.callHistory.Load(userID)
	userCallHistory := make([]time.Time, 0)
	if history != nil {
		userCallHistory = history.([]time.Time)
	} else {
		userCallHistory = append(userCallHistory, time.Now())
		fc.callHistory.Store(userID, userCallHistory)
		return true, nil
	}
	newUserCallHistory := make([]time.Time, 0)
	userCallHistory = append(userCallHistory, time.Now())
	for _, call := range userCallHistory {
		if call.After(time.Now().Add(-fc.callExpireDuration)) {
			newUserCallHistory = append(newUserCallHistory, call)
		}
	}
	userCallHistory = newUserCallHistory
	fc.callHistory.Store(userID, userCallHistory)

	if len(userCallHistory) > fc.callNumber {
		return false, nil
	} else {
		return true, nil
	}
}

func MakeFloodController(callNumber int, callExpireDuration time.Duration) FloodController {
	return FloodController{callNumber: callNumber, callExpireDuration: callExpireDuration}
}
