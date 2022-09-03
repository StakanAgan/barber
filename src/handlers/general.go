package handlers

import (
	"benny/src/repository"
	"context"
	"fmt"
	tele "gopkg.in/telebot.v3"
	"log"
)

func HandleStart() func(c tele.Context) error {
	return func(c tele.Context) error {
		ctx := context.Background()
		store, closer := repository.New(ctx)
		defer closer()

		var barber, missing = store.Barber().GetByTelegramId(uint64(c.Chat().ID))
		if missing == true {
			log.Printf("INFO: User %d try to Start bot", uint64(c.Chat().ID))
			customer, missing := store.Customer().GetByTelegramId(c.Chat().ID)
			if missing == true {
				PhoneRequestKeyboard.Reply(PhoneRequestKeyboard.Row(BtnRequestPhone))
				return c.Send("Заделись цифрами, чтобы записаться на стригу. Просто нажми на <b>☎️ Поделиться цифрами</b> внизу 👇🏼", PhoneRequestKeyboard, tele.ModeHTML)
			}
			MainCustomerKeyboard.Reply(MainCustomerKeyboard.Row(BtnCreateVisit))
			return c.Send(fmt.Sprintf("Йо, тебя зовут %s, твой id %s", customer.FullName, customer.Id), MainCustomerKeyboard)
		}
		MainBarberKeyboard.Reply(MainBarberKeyboard.Row(BtnShifts))
		return c.Send(fmt.Sprintf("Йо, твой тлф %s", barber.Phone), MainBarberKeyboard)
	}
}