package service

import (
	pb "github.com/ride-app/entity-service/api/gen/ride/entity/v1alpha1"

	er "github.com/ride-app/entity-service/repositories/entity"
)

type EntityServiceServer struct {
	pb.UnimplementedEntityServiceServer

	entityRepository er.EntityRepository
}

func New(
	entityRepository er.EntityRepository,
) *EntityServiceServer {
	return &EntityServiceServer{
		entityRepository: entityRepository,
	}
}
