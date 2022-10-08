//go:generate go run github.com/golang/mock/mockgen -destination ../../mocks/$GOFILE -package mocks . EntityRepository

package entityrepository

import (
	"context"
	"time"

	pb "github.com/ride-app/entity-service/api/gen/ride/entity/v1alpha1"
)

type EntityRepository interface {
	CreateEntity(ctx context.Context, entity *pb.Entity) (createTime *time.Time, err error)

	GetEntity(ctx context.Context, id string) (*pb.Entity, error)

	GetEntitys(ctx context.Context, parentId string) ([]*pb.Entity, error)

	UpdateEntity(ctx context.Context, entity *pb.Entity) (createTime *time.Time, err error)

	DeleteEntity(ctx context.Context, id string) (createTime *time.Time, err error)
}

type SomeImpl struct {
}

func NewSomeEntityRepository() (*SomeImpl, error) {
	return &SomeImpl{}, nil
}

func (r *SomeImpl) CreateEntity(ctx context.Context, entity *pb.Entity) (createTime *time.Time, err error) {
	return nil, nil
}

func (r *SomeImpl) GetEntity(ctx context.Context, id string) (*pb.Entity, error) {
	return nil, nil
}

func (r *SomeImpl) GetEntitys(ctx context.Context, parentId string) ([]*pb.Entity, error) {
	return nil, nil
}

func (r *SomeImpl) UpdateEntity(ctx context.Context, entity *pb.Entity) (createTime *time.Time, err error) {
	return nil, nil
}

func (r *SomeImpl) DeleteEntity(ctx context.Context, id string) (createTime *time.Time, err error) {
	return nil, nil
}
