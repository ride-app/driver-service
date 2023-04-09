package service

import (
	"context"
	"strings"

	"github.com/bufbuild/connect-go"
	pb "github.com/ride-app/driver-service/api/gen/ride/driver/v1alpha1"
)

func (service *DriverServiceServer) GetVehicle(ctx context.Context,
	req *connect.Request[pb.GetVehicleRequest]) (*connect.Response[pb.GetVehicleResponse], error) {

	if err := req.Msg.Validate(); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	driverId := strings.Split(req.Msg.Name, "/")[1]

	vehicle, err := service.vehicleRepository.GetVehicle(ctx, driverId)

	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&pb.GetVehicleResponse{
		Vehicle: vehicle,
	}), nil
}
