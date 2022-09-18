package services

import (
	"benny/src/models"
	"benny/src/repository"
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

func NotifyCustomerAboutCancel(b *tele.Bot, barber models.Barber, visit models.Visit) error {
	customerTg := &tele.User{ID: visit.Customer.TelegramId}
	totalPrice, _ := visit.TotalPrice.Get()
	_, err := b.Send(customerTg, fmt.Sprintf("У %s отменилась смена, поэтому отменилась запись\n\n"+
		"%s %d ₽\n<b>%s</b>", barber.FullName, visit.Service.Title, totalPrice,
		visit.PlannedFrom.Add(barber.TimeOffset()).Format("02.01.2006 15:04")), tele.ModeHTML)
	return err
}

func CreateNewBarberShiftOnNextWeek(b *tele.Bot, store *repository.Store) {
	for range time.Tick(time.Hour) {
		now := time.Now().UTC()
		if now.Hour() != 20 {
			log.Println("INFO: Not now, later.")
			continue
		}

		barber, err := store.Barber().GetFirst()
		if barber.Missing() {
			log.Println("WARN: No one barber")
			continue
		}
		if err != nil {
			log.Printf("ERROR: Error while get barber: %s", err)
		}
		todayShift, err := store.Shift().GetLast(barber.Id.String())
		if err != nil {
			log.Printf("ERROR: Error while get today shift: %s", err)
			continue
		}
		if todayShift.Missing() {
			log.Println("INFO: Today without shift")
			continue
		}
		newShift := &models.BarberShift{
			PlannedFrom: todayShift.PlannedFrom.UTC().AddDate(0, 0, 7),
			PlannedTo:   todayShift.PlannedTo.UTC().AddDate(0, 0, 7),
		}
		newShift, err = store.Shift().Create(barber.Id.String(), newShift)
		if err != nil {
			log.Printf("INFO: Shift already created on next week: %s", err)
			continue
		}
		barberTg := &tele.User{ID: barber.TelegramId}
		_, err = b.Send(barberTg, fmt.Sprintf("Создана смена\n\n%s %s - %s",
			newShift.PlannedFrom.Add(barber.TimeOffset()).Format("02.01.2006"),
			newShift.PlannedFrom.Add(barber.TimeOffset()).Format("15:04"),
			newShift.PlannedTo.Add(barber.TimeOffset()).Format("15:04"),
		))
		if err != nil {
			log.Printf("WARN: No notify about new shift barber: %s", err)
		}
	}

}
