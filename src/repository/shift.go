package repository

import (
	models2 "benny/src/models"
	"context"
	"errors"
	"fmt"
	"github.com/edgedb/edgedb-go"
	"log"
	"time"
)

type ShiftRepository interface {
	Create(barberId edgedb.UUID, shift *models2.BarberShift) (*models2.BarberShift, error)
	GetAll(barberId edgedb.UUID) ([]models2.BarberShift, bool)
	GetActual(barberId edgedb.UUID) ([]models2.BarberShift, bool)
	Get(shiftId edgedb.UUID) (models2.BarberShift, bool)
	Delete(shiftId edgedb.UUID) bool
	UpdateStatus(shiftId edgedb.UUID, status models2.ShiftStatus) (models2.BarberShift, bool)
}

type ShiftRepositoryImpl struct {
	ctx    context.Context
	client *edgedb.Client
}

func (r *ShiftRepositoryImpl) Create(barberId edgedb.UUID, shift *models2.BarberShift) (*models2.BarberShift, error) {
	ctx := context.Background()
	client, closer := NewDBClient(ctx)
	defer closer()

	plannedFromStr := shift.PlannedFrom.Format(time.RFC3339)
	plannedToStr := shift.PlannedTo.Format(time.RFC3339)
	var isShiftCrossing bool

	var validateQuery = fmt.Sprintf("select exists(select BarberShift"+
		" filter .barber.id = <uuid>'%s'"+
		" and (.plannedFrom >= <datetime>'%s' and .plannedFrom <= <datetime>'%s')"+
		" or (.plannedTo >= <datetime>'%s' and .plannedTo <= <datetime>'%s'))",
		barberId, plannedFromStr, plannedToStr, plannedFromStr, plannedToStr)
	err := client.QuerySingle(ctx, validateQuery, &isShiftCrossing)
	if isShiftCrossing == true {
		return shift, errors.New("shift will crossing with another shifts")
	}
	var query = fmt.Sprintf("with barberId := <uuid>'%s' insert BarberShift {"+
		"barber := (select Barber filter .id = barberId), "+
		"status := ShiftStatus.Planned,"+
		"plannedFrom := <datetime>'%s',"+
		"plannedTo := <datetime>'%s',"+
		"};", barberId, plannedFromStr, plannedToStr)
	err = client.QuerySingle(ctx, query, shift)
	if err != nil {
		log.Fatal(err)
	}

	return shift, nil
}

func (r *ShiftRepositoryImpl) GetAll(barberId edgedb.UUID) ([]models2.BarberShift, bool) {
	ctx := context.Background()
	client, closer := NewDBClient(ctx)
	defer closer()
	var shifts []models2.BarberShift
	var query = fmt.Sprintf("select BarberShift{id, barber: {fullName, timeZoneOffset}, plannedFrom, plannedTo} filter .barber.id = <uuid>'%s';", barberId)
	err := client.Query(ctx, query, &shifts)
	if err != nil {
		log.Fatal(err)
	}

	return shifts, len(shifts) == 0
}

func (r *ShiftRepositoryImpl) GetActual(barberId edgedb.UUID) ([]models2.BarberShift, bool) {
	ctx := context.Background()
	client, closer := NewDBClient(ctx)
	defer closer()
	var shifts []models2.BarberShift
	var query = fmt.Sprintf("select BarberShift "+
		"{id, barber: {fullName, timeZoneOffset}, plannedFrom, plannedTo}"+
		" filter .barber.id = <uuid>'%s'"+
		" and .status = ShiftStatus.%s or .status = ShiftStatus.%s;", barberId, models2.Planned, models2.Work)
	err := client.Query(ctx, query, &shifts)
	if err != nil {
		log.Fatal(err)
	}

	return shifts, len(shifts) == 0
}

func (r *ShiftRepositoryImpl) Get(shiftId edgedb.UUID) (models2.BarberShift, bool) {
	ctx := context.Background()
	client, closer := NewDBClient(ctx)
	defer closer()
	var shift models2.BarberShift
	var query = fmt.Sprintf("select BarberShift{"+
		"id, barber: {fullName, timeZoneOffset}, status, plannedFrom, plannedTo, actualFrom, actualTo, visits: {"+
		"customer: {fullName, phone}}}"+
		" filter .id = <uuid>'%s';", shiftId)
	err := client.QuerySingle(ctx, query, &shift)
	if err != nil {
		log.Fatal(err)
	}
	return shift, shift.Missing()

}

func (r *ShiftRepositoryImpl) Delete(shiftId edgedb.UUID) bool {
	ctx := context.Background()
	client, closer := NewDBClient(ctx)
	defer closer()
	var shift models2.BarberShift

	var query = fmt.Sprintf("update BarberShift filter .id=<uuid>'%s' set {deleted := true};", shiftId)
	err := client.QuerySingle(ctx, query, &shift)
	if err != nil {
		log.Fatal(err)
	}
	return shift.Missing()
}

func (r *ShiftRepositoryImpl) UpdateStatus(shiftId edgedb.UUID, status models2.ShiftStatus) (models2.BarberShift, bool) {
	ctx := context.Background()
	client, closer := NewDBClient(ctx)
	defer closer()

	var query = fmt.Sprintf("update BarberShift filter .id=<uuid>'%s' set {status := ShiftStatus.%s}", shiftId, status)
	var shift models2.BarberShift
	err := client.QuerySingle(ctx, query, &shift)
	if err != nil {
		log.Fatal(err)
	}
	return shift, shift.Missing()
}
