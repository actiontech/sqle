package resolver

// This file will not be regenerated automatically.
//
// It serves as dependency injection for your app, add any dependencies you require here.

type Resolver struct{}

type MutationResolverImpl struct {
	*mutationResolver
}

type QueryResolverImpl struct {
	*queryResolver
}
