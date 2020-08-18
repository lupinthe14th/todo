package todo

import (
	"context"
	"sync"
	"time"

	"github.com/google/uuid"
)

type DB interface {
	Put(ctx context.Context, todo *TODO) error
	GetAll(ctx context.Context) ([]*TODO, error)
}

type MemoryDB struct {
	sync.RWMutex
	m map[string]*TODO
}

var _ DB = (*MemoryDB)(nil)

func NewMemoryDB() *MemoryDB {
	return &MemoryDB{m: map[string]*TODO{}}
}

func (db *MemoryDB) Put(ctx context.Context, todo *TODO) error {
	if todo.ID == "" {
		id, err := uuid.NewRandom()
		if err != nil {
			return err
		}
		todo.ID = id.String()
	}
	if todo.CreatedAt.IsZero() {
		todo.CreatedAt = time.Now()
	}

	db.Lock()
	db.m[todo.ID] = todo
	db.Unlock()

	return nil
}
func (db *MemoryDB) GetAll(ctx context.Context) ([]*TODO, error) {
	db.RLock()
	defer db.RUnlock()

	todos := make([]*TODO, 0, len(db.m))

	for _, todo := range db.m {
		todos = append(todos, todo)
	}
	return todos, nil
}
