package core

import (
	referral_resolvers "github.com/root9464/Go_GamlerDefi/modules/referral"
	test_resolvers "github.com/root9464/Go_GamlerDefi/modules/test"
	gqlgen "github.com/root9464/Go_GamlerDefi/packages/generated/gql_generated"
)

var _ gqlgen.MutationResolver = &Resolver{}
var _ gqlgen.QueryResolver = &Resolver{}

type Resolver struct {
	*test_resolvers.TestResolver
	*referral_resolvers.ReferralResolver
}

type AppMutationResolvers struct{ *Resolver }
type AppQueryResolvers struct{ *Resolver }

func (r *Resolver) Mutation() gqlgen.MutationResolver {
	return &AppMutationResolvers{r}
}

func (r *Resolver) Query() gqlgen.QueryResolver {
	return &AppQueryResolvers{r}
}

func (app *Core) init_gql_resolvers() {
	app.logger.Info("Initializing Graphql resolvers")
	resolvers := &Resolver{
		TestResolver:     &test_resolvers.TestResolver{},
		ReferralResolver: referral_resolvers.NewReferralResolver(app.database, app.logger),
	}
	app.gql_resolvers = resolvers
	app.logger.Success("Graphql resolvers initialized")
}
