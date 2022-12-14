package models

import (
	"fmt"
	"github.com/edgedb/edgedb-go"
	"time"
)

type Customer struct {
	edgedb.Optional
	Id             edgedb.UUID `edgedb:"id"`
	FullName       string      `edgedb:"fullName"`
	Phone          string      `edgedb:"phone"`
	TelegramId     int64       `edgedb:"telegramId"`
	Visits         []Visit     `edgedb:"visits"`
	TimeZoneOffset int64       `edgedb:"timeZoneOffset"`
}

type Visit struct {
	edgedb.Optional
	BarberShift   BarberShift             `edgedb:"barberShift"`
	Customer      Customer                `edgedb:"customer"`
	Service       Service                 `edgedb:"service"`
	Id            edgedb.UUID             `edgedb:"id"`
	PlannedFrom   time.Time               `edgedb:"plannedFrom"`
	PlannedTo     time.Time               `edgedb:"plannedTo"`
	ActualFrom    edgedb.OptionalDateTime `edgedb:"actualFrom"`
	ActualTo      edgedb.OptionalDateTime `edgedb:"actualTo"`
	Price         int64                   `edgedb:"price"`
	DiscountPrice int64                   `edgedb:"discountPrice"`
	TotalPrice    edgedb.OptionalInt64    `edgedb:"totalPrice"`
	Status        VisitStatus             `edgedb:"status"`
}

type Barber struct {
	edgedb.Optional
	Id             edgedb.UUID   `edgedb:"id"`
	FullName       string        `edgedb:"fullName"`
	Phone          string        `edgedb:"phone"`
	Services       []Service     `edgedb:"services"`
	TelegramId     int64         `edgedb:"telegramId"`
	Shifts         []BarberShift `edgedb:"shifts"`
	TimeZoneOffset int64         `edgedb:"timeZoneOffset"`
}

type BarberShift struct {
	edgedb.Optional
	Barber      Barber                  `edgedb:"barber"`
	Id          edgedb.UUID             `edgedb:"id"`
	Visits      []Visit                 `edgedb:"visits"`
	Status      string                  `edgedb:"status"`
	PlannedFrom time.Time               `edgedb:"plannedFrom"`
	PlannedTo   time.Time               `edgedb:"plannedTo"`
	ActualFrom  edgedb.OptionalDateTime `edgedb:"actualFrom"`
	ActualTo    edgedb.OptionalDateTime `edgedb:"actualTo"`
	Deleted     bool                    `edgedb:"deleted"`
}

type Service struct {
	edgedb.Optional
	Barber   Barber          `edgedb:"barber"`
	Id       edgedb.UUID     `edgedb:"id"`
	Title    string          `edgedb:"title"`
	Price    int64           `edgedb:"price"`
	Duration edgedb.Duration `edgedb:"duration"`
}

func (b Barber) TimeOffset() time.Duration {
	return time.Hour * time.Duration(b.TimeZoneOffset)
}

func (b BarberShift) String() string {
	return fmt.Sprintf("%s %s ???? %s",
		b.PlannedFrom.Add(b.Barber.TimeOffset()).Format("02.01.2006"),
		b.PlannedFrom.Add(b.Barber.TimeOffset()).Format("15:04"),
		b.PlannedTo.Add(b.Barber.TimeOffset()).Format("15:04"))
}

func (s Service) String() string {
	return fmt.Sprintf("%s - %d ?????????? - %d ???", s.Title, s.Duration/60_000_000, s.Price)
}
