package handlers

import (
	"benny/src/repository"
	"fmt"
	tele "gopkg.in/telebot.v3"
	"log"
)

type Handler func(c tele.Context) error

func HandleStart(store *repository.Store) Handler {
	return func(c tele.Context) error {
		barber, err := store.Barber().GetByTelegramId(uint64(c.Chat().ID))
		if err != nil {
			return c.Send("Какая-то ошибка...")
		}
		if barber.Missing() {
			log.Printf("INFO: User %d try to Start bot", uint64(c.Chat().ID))
			customer, err := store.Customer().GetByTelegramId(c.Chat().ID)
			if err != nil {
				return c.Send("Какая-то ошибка...")
			}
			if customer.Missing() {
				PhoneRequestKeyboard.Reply(PhoneRequestKeyboard.Row(BtnRequestPhone))
				return c.Send("Заделись цифрами, чтобы записаться на стригу. Просто нажми на <b>☎️ Поделиться цифрами</b> внизу 👇🏼", PhoneRequestKeyboard, tele.ModeHTML)
			}
			MainCustomerKeyboard.Inline(MainCustomerKeyboard.Row(BtnCreateVisit))
			return c.Send(fmt.Sprintf("Велком, %s", customer.FullName), MainCustomerKeyboard)
		}
		MainBarberKeyboard.Reply(MainBarberKeyboard.Row(BtnShifts), MainBarberKeyboard.Row(BtnServices))
		return c.Send(fmt.Sprintf("Йо, твой тлф %s", barber.Phone), MainBarberKeyboard)
	}
}
