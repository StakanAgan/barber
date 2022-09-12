package main

import (
	"benny/src/handlers"
	"benny/src/repository"
	"context"
	"fmt"
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

	ctx := context.Background()
	client, closer := repository.NewDBClient(ctx)
	var result string
	err = client.QuerySingle(ctx, "SELECT 'EdgeDB connected...'", &result)
	if err != nil {
		log.Fatal(fmt.Sprintf("Can't connect to DB, err: %s", err))
	}

	log.Println(result)
	closer()

	b.Handle("/start", handlers.HandleStart())

	// barber handlers
	b.Handle(&handlers.BtnShifts, handlers.HandleMainShifts())
	b.Handle(&handlers.BtnPlannedShifts, handlers.HandleMainShifts())
	b.Handle(&handlers.BtnAllShifts, handlers.HandleAllShifts())
	b.Handle(&handlers.BtnCreateShift, handlers.HandleStartCreateShift())
	b.Handle(&handlers.BtnGetShift, handlers.HandleGetShift())
	b.Handle(&handlers.BtnStartShift, handlers.HandleStartShift())
	b.Handle(&handlers.BtnFinishShift, handlers.HandleFinishShift())
	b.Handle(&handlers.BtnServices, handlers.HandleMainServices())
	b.Handle(&handlers.BtnCreateService, handlers.HandleStartCreateService())
	b.Handle(tele.OnText, handlers.HandleText())

	// customer handlers
	b.Handle(tele.OnContact, handlers.HandleReceivePhone())
	b.Handle(&handlers.BtnCreateVisit, handlers.HandleStartCreateVisit())
	b.Handle(&handlers.BtnSelectBarber, handlers.HandleSelectBarber())
	b.Handle(&handlers.BtnSelectService, handlers.HandleSelectService())
	b.Handle(&handlers.BtnSelectShiftToVisit, handlers.HandleSelectShift())
	b.Handle(&handlers.BtnSelectTimeToVisit, handlers.HandleSelectTime())

	log.Println("INFO: Bot started...")
	b.Start()
}
