package handlers

import (
	"benny/fsm"
	"benny/models"
	"benny/store"
	"benny/utils"
	"context"
	"fmt"
	"github.com/edgedb/edgedb-go"
	tele "gopkg.in/telebot.v3"
	"log"
	"time"
)

func HandleMainShifts() func(c tele.Context) error {
	return func(c tele.Context) error {
		ctx := context.Background()
		store, closer := store.New(ctx)
		defer closer()

		barber, missing := store.Barber().GetByTelegramId(uint64(c.Chat().ID))
		if missing == true {
			log.Printf("INFO: Customer %d try to get shifts", uint64(c.Chat().ID))
		}
		shifts, _ := store.Shift().GetActual(barber.Id)
		buttons := make([]tele.Btn, len(shifts)+2) // количество смен + кнопка для создания смен + кнопка все смены
		buttons = append(buttons, BtnAllShifts)
		for _, shift := range shifts {
			var btn = BarberShiftsInlineKeyboard.Data(shift.String(), "toShift", shift.Id.String())
			buttons = append(buttons, btn)
		}
		buttons = append(buttons, BtnCreateShift)
		var rows = BarberShiftsInlineKeyboard.Split(1, buttons)
		BarberShiftsInlineKeyboard.Inline(rows...)
		return c.Send("Твои смены:", BarberShiftsInlineKeyboard)
	}
}

func HandleAllShifts() func(c tele.Context) error {
	return func(c tele.Context) error {
		ctx := context.Background()
		store, closer := store.New(ctx)
		defer closer()

		barber, missing := store.Barber().GetByTelegramId(uint64(c.Chat().ID))
		if missing == true {
			log.Printf("INFO: Customer %d try to get shifts", uint64(c.Chat().ID))
		}
		shifts, _ := store.Shift().GetAll(barber.Id)
		buttons := make([]tele.Btn, len(shifts)+2) // количество смен + кнопка для создания смен + кнопка все смены
		buttons = append(buttons, BtnPlannedShifts)
		for _, shift := range shifts {
			var btn = BarberShiftsInlineKeyboard.Data(shift.String(), "toShift", shift.Id.String())
			buttons = append(buttons, btn)
		}
		buttons = append(buttons, BtnCreateShift)
		var rows = BarberShiftsInlineKeyboard.Split(1, buttons)
		BarberShiftsInlineKeyboard.Inline(rows...)
		return c.Send("Твои смены:", BarberShiftsInlineKeyboard)
	}
}

func HandleText() func(c tele.Context) error {
	return func(c tele.Context) error {
		ctx := context.Background()
		store, closer := store.New(ctx)
		defer closer()

		stateManager, managerCloser := fsm.New(ctx, c.Chat().ID)
		defer managerCloser()

		barber, missing := store.Barber().GetByTelegramId(uint64(c.Chat().ID))
		if missing == true {
			return c.Send("ты кто?")
		}
		state := stateManager.State().Get()

		switch state {
		case fsm.ShiftEnter:
			times, err := utils.ParseTimesFromString(c.Text())
			if err != nil {
				return c.Send("Не тот формат, брат. Введи период в формате\n<b>29.08.2022 11:00-19:00</b>", tele.ModeHTML)
			}
			var dtParseFormat = "02.01.2006T15:04-07"
			var dtPassFormat = "%sT%s+0%d"
			plannedFrom, err := time.Parse(dtParseFormat, fmt.Sprintf(dtPassFormat, times.Date, times.TimeFrom, barber.TimeZoneOffset))
			if err != nil {
				log.Println(err)
				return c.Send(fmt.Sprintf("Ошибка где-то здесь <b>%s</b> <b>%s</b>-%s", times.Date, times.TimeFrom, times.TimeTo), tele.ModeHTML)
			}
			plannedTo, err := time.Parse(dtParseFormat, fmt.Sprintf(dtPassFormat, times.Date, times.TimeTo, barber.TimeZoneOffset))
			if err != nil {
				log.Println(err)
				return c.Send(fmt.Sprintf("Ошибка где-то здесь <b>%s</b> %s-<b>%s</b>", times.Date, times.TimeFrom, times.TimeTo), tele.ModeHTML)
			}
			if plannedFrom.Unix() >= plannedTo.Unix() {
				return c.Send(fmt.Sprintf("Время указанное от <b>%s</b> позже либо равно времени до <b>%s</b>\n"+
					"Начало смены должно быть раньше ее конца", times.TimeFrom, times.TimeTo), tele.ModeHTML)
			}

			var shift = &models.BarberShift{PlannedFrom: plannedFrom.UTC(), Barber: *barber, PlannedTo: plannedTo.UTC()}
			shift, err = store.Shift().Create(barber.Id, shift)
			if err != nil {
				return c.Send("Брат, эта смена пересекается с другой\nПопробуй еще раз, но так, чтобы не пересекалась", tele.ModeHTML)
			}
			err = c.Send("Добавили смену работяге")
			if err != nil {
				log.Fatal(err)
			}
			stateManager.State().Reset()
			return HandleMainShifts()(c)
		default:
			return c.Send("Я тебя не понял")
		}
	}
}

