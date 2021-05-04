package util

import (
	"errors"
	"regexp"
	"time"
)

func ConvertCondensedDateString(s string) (time.Time, error) {
	errTime := time.Date(0, 0, 0, 0, 0, 0, 0, time.UTC)

	re := regexp.MustCompile(`^(\d{4})(\d{2})(\d{2})$`)

	ms := re.FindStringSubmatch(s)
	if len(ms) != 4 {
		return errTime, errors.New("string does not match YYYYMMDD format")
	}

	return time.Parse("2006-01-02", ms[1]+"-"+ms[2]+"-"+ms[3])
}
