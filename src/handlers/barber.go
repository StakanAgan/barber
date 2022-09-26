package handlers

import (
	"benny/src/fsm"
	"benny/src/models"
	"benny/src/repository"
	"benny/src/services"
	"benny/src/utils"
	"fmt"
	"github.com/edgedb/edgedb-go"
	tele "gopkg.in/telebot.v3"
	"log"
	"sort"
	"strconv"
	"time"
)

func HandleMainShifts(store *repository.Store) Handler {
	return func(c tele.Context) error {
		barber, err := store.Barber().GetByTelegramId(uint64(c.Chat().ID))
		if err != nil {
			return c.Send("Какая-то ошибка...")
		}
		if barber.Missing() {
			log.Printf("INFO: Customer %d try to get shifts", uint64(c.Chat().ID))
			return c.Send("Тебя не нашлось в списке")
		}
		shifts, err := store.Shift().GetActual(barber.Id.String())
		if err != nil {
			return c.Send("Не получится пока увидеть смены брат")
		}
		buttons := make([]tele.Btn, len(shifts)+2) // количество смен + кнопка для создания смен + кнопка все смены
		buttons = append(buttons, BtnAllShifts)
		sort.Slice(shifts, func(i, j int) bool {
			return shifts[i].PlannedFrom.Before(shifts[j].PlannedFrom)
		})
		for _, shift := range shifts {
			var btn = BarberShiftsInlineKeyboard.Data(shift.String(), "barberToShift", shift.Id.String())
			buttons = append(buttons, btn)
		}
		buttons = append(buttons, BtnCreateShift)
		var rows = BarberShiftsInlineKeyboard.Split(1, buttons)
		BarberShiftsInlineKeyboard.Inline(rows...)
		if c.Callback() != nil {
			return c.Edit("Твои смены:", BarberShiftsInlineKeyboard)
		}
		return c.Send("Твои смены:", BarberShiftsInlineKeyboard)
	}
}

func HandleAllShifts(store *repository.Store) Handler {
	return func(c tele.Context) error {
		barber, err := store.Barber().GetByTelegramId(uint64(c.Chat().ID))
		if err != nil {
			return c.Send("Какая-то ошибка...")
		}
		if barber.Missing() {
			log.Printf("INFO: Customer %d try to get shifts", uint64(c.Chat().ID))
		}
		shifts, err := store.Shift().GetAll(barber.Id.String())
		if err != nil {
			return c.Send("Возникла проблемка, не получилось показать все смены")
		}
		buttons := make([]tele.Btn, len(shifts)+2) // количество смен + кнопка для создания смен + кнопка все смены
		buttons = append(buttons, BtnPlannedShifts)
		sort.Slice(shifts, func(i, j int) bool {
			return shifts[i].PlannedFrom.Before(shifts[j].PlannedFrom)
		})
		for _, shift := range shifts {
			var btn = BarberShiftsInlineKeyboard.Data(shift.String(), "barberToShift", shift.Id.String())
			buttons = append(buttons, btn)
		}
		buttons = append(buttons, BtnCreateShift)
		var rows = BarberShiftsInlineKeyboard.Split(1, buttons)
		BarberShiftsInlineKeyboard.Inline(rows...)
		return c.Edit("Твои смены:", BarberShiftsInlineKeyboard)
	}
}

func HandleText(store *repository.Store, stateManager *fsm.StateManager) Handler {
	return func(c tele.Context) error {
		barber, err := store.Barber().GetByTelegramId(uint64(c.Chat().ID))
		if err != nil {
			return c.Send("Какая-то ошибка...")
		}
		if barber.Missing() {
			return c.Send("ты кто?")
		}
		state := stateManager.State(c.Chat().ID).Get()
		switch state {
		case fsm.ShiftEnter:
			err := HandleShiftEnter(store, stateManager, c, barber)
			if err != nil {
				return err
			}
			return HandleMainShifts(store)(c)
		case fsm.ServiceEnterTitle:
			return HandleServiceEnterTitle(stateManager, c)
		case fsm.ServiceEnterPrice:
			return HandleServiceEnterPrice(stateManager, c)
		case fsm.ServiceEnterDuration:
			err := HandleEndCreateService(store, stateManager, *barber, c)
			if err != nil {
				return err
			}
			stateManager.Data(c.Chat().ID).Reset()
			return HandleMainServices(store)(c)
		default:
			return c.Send("Я тебя не понял")
		}
	}
}

