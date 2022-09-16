package main

import (
	"benny/src/fsm"
	"benny/src/handlers"
	"benny/src/repository"
	"benny/src/services"
	"context"
	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
	"log"
	"os"
	"time"
)

func setHandlers(b *tele.Bot, store *repository.Store, stateManager *fsm.StateManager) {
	b.Handle("/start", tele.HandlerFunc(handlers.HandleStart(store)))

	// barber handlers
	b.Handle(&handlers.BtnShifts, tele.HandlerFunc(handlers.HandleMainShifts(store)))
	b.Handle(&handlers.BtnPlannedShifts, tele.HandlerFunc(handlers.HandleMainShifts(store)))
	b.Handle(&handlers.BtnAllShifts, tele.HandlerFunc(handlers.HandleAllShifts(store)))
	b.Handle(&handlers.BtnCreateShift, tele.HandlerFunc(handlers.HandleStartCreateShift(store, stateManager)))
	b.Handle(&handlers.BtnGetShift, tele.HandlerFunc(handlers.HandleGetShift(store)))
	b.Handle(&handlers.BtnStartShift, tele.HandlerFunc(handlers.HandleStartShift(store)))
	b.Handle(&handlers.BtnFinishShift, tele.HandlerFunc(handlers.HandleFinishShift(store)))
	b.Handle(&handlers.BtnServices, tele.HandlerFunc(handlers.HandleMainServices(store)))
	b.Handle(&handlers.BtnCreateService, tele.HandlerFunc(handlers.HandleStartCreateService(stateManager)))
	b.Handle(tele.OnText, tele.HandlerFunc(handlers.HandleText(store, stateManager)))

	// customer handlers
	b.Handle(tele.OnContact, tele.HandlerFunc(handlers.HandleReceivePhone(store)))
	b.Handle(&handlers.BtnCreateVisit, tele.HandlerFunc(handlers.HandleStartCreateVisit(store, stateManager)))
	//b.Handle(&handlers.BtnSelectBarber, tele.HandlerFunc(handlers.HandleSelectBarber()))
	b.Handle(&handlers.BtnSelectService, tele.HandlerFunc(handlers.HandleSelectService(store, stateManager)))
	b.Handle(&handlers.BtnSelectShiftToVisit, tele.HandlerFunc(handlers.HandleSelectShift(store, stateManager)))
	b.Handle(&handlers.BtnSelectTimeToVisit, tele.HandlerFunc(handlers.HandleSelectTime(store, stateManager)))
	b.Handle(&handlers.BtnAcceptVisit, tele.HandlerFunc(handlers.HandleAcceptVisit(store, stateManager)))
	b.Handle(&handlers.BtnDeclineVisit, tele.HandlerFunc(handlers.HandleDeclineVisit(stateManager)))
}

func main() {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
	}
	pref := tele.Settings{
		Token:  os.Getenv("BOT_TOKEN"),
		Poller: &tele.LongPoller{Timeout: 20 * time.Second},
	}

	ctx := context.Background()
	store, closer := repository.New(ctx)
	defer closer()

	stateManager, managerCloser := fsm.New(ctx)
	defer managerCloser()

	b, err := tele.NewBot(pref)
	if err != nil {
		log.Fatal(err)
	}
	b.Use(middleware.AutoRespond())
	setHandlers(b, store, stateManager)
	log.Println("INFO: Add task for create shifts...")
	go services.CreateNewBarberShiftOnNextWeek(*b)
	log.Println("INFO: Bot started...")
	b.Start()
}
