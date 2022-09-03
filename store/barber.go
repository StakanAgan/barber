package store

import (
	"benny/models"
	"context"
	"fmt"
	"github.com/edgedb/edgedb-go"
	"log"
)

type BarberRepository interface {
	Create(barber *models.Barber) (*models.Barber, bool)
	GetByTelegramId(telegramId uint64) (*models.Barber, bool)
	GetAll() ([]models.Barber, bool)
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
		"availableTypes := [ServiceType.Hair, ServiceType.Beard, ServiceType.HairBeard],"+
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
	var query = fmt.Sprintf("select Barber {id, fullName, phone, availableTypes, telegramId, timeZoneOffset} filter .telegramId = %d;", telegramId)
	err := r.client.QuerySingle(r.ctx, query, &result)
	if err != nil {
		log.Fatal(err)
	}
	return &result, result.Missing()
}

func (r *BarberRepositoryImpl) GetAll() ([]models.Barber, bool) {
	var barbers []models.Barber
	var query = "select Barber{id, fullName, phone, availableTypes, timeZoneOffset} filter len(.shifts) > 0;"
	err := r.client.Query(r.ctx, query, barbers)
	if err != nil {
		log.Fatal(err)
	}
	return barbers, len(barbers) == 0
}
