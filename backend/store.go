package main

import (
	"context"
	"errors"
	"sync"
	"time"
)

var ErrNotFound = errors.New("circle not found")

type Circle struct {
	ID        int64     `json:"id"`
	Name      string    `json:"name"`
	CreatedAt time.Time `json:"created_at"`
}

type Store interface {
	Count(ctx context.Context) (int, error)
	List(ctx context.Context) ([]Circle, error)
	Add(ctx context.Context, name string) (Circle, error)
	Delete(ctx context.Context, id int64) error
}

type MemStore struct {
	mu      sync.Mutex
	circles []Circle
	nextID  int64
}

func NewMemStore() *MemStore {
	return &MemStore{nextID: 1}
}

func (m *MemStore) Count(_ context.Context) (int, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.circles), nil
}

func (m *MemStore) List(_ context.Context) ([]Circle, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]Circle, len(m.circles))
	copy(out, m.circles)
	return out, nil
}

func (m *MemStore) Add(_ context.Context, name string) (Circle, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	c := Circle{ID: m.nextID, Name: name, CreatedAt: time.Now().UTC()}
	m.nextID++
	m.circles = append(m.circles, c)
	return c, nil
}

func (m *MemStore) Delete(_ context.Context, id int64) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	for i, c := range m.circles {
		if c.ID == id {
			m.circles = append(m.circles[:i], m.circles[i+1:]...)
			return nil
		}
	}
	return ErrNotFound
}
