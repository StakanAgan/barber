package fsm

import (
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
}

func New(ctx context.Context, telegramId int64) (*Manager, func() error) {
	client := redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "",
		DB:       0,
	})
	return &Manager{ctx: ctx, client: client, telegramId: telegramId}, client.Close
}

func (m *Manager) State() *StateManagerImpl {
	if m.stateManager == nil {
		m.stateManager = &StateManagerImpl{ctx: m.ctx, client: m.client, key: fmt.Sprintf(stateKey, m.telegramId)}
	}

	return m.stateManager
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
	return s.client.Set(s.ctx, s.key, state, 0).Err()
}

func (s *StateManagerImpl) Reset() {
	s.client.Del(s.ctx, s.key)
}
