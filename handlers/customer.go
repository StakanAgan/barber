package handlers

import (
	"benny/models"
	"benny/store"
	"context"
	"fmt"
	tele "gopkg.in/telebot.v3"
	"log"
)

func HandleReceivePhone() func(c tele.Context) error {
	return func(c tele.Context) error {
		ctx := context.Background()
		store, closer := store.New(ctx)
		defer closer()

		var customer = &models.Customer{
			Phone:      c.Message().Contact.PhoneNumber,
			FullName:   fmt.Sprintf("%s %s", c.Message().Contact.LastName, c.Message().Contact.FirstName),
			TelegramId: uint64(c.Chat().ID),
		}
		customer = store.Customer().Create(customer)
		MainCustomerKeyboard.Reply(MainBarberKeyboard.Row(BtnCreateVisit))
		return c.Send(fmt.Sprintf("%s", customer.Id), MainCustomerKeyboard)
	}
}

func HandleStartCreateVisit() func(c tele.Context) error {
	return func(c tele.Context) error {
		ctx := context.Background()
		store, closer := store.New(ctx)
		defer closer()

		customer, missing := store.Customer().GetByTelegramId(c.Chat().ID)
		if missing == true {
			log.Printf("WARN: Пользователь: %d %s нажал кнопку записаться, но не был залогинен", c.Chat().ID, c.Chat().Username)
			c.Send("wtf")
		}
		barbers, missing := store.Barber().GetAll()
		if missing == true {
			c.Send("Не нашлось парикмахеров")
		}
		buttons := make([]tele.Btn, len(barbers))
		buttons = append(buttons, BtnPlannedShifts)
		for _, barber := range barbers {
			var btn = BarberShiftsInlineKeyboard.Data(, "toShift", barber.Id.String())
			buttons = append(buttons, btn)
		}
		buttons = append(buttons, BtnCreateShift)
		var rows = BarberShiftsInlineKeyboard.Split(1, buttons)
		BarberShiftsInlineKeyboard.Inline(rows...)
	}
}
