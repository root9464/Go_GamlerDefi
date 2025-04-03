package test_resolvers

import (
	"context"
)

var _ TodoQueries = &TestResolver{}
var _ TodoMutations = &TestResolver{}

type TestResolver struct{}

type TodoQueries interface {
	Ping(ctx context.Context) (string, error)
}

type TodoMutations interface {
	Ping(ctx context.Context) (string, error)
}

func (r *TestResolver) Ping(ctx context.Context) (string, error) {
	return "pong", nil
}
