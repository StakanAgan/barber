package repository

import (
	"benny/src/models"
	"context"
	"fmt"
	"github.com/edgedb/edgedb-go"
	"log"
)

type CustomerRepository interface {
	Create(customer *models.Customer) (*models.Customer, error)
	GetByTelegramId(telegramId int64) (models.Customer, error)
}

type CustomerRepositoryImpl struct {
	ctx    context.Context
	client edgedb.Client
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
	if err != nil {
		log.Printf("ERROR: error on get customer by tg id, tgId: %d, err: %s", telegramId, err)
	}
	return customer, err
}
