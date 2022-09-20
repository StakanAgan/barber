package repository

import (
	"benny/src/models"
	"context"
	"fmt"
	"github.com/edgedb/edgedb-go"
	"log"
)

type BarberRepository interface {
	Create(barber *models.Barber) (*models.Barber, error)
	GetByTelegramId(telegramId uint64) (*models.Barber, error)
	Get(barberId string) (*models.Barber, error)
	GetFirst() (*models.Barber, error)
	GetAll() ([]models.Barber, error)
}

type BarberRepositoryImpl struct {
	ctx    context.Context
	client edgedb.Client
}

func (r *BarberRepositoryImpl) Create(barber *models.Barber) (*models.Barber, error) {
	var result models.Barber
	var query = fmt.Sprintf("insert Barber {"+
		"fullName := '%s', "+
		"phone := '%s',"+
		"telegramId := %d"+
		"};", barber.FullName, barber.Phone, barber.TelegramId)
	err := r.client.QuerySingle(r.ctx, query, &result)
	if err != nil {
		log.Printf("ERROR: error on create barber: err: %s", err)
	}

	return &result, err
}

func (r *BarberRepositoryImpl) GetByTelegramId(telegramId uint64) (*models.Barber, error) {
	var result models.Barber
	var query = fmt.Sprintf("select Barber {id, fullName, phone, telegramId, timeZoneOffset} filter .telegramId = %d;", telegramId)
	err := r.client.QuerySingle(r.ctx, query, &result)
	if err != nil {
		log.Printf("ERROR: error on get barber by tg id, barberTgId: %d, err: %s", telegramId, err)
	}
	return &result, err
}

func (r *BarberRepositoryImpl) Get(barberId string) (*models.Barber, error) {
	var barber models.Barber
	var query = fmt.Sprintf("select Barber{"+
		"id, fullName, phone, telegramId, timeZoneOffset,"+
		" services: {id, title, price, duration},"+
		" shifts: {id, barber: {timeZoneOffset}, plannedFrom, plannedTo, status}"+
		"} filter .id = <uuid>'%s';", barberId)
	err := r.client.QuerySingle(r.ctx, query, &barber)
	if err != nil {
		log.Printf("ERROR: error on get barber, barberId: %s, err: %s", barberId, err)
	}
	return &barber, err
}

func (r *BarberRepositoryImpl) GetFirst() (*models.Barber, error) {
	var barber models.Barber
	var query = "select Barber{" +
		"id, fullName, phone, telegramId, timeZoneOffset," +
		" services: {id, title, price, duration}," +
		" shifts: {id, barber: {timeZoneOffset}, plannedFrom, plannedTo, status}" +
		"} limit 1;"
	err := r.client.QuerySingle(r.ctx, query, &barber)
	if err != nil {
		log.Printf("ERROR: error on get first barber, err: %s", err)
	}
	return &barber, err
}

func (r *BarberRepositoryImpl) GetAll() ([]models.Barber, error) {
	var barbers []models.Barber
	var query = "select Barber{" +
		"id, fullName, phone, timeZoneOffset" +
		"}" +
		" filter count(.shifts filter .status = ShiftStatus.Planned or .status = ShiftStatus.Work) > 0" +
		" and count(.services filter .deleted = false) > 0;"
	err := r.client.Query(r.ctx, query, &barbers)
	if err != nil {
		log.Printf("ERROR: error on get all barbers, err: %s", err)
	}
	return barbers, err
}
