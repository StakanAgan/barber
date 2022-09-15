package main

import (
	"benny/src/handlers"
	"benny/src/repository"
	"context"
	"fmt"
	"github.com/joho/godotenv"
	tele "gopkg.in/telebot.v3"
	"gopkg.in/telebot.v3/middleware"
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

	ctx := context.Background()
	client, closer := repository.NewDBClient(ctx)
	var result string
	err = client.QuerySingle(ctx, "SELECT 'EdgeDB connected...'", &result)
	if err != nil {
		log.Fatal(fmt.Sprintf("Can't connect to DB, err: %s", err))
	}

	log.Println(result)
	closer()
	b.Use(middleware.AutoRespond())

	b.Handle("/start", tele.HandlerFunc(handlers.HandleStart()))

	// barber handlers
	b.Handle(&handlers.BtnShifts, tele.HandlerFunc(handlers.HandleMainShifts()))
	b.Handle(&handlers.BtnPlannedShifts, tele.HandlerFunc(handlers.HandleMainShifts()))
	b.Handle(&handlers.BtnAllShifts, tele.HandlerFunc(handlers.HandleAllShifts()))
	b.Handle(&handlers.BtnCreateShift, tele.HandlerFunc(handlers.HandleStartCreateShift()))
	b.Handle(&handlers.BtnGetShift, tele.HandlerFunc(handlers.HandleGetShift()))
	b.Handle(&handlers.BtnStartShift, tele.HandlerFunc(handlers.HandleStartShift()))
	b.Handle(&handlers.BtnFinishShift, tele.HandlerFunc(handlers.HandleFinishShift()))
	b.Handle(&handlers.BtnServices, tele.HandlerFunc(handlers.HandleMainServices()))
	b.Handle(&handlers.BtnCreateService, tele.HandlerFunc(handlers.HandleStartCreateService()))
	b.Handle(tele.OnText, tele.HandlerFunc(handlers.HandleText()))

	// customer handlers
	b.Handle(tele.OnContact, tele.HandlerFunc(handlers.HandleReceivePhone()))
	b.Handle(&handlers.BtnCreateVisit, tele.HandlerFunc(handlers.HandleStartCreateVisit()))
	//b.Handle(&handlers.BtnSelectBarber, tele.HandlerFunc(handlers.HandleSelectBarber()))
	b.Handle(&handlers.BtnSelectService, tele.HandlerFunc(handlers.HandleSelectService()))
	b.Handle(&handlers.BtnSelectShiftToVisit, tele.HandlerFunc(handlers.HandleSelectShift()))
	b.Handle(&handlers.BtnSelectTimeToVisit, tele.HandlerFunc(handlers.HandleSelectTime()))
	b.Handle(&handlers.BtnAcceptVisit, tele.HandlerFunc(handlers.HandleAcceptVisit()))
	b.Handle(&handlers.BtnDeclineVisit, tele.HandlerFunc(handlers.HandleDeclineVisit()))

	log.Println("INFO: Bot started...")
	b.Start()
}
