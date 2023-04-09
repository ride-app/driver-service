package service

import (
	"context"

	"github.com/bufbuild/connect-go"
	pb "github.com/ride-app/driver-service/api/gen/ride/driver/v1alpha1"
)

func (service *DriverServiceServer) UpdateVehicle(ctx context.Context,
	req *connect.Request[pb.UpdateVehicleRequest]) (*connect.Response[pb.UpdateVehicleResponse], error) {

	if err := req.Msg.Validate(); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	_, err := service.vehicleRepository.UpdateVehicle(ctx, req.Msg.Vehicle)

	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&pb.UpdateVehicleResponse{
		Vehicle: req.Msg.Vehicle,
	}), nil
}
