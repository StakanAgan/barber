package store

import (
	"context"
	"github.com/edgedb/edgedb-go"
	"log"
)

type Store struct {
	ctx                context.Context
	client             *edgedb.Client
	customerRepository *CustomerRepositoryImpl
	barberRepository   *BarberRepositoryImpl
	shiftRepository    *ShiftRepositoryImpl
}

func NewDBClient(ctx context.Context) (*edgedb.Client, func()) {
	opts := edgedb.Options{}
	client, err := edgedb.CreateClient(ctx, opts)
	if err != nil {
		log.Fatal(err)
	}
	return client, func() {
		client.Close()
	}
}

func New(ctx context.Context) (*Store, func()) {
	client, closer := NewDBClient(ctx)
	store := &Store{
		ctx:    ctx,
		client: client,
	}
	return store, func() {
		store.barberRepository = nil
		store.customerRepository = nil
		store.shiftRepository = nil
		closer()
	}
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
