package service

import (
	"context"
	"errors"
	"strings"

	"github.com/bufbuild/connect-go"
	pb "github.com/ride-app/driver-service/api/gen/ride/driver/v1alpha1"
)

func (service *DriverServiceServer) UpdateVehicle(ctx context.Context,
	req *connect.Request[pb.UpdateVehicleRequest]) (*connect.Response[pb.UpdateVehicleResponse], error) {

	if err := req.Msg.Validate(); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	driverId := strings.Split(req.Msg.Vehicle.Name, "/")[1]

	if driverId != req.Header().Get("Authorization") {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("permission denied"))
	}

	_, err := service.vehicleRepository.UpdateVehicle(ctx, req.Msg.Vehicle)

	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&pb.UpdateVehicleResponse{
		Vehicle: req.Msg.Vehicle,
	}), nil
}
