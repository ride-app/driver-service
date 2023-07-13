package service

import (
	"context"
	"errors"
	"strings"

	"github.com/bufbuild/connect-go"
	pb "github.com/ride-app/driver-service/api/gen/ride/driver/v1alpha1"
)

func (service *DriverServiceServer) UpdateLocation(ctx context.Context,
	req *connect.Request[pb.UpdateLocationRequest]) (*connect.Response[pb.UpdateLocationResponse], error) {

	if err := req.Msg.Validate(); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	driverId := strings.Split(req.Msg.Parent, "/")[1]

	if driverId != req.Header().Get("Authorization") {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("permission denied"))
	}

	status, err := service.driverRepository.GetStatus(ctx, driverId)

	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if !status.Online {
		return nil, connect.NewError(connect.CodeFailedPrecondition, err)
	}

	_, err = service.driverRepository.UpdateLocation(ctx, driverId, req.Msg.Location)

	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&pb.UpdateLocationResponse{}), nil
}
