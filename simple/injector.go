//go:build wireinject
// +build wireinject

package simple

func InitializedSimpleService() *SimpleService {
	wire.Build(NewSimpleRepository, NewSimpleService)
	return nil
}
