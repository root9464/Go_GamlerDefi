package test_resolvers

import (
	"context"
	"fmt"

	"github.com/root9464/Go_GamlerDefi/modules/test/model"
)

// Интерфейсы для реализации GraphQL запросов и мутаций
type TodoQueries interface {
	Todos(ctx context.Context) ([]*model.Todo, error)
}

type TodoMutations interface {
	CreateTodo(ctx context.Context, input model.NewTodo) (*model.Todo, error)
}

type Resolver struct{}

func (r *Resolver) CreateTodo(ctx context.Context, input model.NewTodo) (*model.Todo, error) {
	panic(fmt.Errorf("not implemented: CreateTodo - createTodo"))
}

func (r *Resolver) Todos(ctx context.Context) ([]*model.Todo, error) {
	panic(fmt.Errorf("not implemented: Todos - todos"))
}
