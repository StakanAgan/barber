package utils

import (
	"benny/src/models"
	"errors"
	tele "gopkg.in/telebot.v3"
	"log"
	"reflect"
	"regexp"
	"runtime"
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

type TimeSlice []models.Visit

func (p TimeSlice) Len() int {
	return len(p)
}

func (p TimeSlice) Less(i, j int) bool {
	return p[i].PlannedFrom.Before(p[j].PlannedFrom)
}

func (p TimeSlice) Swap(i, j int) {
	p[i], p[j] = p[j], p[i]
}

func LogUpdates() tele.MiddlewareFunc {
	return func(next tele.HandlerFunc) tele.HandlerFunc {
		return func(c tele.Context) error {
			if c.Callback() != nil {
				log.Printf("INFO: userId: %d, userName: %s, callback: %s %s",
					c.Chat().ID, c.Chat().Username, c.Callback().Unique, c.Callback().Data)
			}
			if c.Text() != "" {
				log.Printf("INFO: userId: %d, userName: %s, text: %s", c.Chat().ID, c.Chat().Username, c.Text())
			}
			v := runtime.FuncForPC(reflect.ValueOf(next).Pointer()).Name()
			log.Printf("INFO: handler %s", v)
			return next(c)
		}
	}
}
