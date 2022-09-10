package fsm

import (
	"benny/src"
	"context"
	"fmt"
	"github.com/go-redis/redis/v9"
	"log"
)

var (
	stateKey = "%d_state"
	dataKey  = "%d_data"
)

type Manager struct {
	ctx          context.Context
	client       *redis.Client
	telegramId   int64
	stateManager *StateManagerImpl
	dataManager  *DataManagerImpl
}

var config = src.NewRedisConfig()

func New(ctx context.Context, telegramId int64) (*Manager, func() error) {
	client := redis.NewClient(config)
	return &Manager{ctx: ctx, client: client, telegramId: telegramId}, client.Close
}

func (m *Manager) State() *StateManagerImpl {
	if m.stateManager == nil {
		m.stateManager = &StateManagerImpl{ctx: m.ctx, client: m.client, key: fmt.Sprintf(stateKey, m.telegramId)}
	}

	return m.stateManager
}

func (m *Manager) Data() *DataManagerImpl {
	if m.dataManager == nil {
		m.dataManager = &DataManagerImpl{ctx: m.ctx, client: m.client, key: fmt.Sprintf(dataKey, m.telegramId)}
	}

	return m.dataManager
}

type StateManager interface {
	Get() State
	Set(state State) error
	Reset()
}

type StateManagerImpl struct {
	ctx    context.Context
	client *redis.Client
	key    string
}

type DataManager interface {
	Set(key string, value string) error
	Get(key string) (value string)
	Reset()
}

type DataManagerImpl struct {
	ctx    context.Context
	client *redis.Client
	key    string
}

func (s *StateManagerImpl) Get() State {
	state, err := s.client.Get(s.ctx, s.key).Result()
	if err != nil {
		if err != redis.Nil {
			log.Fatal(err)
		}
	}
	return State(state)
}

func (s *StateManagerImpl) Set(state State) error {
	return s.client.Set(s.ctx, s.key, state.String(), 0).Err()
}

func (s *StateManagerImpl) Reset() {
	s.client.Del(s.ctx, s.key)
}

func (d *DataManagerImpl) Set(key string, value string) error {
	userDataKey := fmt.Sprintf("%s_%s", d.key, key)
	return d.client.Set(d.ctx, userDataKey, value, 0).Err()
}

func (d *DataManagerImpl) Get(key string) (value string) {
	userDataKey := fmt.Sprintf("%s_%s", d.key, key)
	value, err := d.client.Get(d.ctx, userDataKey).Result()
	if err != nil {
		if err != redis.Nil {
			log.Fatal(err)
		}
	}
	return value
}

func (d *DataManagerImpl) Reset() {
	keys, err := d.client.Keys(d.ctx, d.key).Result()
	if err != nil {
		if err != redis.Nil {
			log.Fatal(err)
		}
	}
	d.client.Del(d.ctx, keys...)
}
