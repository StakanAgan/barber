package repository

import (
	"benny/src/models"
	"context"
	"fmt"
	"github.com/edgedb/edgedb-go"
	"log"
)

type ServiceRepository interface {
	Create(barberId string, service *models.Service) (*models.Service, error)
	GetAll(barberId string) ([]models.Service, error)
	Get(serviceId string) (models.Service, error)
	//Update(barberId edgedb.UUID, service *models.Service) *models.Service
	//Delete(barberId edgedb.UUID, serviceId edgedb.UUID) error
}

type ServiceRepositoryImpl struct {
	ctx    context.Context
	client edgedb.Client
}

func (r *ServiceRepositoryImpl) Create(barberId string, service *models.Service) (*models.Service, error) {
	var query = fmt.Sprintf(
		"with barberId := <uuid>'%s'"+
			" insert Service{"+
			"barber := (select Barber filter .id = barberId),"+
			" title := '%s',"+
			" price := %d,"+
			" duration := <duration>'%d minutes'"+
			"}",
		barberId, service.Title, service.Price, service.Duration/60_000_000_000,
	)
	err := r.client.QuerySingle(r.ctx, query, service)
	if err != nil {
		log.Printf("ERROR: On create service, barberId: %s, service: %s", barberId, service.String())
	}
	return service, err
}

func (r *ServiceRepositoryImpl) GetAll(barberId string) ([]models.Service, error) {
	var query = fmt.Sprintf("select Service{id, title, price, duration} filter .barber.id = <uuid>'%s'", barberId)
	var services []models.Service
	err := r.client.Query(r.ctx, query, &services)
	if err != nil {
		log.Printf("ERROR: On get all barber services, barberId: %s", barberId)
	}
	return services, err
}

func (r *ServiceRepositoryImpl) Get(serviceId string) (models.Service, error) {
	var query = fmt.Sprintf("select Service{id, title, price, duration} filter .id = <uuid>'%s'", serviceId)
	var service models.Service
	err := r.client.QuerySingle(r.ctx, query, &service)
	if err != nil {
		log.Printf("ERROR: On get service, serviceid: %s", serviceId)
	}
	return service, err
}
