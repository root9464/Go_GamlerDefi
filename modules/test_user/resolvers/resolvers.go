package test_user_resolvers

import (
	"context"
)

var _ TodoQueries = &TestUserResolver{}
var _ TodoMutations = &TestUserResolver{}

type TestUserResolver struct{}

type TodoQueries interface {
	Ping(ctx context.Context) (string, error)
}

type TodoMutations interface{}

func (r *TestUserResolver) Ping(ctx context.Context) (string, error) {
	return "pong", nil
}
