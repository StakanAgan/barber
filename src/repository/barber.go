package repository

import (
	"benny/src/models"
	"context"
	"fmt"
	"github.com/edgedb/edgedb-go"
	"log"
)

type BarberRepository interface {
	Create(barber *models.Barber) (*models.Barber, bool)
	GetByTelegramId(telegramId uint64) (*models.Barber, bool)
	GetWithShifts(barberId edgedb.UUID) (*models.Barber, bool)
	GetAllWithShifts() ([]models.Barber, bool)
}

type BarberRepositoryImpl struct {
	ctx    context.Context
	client *edgedb.Client
}

func (r *BarberRepositoryImpl) Create(barber *models.Barber) (*models.Barber, bool) {
	var result models.Barber
	var query = fmt.Sprintf("insert Barber {"+
		"fullName := '%s', "+
		"phone := '%s',"+
		"telegramId := %d"+
		"};", barber.FullName, barber.Phone, barber.TelegramId)
	err := r.client.QuerySingle(r.ctx, query, &result)
	if err != nil {
		log.Fatal(err)
	}

	return &result, result.Missing()
}

func (r *BarberRepositoryImpl) GetByTelegramId(telegramId uint64) (*models.Barber, bool) {
	var result models.Barber
	var query = fmt.Sprintf("select Barber {id, fullName, phone, telegramId, timeZoneOffset} filter .telegramId = %d;", telegramId)
	err := r.client.QuerySingle(r.ctx, query, &result)
	if err != nil {
		log.Fatal(err)
	}
	return &result, result.Missing()
}

func (r *BarberRepositoryImpl) GetWithShifts(barberId edgedb.UUID) (*models.Barber, bool) {
	var barber models.Barber
	var query = fmt.Sprintf("select Barber{id, fullName, phone, telegramId, timeZoneOffset, shifts: {id, barber: {timeZoneOffset}, plannedFrom, plannedTo, status} filter .status in {ShiftStatus.Planned, ShiftStatus.Work}} filter .id = <uuid>'%s';", barberId)
	err := r.client.QuerySingle(r.ctx, query, &barber)
	if err != nil {
		log.Fatal(err)
	}
	return &barber, barber.Missing()
}

func (r *BarberRepositoryImpl) GetAllWithShifts() ([]models.Barber, bool) {
	var barbers []models.Barber
	var query = "select Barber{id, fullName, phone, timeZoneOffset} filter count(.shifts filter .status = ShiftStatus.Planned or .status = ShiftStatus.Work) > 0;"
	err := r.client.Query(r.ctx, query, &barbers)
	if err != nil {
		log.Fatal(err)
	}
	return barbers, len(barbers) == 0
}