package handlers

import (
	"benny/src/fsm"
	"benny/src/models"
	"benny/src/repository"
	"benny/src/utils"
	"context"
	"fmt"
	tele "gopkg.in/telebot.v3"
	"log"
	"sort"
	"strconv"
	"time"
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
		barbers, missing := store.Barber().GetAll()
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

		stateManager, managerCloser := fsm.New(ctx, c.Chat().ID)
		defer managerCloser()

		_, missing := store.Customer().GetByTelegramId(c.Chat().ID)
		if missing == true {
			log.Printf("WARN: Пользователь: %d %s нажал кнопку записаться, но не был залогинен", c.Chat().ID, c.Chat().Username)
			return c.Send("wtf")
		}
		barberId := c.Callback().Data
		barber, missing := store.Barber().Get(barberId)
		if missing == true {
			return c.Send("Не найден такой барбер")
		}
		stateManager.Data().Set("barberId", c.Callback().Data)
		services, _ := store.Service().GetAll(barber.Id.String())
		buttons := make([]tele.Btn, len(services))
		for _, service := range services {
			var btn = BarberShiftsInlineKeyboard.Data(service.String(), "customerToService", service.Id.String())
			buttons = append(buttons, btn)
		}
		var rows = BarberShiftsInlineKeyboard.Split(1, buttons)
		BarberShiftsInlineKeyboard.Inline(rows...)
		var txt = fmt.Sprintf("<b>%s</b>\n\nВыбери услугу", barber.FullName)
		return c.Send(txt, BarberShiftsInlineKeyboard, tele.ModeHTML)
	}
}

func HandleSelectService() func(c tele.Context) error {
	return func(c tele.Context) error {
		ctx := context.Background()
		store, closer := repository.New(ctx)
		defer closer()

		stateManager, managerCloser := fsm.New(ctx, c.Chat().ID)
		defer managerCloser()
		stateManager.Data().Set("serviceId", c.Callback().Data)

		serviceId := c.Callback().Data
		service, _ := store.Service().Get(serviceId)

		barberId := stateManager.Data().Get("barberId")
		barber, missing := store.Barber().Get(barberId)
		if missing == true {
			return c.Send("Не найден такой барбер")
		}
		buttons := make([]tele.Btn, len(barber.Shifts))
		for _, service := range barber.Shifts {
			var btn = BarberShiftsInlineKeyboard.Data(service.String(), "customerToShift", service.Id.String())
			buttons = append(buttons, btn)
		}
		var rows = BarberShiftsInlineKeyboard.Split(1, buttons)
		BarberShiftsInlineKeyboard.Inline(rows...)
		var txt = fmt.Sprintf("<b>%s</b>\n\nУслуга: <b>%s</b>\n\nВыбери дату", barber.FullName, service.String())
		return c.Send(txt, BarberShiftsInlineKeyboard, tele.ModeHTML)
	}
}

func HandleSelectShift() func(c tele.Context) error {
	return func(c tele.Context) error {
		ctx := context.Background()
		store, closer := repository.New(ctx)
		defer closer()

		stateManager, managerCloser := fsm.New(ctx, c.Chat().ID)
		defer managerCloser()
		shiftId := c.Callback().Data
		stateManager.Data().Set("shiftId", shiftId)

		serviceId := stateManager.Data().Get("serviceId")
		service, _ := store.Service().Get(serviceId)
		_, missing := store.Customer().GetByTelegramId(c.Chat().ID)
		if missing == true {
			log.Printf("WARN: Пользователь: %startOfVisit %s нажал кнопку записаться, но не был залогинен", c.Chat().ID, c.Chat().Username)
			return c.Send("wtf")
		}
		shift, _ := store.Shift().Get(shiftId)

		closedHours := make(map[time.Time]time.Time)
		for _, visit := range shift.Visits {
			closedHours[visit.PlannedFrom] = visit.PlannedTo
		}
		openHours := make(map[time.Time]models.Visit)
		for startOfVisit := shift.PlannedFrom; startOfVisit.After(shift.PlannedTo) == false; startOfVisit = startOfVisit.Add(time.Duration(1) * time.Hour) {
			_, visitRegistered := closedHours[startOfVisit]
			if visitRegistered == false {
				timeOffset := time.Hour * time.Duration(shift.Barber.TimeZoneOffset)
				localEndOfVisit := startOfVisit.Add(time.Duration(service.Duration/60_000_000) * time.Minute).Add(timeOffset)
				localStartOfVisit := startOfVisit.Add(timeOffset)
				visit := models.Visit{PlannedFrom: localStartOfVisit, PlannedTo: localEndOfVisit}
				openHours[startOfVisit] = visit
			}
		}
		dateSortedVisits := make(utils.TimeSlice, 0, len(openHours))
		buttons := make([]tele.Btn, len(openHours))
		for _, visit := range openHours {
			dateSortedVisits = append(dateSortedVisits, visit)
		}
		sort.Sort(dateSortedVisits)
		for index, potentialVisit := range dateSortedVisits {
			stateManager.Data().Set(
				strconv.Itoa(index),
				fmt.Sprintf(
					"%s %s-%s",
					potentialVisit.PlannedFrom.Format("02.01.2006"),
					potentialVisit.PlannedFrom.Format("15:04"),
					potentialVisit.PlannedTo.Format("15:04"),
				),
			)
			var btn = BarberShiftsInlineKeyboard.Data(
				fmt.Sprintf("%s - %s", potentialVisit.PlannedFrom.Format("15:04"), potentialVisit.PlannedTo.Format("15:04")),
				"customerToTime", strconv.Itoa(index))
			buttons = append(buttons, btn)
		}
		var rows = BarberShiftsInlineKeyboard.Split(1, buttons)
		BarberShiftsInlineKeyboard.Inline(rows...)
		return c.Send(
			fmt.Sprintf("<b>%s</b>\n\nУслуга: <b>%s</b>\nДата: <b>%s</b>\n\nВыбери время", shift.Barber.FullName, service.String(), shift.PlannedFrom.Format("02.01.2006")),
			BarberShiftsInlineKeyboard, tele.ModeHTML)
	}
}

func HandleSelectTime() func(c tele.Context) error {
	return func(c tele.Context) error {
		ctx := context.Background()
		store, closer := repository.New(ctx)
		defer closer()

		stateManager, managerCloser := fsm.New(ctx, c.Chat().ID)
		defer managerCloser()
		timeId := c.Callback().Data
		timeStr := stateManager.Data().Get(timeId)
		stateManager.Data().Set("visitPeriod", timeStr)
		times, err := utils.ParseTimesFromString(timeStr)
		if err != nil {
			log.Fatal(err)
		}
		barberId := stateManager.Data().Get("barberId")
		barber, _ := store.Barber().Get(barberId)
		serviceId := stateManager.Data().Get("serviceId")
		service, _ := store.Service().Get(serviceId)
		CustomerShiftsInlineKeyboard.Inline(CustomerShiftsInlineKeyboard.Row(BtnDeclineVisit, BtnAcceptVisit))
		return c.Send(
			fmt.Sprintf("Барбер: <b>%s</b>\nУслуга: <b>%s</b>\nДата: <b>%s</b>\nВремя: <b>%s - %s</b>\n\nПодтвердить?",
				barber.FullName, service.String(), times.Date, times.TimeFrom, times.TimeTo),
			CustomerShiftsInlineKeyboard, tele.ModeHTML,
		)
	}
}
