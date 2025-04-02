package graph

import (
	test_resolvers "github.com/root9464/Go_GamlerDefi/modules/test/resolvers"
	"go.mongodb.org/mongo-driver/mongo"
)

type Resolver struct {
	*test_resolvers.Resolver
}

type AppMutationResolvers struct {
	*Resolver
}

type AppQueryResolvers struct {
	*Resolver
}

func (r *Resolver) Mutation() MutationResolver {
	return &AppMutationResolvers{r}
}

func (r *Resolver) Query() QueryResolver {
	return &AppQueryResolvers{r}
}

func New(mdb *mongo.Client) *Resolver {
	r := &Resolver{}

	return r
}
