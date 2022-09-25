package repository

import (
	"benny/src/models"
	"context"
	"errors"
	"fmt"
	"github.com/edgedb/edgedb-go"
	"log"
)

type CustomerRepository interface {
	Create(customer *models.Customer) (*models.Customer, error)
	GetByTelegramId(telegramId int64) (models.Customer, error)
	GetAll() ([]models.Customer, error)
}

type CustomerRepositoryImpl struct {
	ctx    context.Context
	client *edgedb.Client
}

func (r *CustomerRepositoryImpl) Create(customer *models.Customer) (*models.Customer, error) {
	var query = fmt.Sprintf("insert Customer{"+
		"fullName := '%s', "+
		"phone := '%s', "+
		"telegramId := %d"+
		"};", customer.FullName, customer.Phone, customer.TelegramId)
	err := r.client.QuerySingle(r.ctx, query, customer)
	if err != nil {
		log.Printf("ERROR: error on create customer, err: %s", err)
	}
	return customer, err
}

func (r *CustomerRepositoryImpl) GetByTelegramId(telegramId int64) (models.Customer, error) {
	var customer models.Customer
	var query = fmt.Sprintf("select Customer{id, fullName, phone, timeZoneOffset} filter .telegramId = %d;", telegramId)
	err := r.client.QuerySingle(r.ctx, query, &customer)
	var edbErr edgedb.Error
	if errors.As(err, &edbErr) && edbErr.Category(edgedb.NoDataError) {
		log.Printf("ERROR: on get customer by tg Id, err: %s", edbErr)
	}
	if err != nil {
		log.Printf("ERROR: error on get customer by tg id, tgId: %d, err: %s", telegramId, err)
	}
	return customer, err
}

func (r *CustomerRepositoryImpl) GetAll() ([]models.Customer, error) {
	var customers []models.Customer
	var query = "select Customer{id, fullName, phone}"
	err := r.client.Query(r.ctx, query, &customers)
	if err != nil {
		log.Printf("ERROR: error on get all customers, err: %s", err)
	}
	return customers, err
}
