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

type StateManager struct {
	ctx          context.Context
	client       *redis.Client
	stateManager *UserStateManagerImpl
	dataManager  *UserDataManagerImpl
}

var config = src.NewRedisConfig()

func New(ctx context.Context) (*StateManager, func() error) {
	client := redis.NewClient(config)
	return &StateManager{ctx: ctx, client: client}, client.Close
}

func (m *StateManager) State(telegramId int64) *UserStateManagerImpl {
	if m.stateManager == nil {
		m.stateManager = &UserStateManagerImpl{ctx: m.ctx, client: m.client, key: fmt.Sprintf(stateKey, telegramId)}
	}

	return m.stateManager
}

func (m *StateManager) Data(telegramId int64) *UserDataManagerImpl {
	if m.dataManager == nil {
		m.dataManager = &UserDataManagerImpl{ctx: m.ctx, client: m.client, key: fmt.Sprintf(dataKey, telegramId)}
	}

	return m.dataManager
}

type UserStateManager interface {
	Get() State
	Set(state State) error
	Reset()
}

type UserStateManagerImpl struct {
	ctx    context.Context
	client *redis.Client
	key    string
}

type UserDataManager interface {
	Set(key string, value string) error
	Get(key string) (value string)
	Reset()
}

type UserDataManagerImpl struct {
	ctx    context.Context
	client *redis.Client
	key    string
}

func (s *UserStateManagerImpl) Get() State {
	state, err := s.client.Get(s.ctx, s.key).Result()
	if err != nil {
		if err != redis.Nil {
			log.Printf("ERROR: error on Get state, key: %s, err: %s", s.key, err)
		}
	}
	return State(state)
}

func (s *UserStateManagerImpl) Set(state State) error {
	return s.client.Set(s.ctx, s.key, state.String(), 0).Err()
}

func (s *UserStateManagerImpl) Reset() {
	s.client.Del(s.ctx, s.key)
}

func (d *UserDataManagerImpl) Set(key string, value string) error {
	userDataKey := fmt.Sprintf("%s_%s", d.key, key)
	return d.client.Set(d.ctx, userDataKey, value, 0).Err()
}

func (d *UserDataManagerImpl) Get(key string) (value string) {
	userDataKey := fmt.Sprintf("%s_%s", d.key, key)
	value, err := d.client.Get(d.ctx, userDataKey).Result()
	if err != nil {
		if err != redis.Nil {
			log.Printf("ERROR: error on Get data, key: %s, err: %s", d.key, err)
		}
	}
	return value
}

func (d *UserDataManagerImpl) Reset() {
	keys, err := d.client.Keys(d.ctx, d.key).Result()
	if err != nil {
		if err != redis.Nil {
			log.Printf("ERROR: error on Reset data, key: %s, err: %s", d.key, err)
		}
	}
	d.client.Del(d.ctx, keys...)
}
