package utils

import (
	"errors"
	"time"
)

var layouts = []string{
	time.RFC3339,
	time.DateOnly + "T" + time.TimeOnly,
	time.DateTime,
	time.DateOnly,
}

var errUnableToParseTime = errors.New("unable to parse time")

func ParseTimeFromString(input string) (time.Time, error) {
	var parsed time.Time
	var err error

	for _, layout := range layouts {
		parseInLocal := layout != time.RFC3339

		if parseInLocal {
			parsed, err = time.ParseInLocation(layout, input, time.Local)
		} else {
			parsed, err = time.Parse(layout, input)
		}

		if err != nil {
			// let's try all of the format's first
			continue
		} else {
			return parsed, nil
		}
	}

	return time.Time{}, errUnableToParseTime
}
