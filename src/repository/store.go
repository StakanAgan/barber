package repository

import (
	"benny/src"
	"context"
	"fmt"
	"github.com/edgedb/edgedb-go"
	"log"
	"os"
	"time"
)

type Store struct {
	ctx                context.Context
	client             *edgedb.Client
	customerRepository *CustomerRepositoryImpl
	barberRepository   *BarberRepositoryImpl
	shiftRepository    *ShiftRepositoryImpl
	serviceRepository  *ServiceRepositoryImpl
	visitRepository    *VisitRepositoryImpl
}

var config = src.NewDBConfig()

func NewDBClient(ctx context.Context) (*edgedb.Client, func()) {
	opts := edgedb.Options{}
	if os.Getenv("ENV") != "local" {
		opts = edgedb.Options{
			Database:           config.DBName,
			Host:               config.Host,
			User:               config.User,
			Password:           edgedb.NewOptionalStr(config.Password),
			WaitUntilAvailable: 3 * time.Second,
			ConnectTimeout:     5 * time.Second,
			Port:               config.Port,
			TLSOptions: edgedb.TLSOptions{
				SecurityMode: edgedb.TLSModeInsecure,
			},
			Concurrency: 4,
		}
	}

	client, err := edgedb.CreateClient(ctx, opts)
	if err != nil {
		log.Fatal(err)
	}
	DBHealthCheck(client, ctx)
	go func() {
		for range time.Tick(time.Minute * 5) {
			DBHealthCheck(client, ctx)
		}
	}()

	return client, func() {
		err := client.Close()
		if err != nil {
			log.Printf("ERROR: while close DB, err: %s", err)
		}
	}
}

func DBHealthCheck(client *edgedb.Client, ctx context.Context) {
	var result string
	err := client.QuerySingle(ctx, "SELECT 'EdgeDB connected...'", &result)
	if err != nil {
		log.Fatal(fmt.Sprintf("Can't connect to DB, err: %s", err))
	}
	log.Printf("INFO: %s", result)
}

func New(ctx context.Context) (*Store, func()) {
	client, closer := NewDBClient(ctx)
	store := &Store{
		ctx:    ctx,
		client: client,
	}
	return store, closer
}

func (s *Store) Barber() BarberRepository {
	if s.barberRepository == nil {
		s.barberRepository = &BarberRepositoryImpl{ctx: s.ctx, client: s.client}
	}

	return s.barberRepository
}

func (s *Store) Shift() ShiftRepository {
	if s.shiftRepository == nil {
		s.shiftRepository = &ShiftRepositoryImpl{ctx: s.ctx, client: s.client}
	}

	return s.shiftRepository
}

func (s *Store) Customer() CustomerRepository {
	if s.customerRepository == nil {
		s.customerRepository = &CustomerRepositoryImpl{ctx: s.ctx, client: s.client}
	}

	return s.customerRepository
}

func (s *Store) Service() ServiceRepository {
	if s.serviceRepository == nil {
		s.serviceRepository = &ServiceRepositoryImpl{ctx: s.ctx, client: s.client}
	}

	return s.serviceRepository
}

func (s *Store) Visit() VisitRepository {
	if s.visitRepository == nil {
		s.visitRepository = &VisitRepositoryImpl{ctx: s.ctx, client: s.client}
	}

	return s.visitRepository
}
