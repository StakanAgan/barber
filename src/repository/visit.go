package repository

import (
	"benny/src/models"
	"context"
	"fmt"
	"github.com/edgedb/edgedb-go"
	"time"
)

type VisitRepository interface {
	Create(visit *models.Visit) (*models.Visit, error)
	//Get(visitId string) (models.Barber, bool)
	//GetAllByShiftId(shiftId string) ([]models.Visit, bool)
	//GetAllByCustomerId(customerId string) ([]models.Visit, bool)
}

type VisitRepositoryImpl struct {
	ctx    context.Context
	client edgedb.Client
}

func (r *VisitRepositoryImpl) Create(visit *models.Visit) (*models.Visit, error) {
	var query = fmt.Sprintf("with"+
		" shiftId := <uuid>'%s',"+
		" customerId := <uuid>'%s',"+
		" serviceId := <uuid>'%s'"+
		"insert Visit{plannedFrom := <datetime>'%s', plannedTo := <datetime>'%s', "+
		"barberShift := (select BarberShift filter .id = shiftId),"+
		"service := (select Service filter .id = serviceId),"+
		"customer := (select Customer filter .id = customerId),"+
		"price := %d, status := VisitStatus.%s"+
		"};", visit.BarberShift.Id, visit.Customer.Id, visit.Service.Id,
		visit.PlannedFrom.Format(time.RFC3339), visit.PlannedTo.Format(time.RFC3339),
		visit.Price, visit.Status,
	)
	err := r.client.QuerySingle(r.ctx, query, visit)

	return visit, err
}
