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
	TelegramId     uint64      `edgedb:"telegramId"`
	Visits         []Visit     `edgedb:"visits"`
	TimeZoneOffset int64       `edgedb:"timeZoneOffset"`
}

type Visit struct {
	edgedb.Optional
	BarberShift   BarberShift             `edgedb:"barberShift"`
	Customer      Customer                `edgedb:"customer"`
	Id            edgedb.UUID             `edgedb:"id"`
	CustomerId    edgedb.UUID             `edgedb:"customerId"`
	BarberShiftId edgedb.UUID             `edgedb:"barberShiftId"`
	PlannedFrom   time.Time               `edgedb:"plannedFrom"`
	PlannedTo     time.Time               `edgedb:"plannedTo"`
	ActualFrom    edgedb.OptionalDateTime `edgedb:"actualFrom"`
	ActualTo      edgedb.OptionalDateTime `edgedb:"actualTo"`
	ServiceType   string                  `edgedb:"serviceType"`
	Price         uint64                  `edgedb:"price"`
	DiscountPrice uint64                  `edgedb:"discountPrice"`
	TotalPrice    uint64                  `edgedb:"totalPrice"`
	Status        VisitStatus             `edgedb:"status"`
}

type Barber struct {
	edgedb.Optional
	Id             edgedb.UUID   `edgedb:"id"`
	FullName       string        `edgedb:"fullName"`
	Phone          string        `edgedb:"phone"`
	AvailableTypes []string      `edgedb:"availableTypes"`
	TelegramId     int64         `edgedb:"telegramId"`
	Shifts         []BarberShift `edgedb:"shifts"`
	TimeZoneOffset int64         `edgedb:"timeZoneOffset"`
}

type BarberShift struct {
	edgedb.Optional
	Barber      Barber                  `edgedb:"barber"`
	Id          edgedb.UUID             `edgedb:"id"`
	BarberId    edgedb.UUID             `edgedb:"barberId"`
	Visits      []Visit                 `edgedb:"visits"`
	Status      string                  `edgedb:"status"`
	PlannedFrom time.Time               `edgedb:"plannedFrom"`
	PlannedTo   time.Time               `edgedb:"plannedTo"`
	ActualFrom  edgedb.OptionalDateTime `edgedb:"actualFrom"`
	ActualTo    edgedb.OptionalDateTime `edgedb:"actualTo"`
	Deleted     bool                    `edgedb:"deleted"`
}

func (b BarberShift) String() string {
	return fmt.Sprintf("%s %s до %s",
		b.PlannedFrom.Format("02.01.2006"),
		b.PlannedFrom.Add(time.Hour*time.Duration(b.Barber.TimeZoneOffset)).Format("15:04"),
		b.PlannedTo.Add(time.Hour*time.Duration(b.Barber.TimeZoneOffset)).Format("15:04"))
}
