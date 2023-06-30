package service

import (
	"context"
	"errors"

	"github.com/bufbuild/connect-go"
	pb "github.com/ride-app/driver-service/api/gen/ride/driver/v1alpha1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (service *DriverServiceServer) UpdateDriver(ctx context.Context,
	req *connect.Request[pb.UpdateDriverRequest]) (*connect.Response[pb.UpdateDriverResponse], error) {

	if err := req.Msg.Validate(); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	if req.Msg.Driver.Name == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("name cannot be empty"))
	}

	updateTime, err := service.driverRepository.UpdateDriver(ctx, req.Msg.Driver)

	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	req.Msg.Driver.UpdateTime = timestamppb.New(*updateTime)

	if err := req.Msg.Validate(); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&pb.UpdateDriverResponse{
		Driver: req.Msg.Driver,
	}), nil
}
