package services

import (
	"benny/src/models"
	"benny/src/repository"
	"context"
	"fmt"
	tele "gopkg.in/telebot.v3"
	"log"
	"time"
)

func NotifyBarberAboutCreated(b *tele.Bot, barberTelegramId int64, visit models.Visit) error {
	barberTg := &tele.User{ID: barberTelegramId}
	_, err := b.Send(barberTg, fmt.Sprintf("К тебе записались\n\n"+
		"%s %d ₽\n<b>%s</b>\n%s +%s", visit.Service.Title, visit.Price-visit.DiscountPrice,
		visit.PlannedFrom.Format("02.01.2006 15:04"),
		visit.Customer.FullName, visit.Customer.Phone), tele.ModeHTML)
	return err
}

func CreateNewBarberShiftOnNextWeek(b tele.Bot) {
	for range time.Tick(time.Hour) {
		now := time.Now().UTC()
		if now.Hour() != 20 {
			log.Println("INFO: Not now, later.")
			continue
		}
		ctx := context.Background()
		store, closer := repository.New(ctx)
		barber, missing := store.Barber().GetFirst()
		if missing == true {
			log.Println("WARN: No one barber")
			closer()
			continue
		}
		todayShift, missing := store.Shift().GetToday(barber.Id.String())
		if missing == true {
			log.Println("INFO: Today without shift")
			closer()
			continue
		}
		newShift := &models.BarberShift{
			PlannedFrom: todayShift.PlannedFrom.UTC().AddDate(0, 0, 7),
			PlannedTo:   todayShift.PlannedTo.UTC().AddDate(0, 0, 7),
		}
		newShift, err := store.Shift().Create(barber.Id.String(), newShift)
		if err != nil {
			log.Println("INFO: Shift already created on next week")
			closer()
			continue
		}
		barberTg := &tele.User{ID: barber.TelegramId}
		_, err = b.Send(barberTg, fmt.Sprintf("Создана смена\n\n%s %s - %s",
			newShift.PlannedFrom.Add(barber.TimeOffset()).Format("02.01.2006"),
			newShift.PlannedFrom.Add(barber.TimeOffset()).Format("15:04"),
			newShift.PlannedTo.Add(barber.TimeOffset()).Format("15:04"),
		))
		if err != nil {
			log.Println("WARN: No notify about new shift barber")
		}
		closer()
	}

}
