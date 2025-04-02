package graph

import (
	graphql_out "github.com/root9464/Go_GamlerDefi/graph"
	test_resolvers "github.com/root9464/Go_GamlerDefi/modules/test/resolvers"
	"go.mongodb.org/mongo-driver/mongo"
)

type Resolver struct {
	*test_resolvers.Resolver
}

type MutationResolver struct {
	*Resolver
}

type QueryResolver struct {
	*Resolver
}

func (r *Resolver) Mutation() graphql_out.MutationResolver {
	return &MutationResolver{r}
}

func (r *Resolver) Query() graphql_out.QueryResolver {
	return &QueryResolver{r}
}

func New(mdb *mongo.Client) *Resolver {
	r := &Resolver{}

	return r
}
