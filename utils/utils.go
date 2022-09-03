package utils

import (
	"errors"
	"regexp"
)

type Times struct {
	Date     string
	TimeFrom string
	TimeTo   string
}

func ParseTimesFromString(stringPeriod string) (times Times, err error) {
	re, _ := regexp.Compile(`(?P<date>\d*\.*\d*.*\d{4})\s*(?P<timeFrom>\d{1,2}\:\d{1,2})\s*\-\s*(?P<timeTo>\d{1,2}\:\d{1,2})`)
	matched := re.FindStringSubmatch(stringPeriod)
	if matched == nil {
		return times, errors.New("no match")
	}
	return Times{
		Date:     matched[1],
		TimeFrom: matched[2],
		TimeTo:   matched[3],
	}, nil
}
