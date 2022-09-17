package handlers

import tele "gopkg.in/telebot.v3"

var (
	BtnShifts          = MainBarberKeyboard.Text("🗓 Смены")
	BtnServices        = MainBarberKeyboard.Text("🧾 Цены")
	MainBarberKeyboard = &tele.ReplyMarkup{ResizeKeyboard: true}

	BarberShiftsInlineKeyboard = &tele.ReplyMarkup{}
	BtnAllShifts               = BarberShiftsInlineKeyboard.Data("🗄 Показать все смены", "all")
	BtnPlannedShifts           = BarberShiftsInlineKeyboard.Data("📖 Показать только актуальные", "planned")
	BtnCreateShift             = BarberShiftsInlineKeyboard.Data("📆 Создать смену", "create")
	BtnGetShift                = BarberShiftsInlineKeyboard.Data("Перейти к смене", "barberToShift")
	BtnStartShift              = BarberShiftsInlineKeyboard.Data("✅ Начать смену", "start")
	BtnFinishShift             = BarberShiftsInlineKeyboard.Data("❎ Закончить смену", "finish")
	BtnCancelShift             = BarberShiftsInlineKeyboard.Data("🚫 Отменить смену", "canceled")

	BarberServicesInlineKeyboard = &tele.ReplyMarkup{}
	BtnGetService                = BarberServicesInlineKeyboard.Data("Перейти к сервису", "barberToService")
	BtnCreateService             = BarberServicesInlineKeyboard.Data("🖋 Добавить услугу", "createService")

	PhoneRequestKeyboard = &tele.ReplyMarkup{ResizeKeyboard: true, RemoveKeyboard: true, OneTimeKeyboard: true}
	BtnRequestPhone      = PhoneRequestKeyboard.Contact("☎️ Поделиться цифрами")

	MainCustomerKeyboard = &tele.ReplyMarkup{ResizeKeyboard: true}
	BtnCreateVisit       = MainCustomerKeyboard.Data("Записаться на стригу", "createVisit")

	CustomerShiftsInlineKeyboard = &tele.ReplyMarkup{}
	BtnSelectBarber              = CustomerShiftsInlineKeyboard.Data("Выбрать барбера", "customerToBarber")
	BtnSelectService             = CustomerShiftsInlineKeyboard.Data("Выбрать услугу", "customerToService")
	BtnSelectShiftToVisit        = CustomerShiftsInlineKeyboard.Data("Перейти к смене", "customerToShift")
	BtnSelectTimeToVisit         = CustomerShiftsInlineKeyboard.Data("Выбрать время", "customerToTime")
	BtnDeclineVisit              = CustomerShiftsInlineKeyboard.Data("❌", "customerDeclineVisit")
	BtnAcceptVisit               = CustomerShiftsInlineKeyboard.Data("✅", "customerAcceptVisit")
)
