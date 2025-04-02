package test_resolvers

import (
	"context"
	"fmt"

	"github.com/root9464/Go_GamlerDefi/modules/test/model"
)

var _ TodoQueries = &TestResolver{}
var _ TodoMutations = &TestResolver{}

type TestResolver struct{}

type TodoQueries interface {
	Todos(ctx context.Context) ([]*model.Todo, error)
	Todo(ctx context.Context, id string) (*model.Todo, error)
}

type TodoMutations interface {
	CreateTodo(ctx context.Context, input model.NewTodo) (*model.Todo, error)
}

func (r *TestResolver) CreateTodo(ctx context.Context, input model.NewTodo) (*model.Todo, error) {
	panic(fmt.Errorf("not implemented: CreateTodo - createTodo"))
}

func (r *TestResolver) Todos(ctx context.Context) ([]*model.Todo, error) {

	return []*model.Todo{
		{
			ID:   "1",
			Text: "Test",
			Done: false,
		},
	}, nil
}

func (r *TestResolver) Todo(ctx context.Context, id string) (*model.Todo, error) {
	panic(fmt.Errorf("not implemented: Todo - todo"))
}
