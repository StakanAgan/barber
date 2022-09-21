package handlers

import (
	"benny/src/repository"
	"benny/src/utils"
	"fmt"
	tele "gopkg.in/telebot.v3"
	"log"
	"sort"
)

type Handler func(c tele.Context) error

func HandleStart(store *repository.Store) Handler {
	return func(c tele.Context) error {
		log.Printf("INFO: somebody (%d) press /start", c.Chat().ID)
		barber, err := store.Barber().GetByTelegramId(uint64(c.Chat().ID))
		if err != nil {
			return c.Send("Какая-то ошибка...")
		}
		log.Println("INFO: Point 1")
		if barber.Missing() {
			log.Printf("INFO: User %d try to Start bot", uint64(c.Chat().ID))
			customer, err := store.Customer().GetByTelegramId(c.Chat().ID)
			log.Println("INFO: Point 7")
			if err != nil {
				return c.Send("Какая-то ошибка...")
			}
			log.Println("INFO: Point 8")
			if customer.Missing() {
				PhoneRequestKeyboard.Reply(PhoneRequestKeyboard.Row(BtnRequestPhone))
				return c.Send("Заделись цифрами, чтобы записаться на стригу. Просто нажми на <b>☎️ Поделиться цифрами</b> внизу 👇🏼", PhoneRequestKeyboard, tele.ModeHTML)
			}
			MainCustomerKeyboard.Inline(MainCustomerKeyboard.Row(BtnCreateVisit))
			return c.Send(fmt.Sprintf("Велком, %s", customer.FullName), MainCustomerKeyboard)
		}
		log.Println("INFO: Point 2")
		nextShift, err := store.Shift().GetNext(barber.Id.String())
		log.Println("INFO: Point 3")
		if err != nil {
			return c.Send("Ошибочка вышла")
		}
		log.Println("INFO: Point 4")
		txt := fmt.Sprintf("Салют, %s", barber.FullName)
		log.Println("INFO: Point 5")
		if !nextShift.Missing() {
			log.Println("INFO: Point 6")
			txt += fmt.Sprintf("\nСледующая смена <b>%s</b>\nс <b>%s до %s</b>\n\n",
				nextShift.PlannedFrom.Add(barber.TimeOffset()).Format("02.01.2006"),
				nextShift.PlannedFrom.Add(barber.TimeOffset()).Format("15:04"),
				nextShift.PlannedTo.Add(barber.TimeOffset()).Format("15:04"),
			)
			if len(nextShift.Visits) == 0 {
				goto send
			}

			dateSortedVisits := make(utils.TimeSlice, 0, len(nextShift.Visits))
			sort.Sort(dateSortedVisits)
			for index, visit := range nextShift.Visits {
				totalPrice, _ := visit.TotalPrice.Get()
				visitTxt := fmt.Sprintf("\n<b>%d. %s - %s</b>\n%s %d ₽\n%s +%s\n",
					index+1, visit.PlannedFrom.Add(barber.TimeOffset()).Format("15:04"), visit.PlannedTo.Add(barber.TimeOffset()).Format("15:04"),
					visit.Service.Title, totalPrice,
					visit.Customer.FullName, visit.Customer.Phone)
				txt += visitTxt
			}
		}
	send:
		MainBarberKeyboard.Reply(MainBarberKeyboard.Row(BtnShifts), MainBarberKeyboard.Row(BtnServices))
		return c.Send(txt, MainBarberKeyboard, tele.ModeHTML)
	}
}
