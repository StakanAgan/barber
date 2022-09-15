package handlers

import (
	"benny/src/fsm"
	"benny/src/models"
	"benny/src/notifications"
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

func HandleReceivePhone() Handler {
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
		MainCustomerKeyboard.Inline(MainBarberKeyboard.Row(BtnCreateVisit))
		return c.Send(fmt.Sprintf("Велком, %s", customer.FullName), MainCustomerKeyboard)
	}
}

//
//func HandleStartCreateVisit() Handler {
//	return func(c tele.Context) error {
//		ctx := context.Background()
//		store, closer := repository.New(ctx)
//		defer closer()
//
//		_, missing := store.Customer().GetByTelegramId(c.Chat().ID)
//		if missing == true {
//			log.Printf("WARN: Пользователь: %d %s нажал кнопку записаться, но не был залогинен", c.Chat().ID, c.Chat().Username)
//			return c.Send("wtf")
//		}
//		barbers, missing := store.Barber().GetAll()
//		if missing == true {
//			return c.Send("Не нашлось парикмахеров")
//		}
//		buttons := make([]tele.Btn, len(barbers))
//		for _, barber := range barbers {
//			var btn = BarberShiftsInlineKeyboard.Data(fmt.Sprintf("%s", barber.FullName), "customerToBarber", barber.Id.String())
//			buttons = append(buttons, btn)
//		}
//		var rows = BarberShiftsInlineKeyboard.Split(1, buttons)
//		BarberShiftsInlineKeyboard.Inline(rows...)
//		return c.Send("Барберы", BarberShiftsInlineKeyboard)
//	}
//}

func HandleStartCreateVisit() Handler {
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
		//barberId := c.Callback().Data
		barber, missing := store.Barber().GetFirst()
		if missing == true {
			return c.Send("Не найден такой барбер")
		}
		stateManager.Data().Set("barberId", barber.Id.String())
		_, missing = store.Shift().GetActual(barber.Id.String())
		if missing == true {
			return c.Send("У Бена нет актуальных смен")
		}
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

func HandleSelectService() Handler {
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
		shifts, missing := store.Shift().GetActual(barberId)
		if missing == true {
			return c.Send("У Бена нет актуальных смен")
		}
		buttons := make([]tele.Btn, len(shifts))
		for _, shift := range shifts {
			var btn = BarberShiftsInlineKeyboard.Data(shift.String(), "customerToShift", shift.Id.String())
			buttons = append(buttons, btn)
		}
		var rows = BarberShiftsInlineKeyboard.Split(1, buttons)
		BarberShiftsInlineKeyboard.Inline(rows...)
		var txt = fmt.Sprintf("<b>%s</b>\n\nУслуга: <b>%s</b>\n\nВыбери дату", barber.FullName, service.String())
		return c.Send(txt, BarberShiftsInlineKeyboard, tele.ModeHTML)
	}
}

func HandleSelectShift() Handler {
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
		startOfOpenHours := shift.PlannedFrom
		if time.Now().UTC().After(startOfOpenHours) == true {
			d := 60 * time.Minute
			startOfOpenHours = time.Now().UTC().Round(d)
		}
		for startOfVisit := startOfOpenHours; startOfVisit.After(shift.PlannedTo) == false; startOfVisit = startOfVisit.Add(time.Duration(1) * time.Hour) {
			_, visitRegistered := closedHours[startOfVisit]
			if visitRegistered == false {
				timeOffset := time.Hour * time.Duration(shift.Barber.TimeZoneOffset)
				localEndOfVisit := startOfVisit.Add(time.Duration(service.Duration/60_000_000) * time.Minute).Add(timeOffset)
				localStartOfVisit := startOfVisit.Add(timeOffset)
				visit := models.Visit{PlannedFrom: localStartOfVisit, PlannedTo: localEndOfVisit}
				openHours[startOfVisit] = visit
			}
		}
		if len(openHours) == 0 {
			return c.Send("Не осталось свободных часов для записи, попробуй другую дату")
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

func HandleSelectTime() Handler {
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
			fmt.Sprintf("Барбер: <b>%s</b>\n<b>%s</b>\nЦена: <b>%d ₽</b>\nДата: <b>%s</b>\nВремя: <b>%s - %s</b>\n\nПодтвердить?",
				barber.FullName, service.Title, service.Price, times.Date, times.TimeFrom, times.TimeTo),
			CustomerShiftsInlineKeyboard, tele.ModeHTML,
		)
	}
}

func HandleAcceptVisit() Handler {
	return func(c tele.Context) error {
		ctx := context.Background()
		store, closer := repository.New(ctx)
		defer closer()

		stateManager, managerCloser := fsm.New(ctx, c.Chat().ID)
		defer managerCloser()
		defer stateManager.Data().Reset()
		defer c.Bot().EditReplyMarkup(c.Message(), nil)

		barberId := stateManager.Data().Get("barberId")
		serviceId := stateManager.Data().Get("serviceId")
		shiftId := stateManager.Data().Get("shiftId")
		visitPeriod := stateManager.Data().Get("visitPeriod")

		barber, _ := store.Barber().Get(barberId)
		service, _ := store.Service().Get(serviceId)
		shift, _ := store.Shift().Get(shiftId)
		visitTimes, _ := utils.ParseTimesFromString(visitPeriod)
		customer, _ := store.Customer().GetByTelegramId(c.Chat().ID)

		var dtParseFormat = "02.01.2006T15:04-07"
		var dtPassFormat = "%sT%s+0%d"
		plannedFrom, _ := time.Parse(dtParseFormat, fmt.Sprintf(dtPassFormat, visitTimes.Date, visitTimes.TimeFrom, barber.TimeZoneOffset))
		plannedTo, _ := time.Parse(dtParseFormat, fmt.Sprintf(dtPassFormat, visitTimes.Date, visitTimes.TimeTo, barber.TimeZoneOffset))

		var visit = &models.Visit{
			Service:       service,
			BarberShift:   shift,
			Customer:      customer,
			PlannedFrom:   plannedFrom,
			PlannedTo:     plannedTo,
			Price:         service.Price,
			DiscountPrice: 0,
			Status:        models.Created,
		}
		visit, err := store.Visit().Create(visit)
		if err != nil {
			MainCustomerKeyboard.Inline(MainCustomerKeyboard.Row(BtnCreateVisit))
			return c.Send("Кто-то записался раньше тебя. Попробуй на другое время", MainCustomerKeyboard)
		}
		err = notifications.NotifyBarberAboutCreated(c.Bot(), barber.TelegramId, *visit)
		if err != nil {
			c.Send("Бен пока не получил уведомление, но зайдет и прочитает")
			log.Println("WARN: Not found Benny telegramId")
		}
		MainCustomerKeyboard.Inline(MainCustomerKeyboard.Row(BtnCreateVisit))

		return c.Send("Записано", MainCustomerKeyboard)
	}
}

func HandleDeclineVisit() Handler {
	return func(c tele.Context) error {
		ctx := context.Background()

		stateManager, managerCloser := fsm.New(ctx, c.Chat().ID)
		defer managerCloser()
		defer stateManager.Data().Reset()
		defer c.Bot().EditReplyMarkup(c.Message(), nil)
		MainCustomerKeyboard.Inline(MainCustomerKeyboard.Row(BtnCreateVisit))

		return c.Send("Запись отменена", MainCustomerKeyboard)
	}
}
