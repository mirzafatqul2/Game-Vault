package date

import (
	"errors"
	"time"
)

func ParseDateOfBirth(dateStr string) (time.Time, error) {

	layout := "2006-01-02"

	t, err := time.Parse(layout, dateStr)
	if err != nil {
		return time.Time{}, errors.New("date format must be YYYY-MM-DD")
	}

	return t, nil
}
