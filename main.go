package main

import (
	"benny/src/handlers"
	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v3"
	"log"
	"os"
	"time"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	pref := tele.Settings{
		Token:  os.Getenv("BOT_TOKEN"),
		Poller: &tele.LongPoller{Timeout: 20 * time.Second},
	}
	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
		return
	}

	b.Handle("/start", handlers.HandleStart())

	// barber handlers
	b.Handle(&handlers.BtnShifts, handlers.HandleMainShifts())
	b.Handle(&handlers.BtnPlannedShifts, handlers.HandleMainShifts())
	b.Handle(&handlers.BtnAllShifts, handlers.HandleAllShifts())
	b.Handle(&handlers.BtnCreateShift, handlers.HandleStartCreateShift())
	b.Handle(&handlers.BtnGetShift, handlers.HandleGetShift())
	b.Handle(&handlers.BtnStartShift, handlers.HandleStartShift())
	b.Handle(&handlers.BtnFinishShift, handlers.HandleFinishShift())
	b.Handle(tele.OnText, handlers.HandleText())

	// customer handlers
	b.Handle(tele.OnContact, handlers.HandleReceivePhone())
	b.Handle(&handlers.BtnCreateVisit, handlers.HandleStartCreateVisit())
	b.Handle(&handlers.BtnSelectBarber, handlers.HandleSelectBarber())

	log.Println("INFO: Bot started...")
	b.Start()
}
