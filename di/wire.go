//go:build wireinject

package di

import (
	"github.com/google/wire"
	entityrepository "github.com/ride-app/entity-service/repositories/entity"
	"github.com/ride-app/entity-service/service"
)

func InitializeService() (*service.EntityServiceServer, error) {
	panic(
		wire.Build(
			entityrepository.NewSomeEntityRepository,
			wire.Bind(
				new(entityrepository.EntityRepository),
				new(*entityrepository.SomeImpl),
			),
			service.New,
		),
	)
}