func HandleShiftEnter(store *repository.Store, stateManager *fsm.StateManager, c tele.Context, barber *models.Barber) error {
	times, err := utils.ParseTimesFromString(c.Text())
	now := time.Now()

	if err != nil {
		return c.Send(fmt.Sprintf("Не тот формат, брат. Введи период в формате\n<b>%s 11:00-19:00</b>", now.Format("02.01.2006")), tele.ModeHTML)
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
	shift, err = store.Shift().Create(barber.Id.String(), shift)
	if err != nil {
		return c.Send("Брат, эта смена пересекается с другой\nПопробуй еще раз, но так, чтобы не пересекалась", tele.ModeHTML)
	}
	err = c.Send("Добавили смену работяге")
	if err != nil {
		log.Fatal(err)
	}
	stateManager.State(c.Chat().ID).Reset()
	return err
}

func HandleStartCreateShift(store *repository.Store, stateManager *fsm.StateManager) Handler {
	return func(c tele.Context) error {
		customerTgId := &c.Chat().ID
		barber, err := store.Barber().GetByTelegramId(uint64(*customerTgId))
		if err != nil {
			return c.Send("Какая-то ошибка...")
		}
		if barber.Missing() {
			log.Println("INFO: Чел как-то нажал кнопку создать смену")
			return c.Send("ты кто?")
		}
		now := time.Now()
		tgerr := c.Send(fmt.Sprintf("Напиши дату смены и со скольких до скольких в формате:\n<b>%s 11:00-19:00</b>", now.Format("02.01.2006")), tele.ModeHTML)
		if tgerr != nil {
			log.Fatal(tgerr)
		}
		err = stateManager.State(*customerTgId).Set(fsm.ShiftEnter)
		if err != nil {
			log.Fatal(err)
		}
		return tgerr
	}
}

func HandleGetShift(store *repository.Store) Handler {
	return func(c tele.Context) error {
		barber, err := store.Barber().GetByTelegramId(uint64(c.Chat().ID))
		if err != nil {
			return c.Send("Какая-то ошибка...")
		}
		if barber.Missing() {
			log.Println("INFO: Чел как-то нажал не на свою смену")
			return c.Send("ты кто?")
		}
		shiftId := c.Callback().Data
		shift, err := store.Shift().Get(shiftId)
		if err != nil {
			return c.Send("Какая-то ошибка, попробуй позже")
		}
		if shift.Missing() {
			return c.Send("Не та смена")
		}
		var (
			needBtn   bool
			btnAction tele.Btn
		)
		switch shift.Status {
		case string(models.Planned):
			btnAction = BarberShiftsInlineKeyboard.Data("Начать смену", "start", shiftId)
			needBtn = true
		case string(models.Work):
			btnAction = BarberShiftsInlineKeyboard.Data("Завершить смену", "finish", shiftId)
			needBtn = true
		default:
			needBtn = false
		}
		BarberShiftsInlineKeyboard.Inline()
		if needBtn == true {
			BarberShiftsInlineKeyboard.Inline(
				BarberShiftsInlineKeyboard.Row(
					BarberShiftsInlineKeyboard.Data(BtnCancelShift.Text, BtnCancelShift.Unique, shiftId), btnAction),
			)
		}
		var txt = fmt.Sprintf("<b>%s</b>\n\nСтатус: %s\n", shift.String(), shift.Status)
		dateSortedVisits := make(utils.TimeSlice, 0, len(shift.Visits))
		sort.Sort(dateSortedVisits)
		for index, visit := range shift.Visits {
			totalPrice, _ := visit.TotalPrice.Get()
			visitTxt := fmt.Sprintf("\n<b>%d. %s - %s</b>\n%s %d ₽\n%s +%s\n",
				index+1, visit.PlannedFrom.Add(barber.TimeOffset()).Format("15:04"), visit.PlannedTo.Add(barber.TimeOffset()).Format("15:04"),
				visit.Service.Title, totalPrice,
				visit.Customer.FullName, visit.Customer.Phone)
			txt += visitTxt
		}
		return c.Send(txt, BarberShiftsInlineKeyboard, tele.ModeHTML)
	}
}

func HandleStartShift(store *repository.Store) Handler {
	return func(c tele.Context) error {
		barber, err := store.Barber().GetByTelegramId(uint64(c.Chat().ID))
		if err != nil {
			return c.Send("Какая-то ошибка...")
		}
		if barber.Missing() {
			log.Println("INFO: Чел как-то нажал не на свою смену")
			return c.Send("ты кто?")
		}
		shiftId := c.Callback().Data
		_, err = store.Shift().UpdateStatus(shiftId, models.Work)
		if err != nil {
			return c.Send("Какая-то ошибка...")
		}
		return HandleGetShift(store)(c)
	}
}

func HandleFinishShift(store *repository.Store) Handler {
	return func(c tele.Context) error {
		barber, err := store.Barber().GetByTelegramId(uint64(c.Chat().ID))
		if err != nil {
			return c.Send("Какая-то ошибка...")
		}
		if barber.Missing() {
			log.Println("INFO: Чел как-то нажал не на свою смену")
			return c.Send("ты кто?")
		}
		shiftId := c.Callback().Data
		_, err = store.Shift().UpdateStatus(shiftId, models.Finished)
		if err != nil {
			return c.Send("Какая-то ошибка...")
		}
		return HandleGetShift(store)(c)
	}
}

func HandleCancelShift(store *repository.Store) Handler {
	return func(c tele.Context) error {
		barber, err := store.Barber().GetByTelegramId(uint64(c.Chat().ID))
		if err != nil {
			return c.Send("Какая-то ошибка...")
		}
		if barber.Missing() {
			log.Println("INFO: Чел как-то нажал не на свою смену")
			return c.Send("ты кто?")
		}
		shiftId := c.Callback().Data
		shift, err := store.Shift().Cancel(shiftId)
		if err != nil {
			return c.Send("Какая-то ошибка...")
		}
		for _, visit := range shift.Visits {
			err := services.NotifyCustomerAboutCancel(c.Bot(), *barber, visit)
			if err != nil {
				c.Send(fmt.Sprintf("Не получилось оповестить %s +%s", visit.Customer.FullName, visit.Customer.Phone))
			}
		}
		c.Send("Смена отменена, все записи тоже, клиентов оповестили (если не сказано иное), все тип-топ")
		return HandleGetShift(store)(c)
	}
}

func HandleMainServices(store *repository.Store) Handler {
	return func(c tele.Context) error {
		barber, err := store.Barber().GetByTelegramId(uint64(c.Chat().ID))
		if err != nil {
			return c.Send("Какая-то ошибка...")
		}
		if barber.Missing() {
			log.Println("INFO: Чел как-то нажал не на свою смену")
			return c.Send("ты кто?")
		}
		barberServices, err := store.Service().GetAll(barber.Id.String())
		if err != nil {
			return c.Send("Какая-то ошибка...")
		}
		buttons := make([]tele.Btn, len(barberServices)+1) // количество смен + кнопка для создания смен + кнопка все смены
		for _, service := range barberServices {
			var btn = BarberShiftsInlineKeyboard.Data(service.String(), "barberToService", service.Id.String())
			buttons = append(buttons, btn)
		}
		buttons = append(buttons, BtnCreateService)
		var rows = BarberShiftsInlineKeyboard.Split(1, buttons)
		BarberShiftsInlineKeyboard.Inline(rows...)
		return c.Send("Твои прайсы", BarberShiftsInlineKeyboard)
	}
}

func HandleStartCreateService(stateManager *fsm.StateManager) Handler {
	return func(c tele.Context) error {
		err := c.Send("Введи название услуги")
		stateManager.State(c.Chat().ID).Set(fsm.ServiceEnterTitle)
		return err
	}
}

func HandleServiceEnterTitle(manager *fsm.StateManager, c tele.Context) error {
	manager.Data(c.Chat().ID).Set("title", c.Text())
	err := c.Send("Введи стоимость услуги в рублях целым числом\nнапример <b>1000</b>", tele.ModeHTML)
	manager.State(c.Chat().ID).Set(fsm.ServiceEnterPrice)
	return err
}

func HandleServiceEnterPrice(manager *fsm.StateManager, c tele.Context) error {
	price, err := strconv.Atoi(c.Text())
	if err != nil {
		return c.Send("Невалидная стоимость. Укажи целое значение, например <b>1000</b>", tele.ModeHTML)
	}
	manager.Data(c.Chat().ID).Set("price", strconv.Itoa(price))
	err = c.Send("Введи продолжительность в минутах одним числом\nнапример <b>60</b>, если услуга длится час\n"+
		"или <b>90</b>, если услуга займет полтора часа", tele.ModeHTML)
	manager.State(c.Chat().ID).Set(fsm.ServiceEnterDuration)
	return err
}

func HandleEndCreateService(store *repository.Store, manager *fsm.StateManager, barber models.Barber, c tele.Context) error {
	duration, err := strconv.Atoi(c.Text())
	if err != nil {
		return c.Send("Невалидная продолжительность. Укажи целое значение,"+
			"\nнапример <b>60</b>, если услуга длится час,"+
			"\nили <b>90</b>, если услуга займет полтора часа", tele.ModeHTML)
	}
	title := manager.Data(c.Chat().ID).Get("title")
	priceStr := manager.Data(c.Chat().ID).Get("price")
	price, _ := strconv.Atoi(priceStr)
	log.Println(priceStr)
	log.Println(price)
	log.Println(time.Minute * time.Duration(duration))

	var service = &models.Service{
		Title:    title,
		Price:    int64(price),
		Duration: edgedb.Duration(time.Minute * time.Duration(duration)),
	}
	service, err = store.Service().Create(barber.Id.String(), service)
	if err != nil {
		return c.Send("Какая-то ошибка...")
	}
	return c.Send(
		fmt.Sprintf("Создана услуга\n\n<b>%s</b>\nЦена: <b>%d ₽</b>\nПродолжительность: <b>%d минут</b>", title, price, duration),
		tele.ModeHTML,
	)
}

func HandleCustomers(store *repository.Store) Handler {
	return func(c tele.Context) error {
		customers, err := store.Customer().GetAll()
		if err != nil {
			return c.Send("Не получилось")
		}
		txt := "👨‍👦 Клиентосы"
		for _, customer := range customers {
			txt += fmt.Sprintf("\n\n%s. +%s", customer.FullName, customer.Phone)
		}
		return c.Send(txt)
	}
}
