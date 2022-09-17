package handlers

import (
	"benny/src/fsm"
	"benny/src/models"
	"benny/src/notifications"
	"benny/src/repository"
	"benny/src/utils"
	"fmt"
	tele "gopkg.in/telebot.v3"
	"log"
	"sort"
	"strconv"
	"time"
)

func HandleReceivePhone(store *repository.Store) Handler {
	return func(c tele.Context) error {
		var customer = &models.Customer{
			Phone:      c.Message().Contact.PhoneNumber,
			FullName:   fmt.Sprintf("%s %s", c.Message().Contact.LastName, c.Message().Contact.FirstName),
			TelegramId: uint64(c.Chat().ID),
		}
		customer, err := store.Customer().Create(customer)
		if err != nil {
			return c.Send("Не удалось залогиниться. Попробуй написать @ctxkn")
		}
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

func HandleStartCreateVisit(store *repository.Store, stateManager *fsm.StateManager) Handler {
	return func(c tele.Context) error {
		customer, err := store.Customer().GetByTelegramId(c.Chat().ID)
		if err != nil {
			return c.Send("Какая-то ошибка...")
		}
		if customer.Missing() {
			log.Printf("WARN: Пользователь: %d %s нажал кнопку записаться, но не был залогинен", c.Chat().ID, c.Chat().Username)
			return c.Send("wtf")
		}
		//barberId := c.Callback().Data
		barber, err := store.Barber().GetFirst()
		if err != nil {
			return c.Send("Какая-то ошибка...")
		}
		if barber.Missing() {
			return c.Send("Не найден такой барбер")
		}
		stateManager.Data(c.Chat().ID).Set("barberId", barber.Id.String())
		shifts, err := store.Shift().GetActual(barber.Id.String())
		if err != nil {
			return c.Send("Какая-то ошибка, попробуй позже")
		}
		if len(shifts) == 0 {
			return c.Send("У Бена нет актуальных смен")
		}
		services, err := store.Service().GetAll(barber.Id.String())
		if err != nil {
			return c.Send("Какая-то ошибка, попробуй позже")
		}
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

func HandleSelectService(store *repository.Store, stateManager *fsm.StateManager) Handler {
	return func(c tele.Context) error {
		stateManager.Data(c.Chat().ID).Set("serviceId", c.Callback().Data)

		serviceId := c.Callback().Data
		service, err := store.Service().Get(serviceId)
		if err != nil {
			return c.Send("Какая-то ошибка, попробуй позже")
		}

		barberId := stateManager.Data(c.Chat().ID).Get("barberId")
		barber, err := store.Barber().Get(barberId)
		if err != nil {
			return c.Send("Какая-то ошибка...")
		}
		if barber.Missing() {
			return c.Send("Не найден такой барбер")
		}
		shifts, err := store.Shift().GetActual(barberId)
		if err != nil {
			return c.Send("Какая-то ошибка, попробуй позже")
		}
		if len(shifts) == 0 {
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

func HandleSelectShift(store *repository.Store, stateManager *fsm.StateManager) Handler {
	return func(c tele.Context) error {
		customerTgId := &c.Chat().ID
		shiftId := c.Callback().Data
		stateManager.Data(*customerTgId).Set("shiftId", shiftId)

		serviceId := stateManager.Data(*customerTgId).Get("serviceId")
		service, err := store.Service().Get(serviceId)
		if err != nil {
			return c.Send("Какая-то ошибка, попробуй позже")
		}
		customer, err := store.Customer().GetByTelegramId(c.Chat().ID)
		if err != nil {
			return c.Send("Какая-то ошибка...")
		}
		if customer.Missing() {
			log.Printf("WARN: Пользователь: %startOfVisit %s нажал кнопку записаться, но не был залогинен", c.Chat().ID, c.Chat().Username)
			return c.Send("Какая-то ошибка...")
		}
		shift, err := store.Shift().Get(shiftId)
		if err != nil {
			return c.Send("Какая-то ошибка")
		}
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
			stateManager.Data(*customerTgId).Set(
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

func HandleSelectTime(store *repository.Store, stateManager *fsm.StateManager) Handler {
	return func(c tele.Context) error {
		customerTgId := c.Chat().ID
		timeId := c.Callback().Data
		timeStr := stateManager.Data(customerTgId).Get(timeId)
		stateManager.Data(customerTgId).Set("visitPeriod", timeStr)
		times, err := utils.ParseTimesFromString(timeStr)
		if err != nil {
			log.Fatal(err)
		}
		barberId := stateManager.Data(customerTgId).Get("barberId")
		barber, err := store.Barber().Get(barberId)
		if err != nil {
			return c.Send("Какая-то ошибка...")
		}
		serviceId := stateManager.Data(customerTgId).Get("serviceId")
		service, err := store.Service().Get(serviceId)
		if err != nil {
			return c.Send("Какая-то ошибка, попробуй позже")
		}
		CustomerShiftsInlineKeyboard.Inline(CustomerShiftsInlineKeyboard.Row(BtnDeclineVisit, BtnAcceptVisit))
		return c.Send(
			fmt.Sprintf("Барбер: <b>%s</b>\n<b>%s</b>\nЦена: <b>%d ₽</b>\nДата: <b>%s</b>\nВремя: <b>%s - %s</b>\n\nПодтвердить?",
				barber.FullName, service.Title, service.Price, times.Date, times.TimeFrom, times.TimeTo),
			CustomerShiftsInlineKeyboard, tele.ModeHTML,
		)
	}
}

func HandleAcceptVisit(store *repository.Store, stateManager *fsm.StateManager) Handler {
	return func(c tele.Context) error {
		customerTgId := c.Chat().ID
		defer stateManager.Data(customerTgId).Reset()
		defer c.Bot().EditReplyMarkup(c.Message(), nil)

		barberId := stateManager.Data(customerTgId).Get("barberId")
		serviceId := stateManager.Data(customerTgId).Get("serviceId")
		shiftId := stateManager.Data(customerTgId).Get("shiftId")
		visitPeriod := stateManager.Data(customerTgId).Get("visitPeriod")

		barber, err := store.Barber().Get(barberId)
		if err != nil {
			return c.Send("Какая-то ошибка...")
		}
		service, err := store.Service().Get(serviceId)
		if err != nil {
			return c.Send("Какая-то ошибка, попробуй позже")
		}
		shift, err := store.Shift().Get(shiftId)
		if err != nil {
			return c.Send("Какая-то ошибка...")
		}
		visitTimes, _ := utils.ParseTimesFromString(visitPeriod)
		customer, err := store.Customer().GetByTelegramId(c.Chat().ID)
		if err != nil {
			return c.Send("Какая-то ошибка...")
		}

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
		visit, err = store.Visit().Create(visit)
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

func HandleDeclineVisit(stateManager *fsm.StateManager) Handler {
	return func(c tele.Context) error {
		defer stateManager.Data(c.Chat().ID).Reset()
		defer c.Bot().EditReplyMarkup(c.Message(), nil)
		MainCustomerKeyboard.Inline(MainCustomerKeyboard.Row(BtnCreateVisit))

		return c.Send("Запись отменена", MainCustomerKeyboard)
	}
}
