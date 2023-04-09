package service

import (
	"context"

	"github.com/bufbuild/connect-go"
	pb "github.com/ride-app/driver-service/api/gen/ride/driver/v1alpha1"
)

func (service *DriverServiceServer) UpdateDriver(ctx context.Context,
	req *connect.Request[pb.UpdateDriverRequest]) (*connect.Response[pb.UpdateDriverResponse], error) {

	if err := req.Msg.Validate(); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	_, err := service.driverRepository.UpdateDriver(ctx, req.Msg.Driver)

	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&pb.UpdateDriverResponse{
		Driver: req.Msg.Driver,
	}), nil
}
