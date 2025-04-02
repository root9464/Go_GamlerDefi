package test_user_resolvers

import (
	"context"
	"fmt"

	"github.com/root9464/Go_GamlerDefi/modules/test/model"
)

var _ TodoQueries = &Resolver{}
var _ TodoMutations = &Resolver{}

type Resolver struct{}

type TodoQueries interface {
	Todos(ctx context.Context) ([]*model.Todo, error)
	Todo(ctx context.Context, id string) (*model.Todo, error)
}

type TodoMutations interface {
	CreateTodo(ctx context.Context, input model.NewTodo) (*model.Todo, error)
}

func (r *Resolver) CreateTodo(ctx context.Context, input model.NewTodo) (*model.Todo, error) {
	panic(fmt.Errorf("not implemented: CreateTodo - createTodo"))
}

func (r *Resolver) Todos(ctx context.Context) ([]*model.Todo, error) {

	return []*model.Todo{
		{
			ID:   "1",
			Text: "Test",
			Done: false,
		},
	}, nil
}

func (r *Resolver) Todo(ctx context.Context, id string) (*model.Todo, error) {
	panic(fmt.Errorf("not implemented: Todo - todo"))
}
