package handlers

import tele "gopkg.in/telebot.v3"

var (
	MainBarberKeyboard = &tele.ReplyMarkup{ResizeKeyboard: true}
	BtnShifts          = MainBarberKeyboard.Text("🗓 Мои смены")

	BarberShiftsInlineKeyboard = &tele.ReplyMarkup{}
	BtnAllShifts               = BarberShiftsInlineKeyboard.Data("🗄 Показать все смены", "all")
	BtnPlannedShifts           = BarberShiftsInlineKeyboard.Data("📖 Показать только запланированные", "planned")
	BtnCreateShift             = BarberShiftsInlineKeyboard.Data("📆 Создать смену", "create")
	BtnGetShift                = BarberShiftsInlineKeyboard.Data("Перейти к смене", "barberToShift")
	BtnStartShift              = BarberShiftsInlineKeyboard.Data("✅ Начать смену", "start")
	BtnFinishShift             = BarberShiftsInlineKeyboard.Data("❎ Закончить смену", "finish")
	BtnCancelShift             = BarberShiftsInlineKeyboard.Data("🚫 Отменить смену", "canceled")

	PhoneRequestKeyboard = &tele.ReplyMarkup{ResizeKeyboard: true}
	BtnRequestPhone      = PhoneRequestKeyboard.Contact("☎️ Поделиться цифрами")

	MainCustomerKeyboard = &tele.ReplyMarkup{ResizeKeyboard: true}
	BtnCreateVisit       = MainCustomerKeyboard.Text("Записаться на стригу")

	CustomerShiftsInlineKeyboard = &tele.ReplyMarkup{}
	BtnGetShiftToVisit           = CustomerShiftsInlineKeyboard.Data("Перейти к смене", "customerToShift")
)
