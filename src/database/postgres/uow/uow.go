package uow

import (
	"context"
	"errors"
	"fmt"
	"gorm.io/gorm"
	"sync"
)

var (
	ErrTxAlreadyStarted   = errors.New("transaction already started")
	ErrTxNotStarted       = errors.New("no transaction started")
	ErrRepositoryNotFound = errors.New("repository not registered")
)

type RepositoryFactory func(tx *gorm.DB) (interface{}, error)

type Uow interface {
	RegisterRepository(name string, factory RepositoryFactory)
	GetRepository(ctx context.Context, name string) (interface{}, error)
	Do(ctx context.Context, fn func(ctx context.Context) error) error
	Begin(ctx context.Context) error
	Commit() error
	Rollback() error
}

type uow struct {
	db           *gorm.DB
	tx           *gorm.DB
	repositories map[string]RepositoryFactory
	mu           sync.RWMutex
}

func New(db *gorm.DB) Uow {
	return &uow{
		db:           db,
		repositories: make(map[string]RepositoryFactory),
	}
}

func (u *uow) RegisterRepository(name string, factory RepositoryFactory) {
	u.mu.Lock()
	defer u.mu.Unlock()
	u.repositories[name] = factory
}

func (u *uow) GetRepository(ctx context.Context, name string) (interface{}, error) {
	u.mu.RLock()
	factory, exists := u.repositories[name]
	u.mu.RUnlock()

	if !exists {
		return nil, fmt.Errorf("%w: %s", ErrRepositoryNotFound, name)
	}

	if u.tx == nil {
		return nil, ErrTxNotStarted
	}

	return factory(u.tx)
}

func (u *uow) Begin(ctx context.Context) error {
	if u.tx != nil {
		return ErrTxAlreadyStarted
	}

	u.tx = u.db.Begin()
	if u.tx.Error != nil {
		return fmt.Errorf("failed to begin transaction: %w", u.tx.Error)
	}

	return nil
}

func (u *uow) Commit() error {
	if u.tx == nil {
		return ErrTxNotStarted
	}

	if err := u.tx.Commit().Error; err != nil {
		return fmt.Errorf("failed to commit transaction: %w", err)
	}

	u.tx = nil
	return nil
}

func (u *uow) Rollback() error {
	if u.tx == nil {
		return ErrTxNotStarted
	}

	if err := u.tx.Rollback().Error; err != nil {
		return fmt.Errorf("failed to rollback transaction: %w", err)
	}

	u.tx = nil
	return nil
}

func (u *uow) Do(ctx context.Context, fn func(ctx context.Context) error) error {
	if u.tx != nil {
		return ErrTxAlreadyStarted
	}

	if err := u.Begin(ctx); err != nil {
		return err
	}

	var fnErr error
	defer func() {
		if p := recover(); p != nil {
			_ = u.Rollback()
			panic(p)
		}

		if fnErr != nil {
			if rbErr := u.Rollback(); rbErr != nil {
				fnErr = fmt.Errorf("original error: %w, rollback error: %v", fnErr, rbErr)
			}
		}
	}()

	ctx = context.WithValue(ctx, "uow", u)
	fnErr = fn(ctx)

	if fnErr != nil {
		return fnErr
	}

	if err := u.Commit(); err != nil {
		return fmt.Errorf("commit failed: %w", err)
	}

	return nil
}
