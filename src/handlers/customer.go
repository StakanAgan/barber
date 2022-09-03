package handlers

import (
	"benny/src/models"
	"benny/src/repository"
	"context"
	"fmt"
	"github.com/edgedb/edgedb-go"
	tele "gopkg.in/telebot.v3"
	"log"
)

func HandleReceivePhone() func(c tele.Context) error {
	return func(c tele.Context) error {
		ctx := context.Background()
		store, closer := repository.New(ctx)
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
		store, closer := repository.New(ctx)
		defer closer()

		_, missing := store.Customer().GetByTelegramId(c.Chat().ID)
		if missing == true {
			log.Printf("WARN: Пользователь: %d %s нажал кнопку записаться, но не был залогинен", c.Chat().ID, c.Chat().Username)
			return c.Send("wtf")
		}
		barbers, missing := store.Barber().GetAllWithShifts()
		if missing == true {
			return c.Send("Не нашлось парикмахеров")
		}
		buttons := make([]tele.Btn, len(barbers))
		for _, barber := range barbers {
			var btn = BarberShiftsInlineKeyboard.Data(fmt.Sprintf("%s", barber.FullName), "customerToBarber", barber.Id.String())
			buttons = append(buttons, btn)
		}
		var rows = BarberShiftsInlineKeyboard.Split(1, buttons)
		BarberShiftsInlineKeyboard.Inline(rows...)
		return c.Send("Барберы", BarberShiftsInlineKeyboard)
	}
}

func HandleSelectBarber() func(c tele.Context) error {
	return func(c tele.Context) error {
		ctx := context.Background()
		store, closer := repository.New(ctx)
		defer closer()

		_, missing := store.Customer().GetByTelegramId(c.Chat().ID)
		if missing == true {
			log.Printf("WARN: Пользователь: %d %s нажал кнопку записаться, но не был залогинен", c.Chat().ID, c.Chat().Username)
			return c.Send("wtf")
		}
		shiftId := &edgedb.UUID{}
		err := shiftId.UnmarshalText([]byte(c.Callback().Data))
		if err != nil {
			log.Fatal(err)
		}
		barber, missing := store.Barber().GetWithShifts(*shiftId)
		if missing == true {
			return c.Send("Не найден такой барбер")
		}
		buttons := make([]tele.Btn, len(barber.Shifts))
		for _, shift := range barber.Shifts {
			var btn = BarberShiftsInlineKeyboard.Data(shift.String(), "customerToBarber", shift.Id.String())
			buttons = append(buttons, btn)
		}
		var rows = BarberShiftsInlineKeyboard.Split(1, buttons)
		BarberShiftsInlineKeyboard.Inline(rows...)
		var txt = fmt.Sprintf("<b>%s</b>\n\nТелефон: +%s", barber.FullName, barber.Phone)
		return c.Send(txt, BarberShiftsInlineKeyboard, tele.ModeHTML)
	}
}
