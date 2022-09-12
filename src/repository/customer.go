package repository

import (
	"benny/src/models"
	"context"
	"fmt"
	"github.com/edgedb/edgedb-go"
	"log"
)

type CustomerRepository interface {
	Create(customer *models.Customer) *models.Customer
	GetByTelegramId(telegramId int64) (models.Customer, bool)
}

type CustomerRepositoryImpl struct {
	ctx    context.Context
	client *edgedb.Client
}

func (r *CustomerRepositoryImpl) Create(customer *models.Customer) *models.Customer {
	var query = fmt.Sprintf("insert Customer{"+
		"fullName := '%s', "+
		"phone := '%s', "+
		"telegramId := %d"+
		"};", customer.FullName, customer.Phone, customer.TelegramId)
	err := r.client.QuerySingle(r.ctx, query, customer)
	if err != nil {
		log.Fatal(err)
	}
	return customer
}

func (r *CustomerRepositoryImpl) GetByTelegramId(telegramId int64) (models.Customer, bool) {
	var customer models.Customer
	var query = fmt.Sprintf("select Customer{id, fullName, phone, timeZoneOffset} filter .telegramId = %d;", telegramId)
	err := r.client.QuerySingle(r.ctx, query, &customer)
	if err != nil {
		log.Fatal(err)
	}
	return customer, customer.Missing()
}