func HandleStartCreateShift() func(c tele.Context) error {
	return func(c tele.Context) error {
		ctx := context.Background()
		store, closer := store.New(ctx)
		defer closer()

		stateManager, managerCloser := fsm.New(ctx, c.Chat().ID)
		defer managerCloser()

		_, missing := store.Barber().GetByTelegramId(uint64(c.Chat().ID))
		if missing == true {
			log.Println("INFO: Чел как-то нажал кнопку создать смену")
			return c.Send("ты кто?")
		}
		tgerr := c.Send("Напиши дату смены и со скольких до скольких в формате:\n<b>29.08.2022 11:00-19:00</b>", tele.ModeHTML)
		if tgerr != nil {
			log.Fatal(tgerr)
		}
		err := stateManager.State().Set(fsm.ShiftEnter)
		if err != nil {
			log.Fatal(err)
		}
		return tgerr
	}
}

func HandleGetShift() func(c tele.Context) error {
	return func(c tele.Context) error {
		ctx := context.Background()
		store, closer := store.New(ctx)
		defer closer()

		_, missing := store.Barber().GetByTelegramId(uint64(c.Chat().ID))
		if missing == true {
			log.Println("INFO: Чел как-то нажал не на свою смену")
			return c.Send("ты кто?")
		}
		shiftId := &edgedb.UUID{}
		err := shiftId.UnmarshalText([]byte(c.Callback().Data))
		if err != nil {
			log.Fatal(err)
		}
		shift, missing := store.Shift().Get(*shiftId)
		if missing == true {
			return c.Send("Не та смена")
		}
		var (
			needBtn   bool
			btnAction tele.Btn
		)
		switch shift.Status {
		case string(models.Planned):
			btnAction = BarberShiftsInlineKeyboard.Data("Начать смену", "start", shiftId.String())
			needBtn = true
		case string(models.Work):
			btnAction = BarberShiftsInlineKeyboard.Data("Завершить смену", "finish", shiftId.String())
			needBtn = true
		default:
			needBtn = false
		}
		if needBtn == true {
			BarberShiftsInlineKeyboard.Inline(BarberShiftsInlineKeyboard.Row(BtnCancelShift, btnAction))
		}
		var txt = fmt.Sprintf("<b>%s</b>\n\nСтатус: %s", shift.String(), shift.Status)
		return c.Send(txt, BarberShiftsInlineKeyboard, tele.ModeHTML)
	}
}

func HandleStartShift() func(c tele.Context) error {
	return func(c tele.Context) error {
		ctx := context.Background()
		store, closer := store.New(ctx)
		defer closer()

		_, missing := store.Barber().GetByTelegramId(uint64(c.Chat().ID))
		if missing == true {
			log.Println("INFO: Чел как-то нажал не на свою смену")
			return c.Send("ты кто?")
		}
		shiftId := &edgedb.UUID{}
		err := shiftId.UnmarshalText([]byte(c.Callback().Data))
		if err != nil {
			log.Fatal(err)
		}
		_, missing = store.Shift().UpdateStatus(*shiftId, models.Work)
		return HandleGetShift()(c)
	}
}

func HandleFinishShift() func(c tele.Context) error {
	return func(c tele.Context) error {
		ctx := context.Background()
		store, closer := store.New(ctx)
		defer closer()

		_, missing := store.Barber().GetByTelegramId(uint64(c.Chat().ID))
		if missing == true {
			log.Println("INFO: Чел как-то нажал не на свою смену")
			return c.Send("ты кто?")
		}
		shiftId := &edgedb.UUID{}
		err := shiftId.UnmarshalText([]byte(c.Callback().Data))
		if err != nil {
			log.Fatal(err)
		}
		_, missing = store.Shift().UpdateStatus(*shiftId, models.Finished)
		return HandleGetShift()(c)
	}
}
