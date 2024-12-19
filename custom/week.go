package custom

import (
	"time"
)

type WeekUtil interface {
	GetNextStart() int
}

type weekUtil struct {
	days int
}

func NewWeekUtil(days int) WeekUtil {
	return &weekUtil{
		days: days,
	}
}

func (u *weekUtil) GetNextStart() int {
	currentDate := time.Now().Round(time.Hour)
	modifiedDate := currentDate.AddDate(0, 0, -u.days)
	daysUntilNextMonday := (8 - int(modifiedDate.Weekday())) % 7
	totalAdjustedDays := u.days - daysUntilNextMonday
	return totalAdjustedDays
}
