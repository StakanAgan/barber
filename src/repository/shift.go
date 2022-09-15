package repository

import (
	"benny/src/models"
	"context"
	"errors"
	"fmt"
	"github.com/edgedb/edgedb-go"
	"log"
	"time"
)

type ShiftRepository interface {
	Create(barberId string, shift *models.BarberShift) (*models.BarberShift, error)
	GetAll(barberId string) ([]models.BarberShift, bool)
	GetActual(barberId string) ([]models.BarberShift, bool)
	Get(shiftId string) (models.BarberShift, bool)
	Delete(shiftId string) bool
	UpdateStatus(shiftId string, status models.ShiftStatus) (models.BarberShift, bool)
}

type ShiftRepositoryImpl struct {
	ctx    context.Context
	client *edgedb.Client
}

func (r *ShiftRepositoryImpl) Create(barberId string, shift *models.BarberShift) (*models.BarberShift, error) {
	plannedFromStr := shift.PlannedFrom.Format(time.RFC3339)
	plannedToStr := shift.PlannedTo.Format(time.RFC3339)
	var isShiftCrossing bool

	var validateQuery = fmt.Sprintf("select exists(select BarberShift"+
		" filter .barber.id = <uuid>'%s'"+
		" and (.plannedFrom >= <datetime>'%s' and .plannedFrom <= <datetime>'%s')"+
		" or (.plannedTo >= <datetime>'%s' and .plannedTo <= <datetime>'%s'))",
		barberId, plannedFromStr, plannedToStr, plannedFromStr, plannedToStr)
	err := r.client.QuerySingle(r.ctx, validateQuery, &isShiftCrossing)
	if isShiftCrossing == true {
		return shift, errors.New("shift will crossing with another shifts")
	}
	var query = fmt.Sprintf("with barberId := <uuid>'%s' insert BarberShift {"+
		"barber := (select Barber filter .id = barberId), "+
		"status := ShiftStatus.Planned,"+
		"plannedFrom := <datetime>'%s',"+
		"plannedTo := <datetime>'%s',"+
		"};", barberId, plannedFromStr, plannedToStr)
	err = r.client.QuerySingle(r.ctx, query, shift)
	if err != nil {
		log.Fatal(err)
	}

	return shift, nil
}

func (r *ShiftRepositoryImpl) GetAll(barberId string) ([]models.BarberShift, bool) {
	var shifts []models.BarberShift
	var query = fmt.Sprintf("select BarberShift{id, barber: {fullName, timeZoneOffset}, plannedFrom, plannedTo} filter .barber.id = <uuid>'%s';", barberId)
	err := r.client.Query(r.ctx, query, &shifts)
	if err != nil {
		log.Fatal(err)
	}

	return shifts, len(shifts) == 0
}

func (r *ShiftRepositoryImpl) GetActual(barberId string) ([]models.BarberShift, bool) {
	var shifts []models.BarberShift
	var query = fmt.Sprintf("select BarberShift "+
		"{id, barber: {fullName, timeZoneOffset}, plannedFrom, plannedTo}"+
		" filter .barber.id = <uuid>'%s' and .plannedTo > datetime_current()"+
		" and .status = ShiftStatus.%s or .status = ShiftStatus.%s;", barberId, models.Planned, models.Work)
	err := r.client.Query(r.ctx, query, &shifts)
	if err != nil {
		log.Fatal(err)
	}

	return shifts, len(shifts) == 0
}

func (r *ShiftRepositoryImpl) Get(shiftId string) (models.BarberShift, bool) {
	var shift models.BarberShift
	var query = fmt.Sprintf("select BarberShift{"+
		"id, barber: {fullName, timeZoneOffset},"+
		" status, plannedFrom, plannedTo, actualFrom, actualTo, visits: {"+
		"customer: {fullName, phone}, plannedFrom, plannedTo, totalPrice, service: {title}"+
		"}"+
		"} filter .id = <uuid>'%s';", shiftId)
	err := r.client.QuerySingle(r.ctx, query, &shift)
	if err != nil {
		log.Fatal(err)
	}
	return shift, shift.Missing()

}

func (r *ShiftRepositoryImpl) Delete(shiftId string) bool {
	var shift models.BarberShift

	var query = fmt.Sprintf("update BarberShift filter .id=<uuid>'%s' set {deleted := true};", shiftId)
	err := r.client.QuerySingle(r.ctx, query, &shift)
	if err != nil {
		log.Fatal(err)
	}
	return shift.Missing()
}

func (r *ShiftRepositoryImpl) UpdateStatus(shiftId string, status models.ShiftStatus) (models.BarberShift, bool) {
	var query = fmt.Sprintf("update BarberShift filter .id=<uuid>'%s' set {status := ShiftStatus.%s}", shiftId, status)
	var shift models.BarberShift
	err := r.client.QuerySingle(r.ctx, query, &shift)
	if err != nil {
		log.Fatal(err)
	}
	return shift, shift.Missing()
}
