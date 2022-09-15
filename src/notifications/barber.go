package notifications

import (
	"benny/src/models"
	"fmt"
	tele "gopkg.in/telebot.v3"
)

func NotifyBarberAboutCreated(b *tele.Bot, barberTelegramId int64, visit models.Visit) error {
	barberTg := &tele.User{ID: barberTelegramId}
	_, err := b.Send(barberTg, fmt.Sprintf("К тебе записались\n\n"+
		"%s %d ₽\n<b>%s</b>\n%s +%s", visit.Service.Title, visit.Price-visit.DiscountPrice,
		visit.PlannedFrom.Format("02.01.2006 15:04"),
		visit.Customer.FullName, visit.Customer.Phone), tele.ModeHTML)
	return err
}
