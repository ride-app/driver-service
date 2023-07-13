package service

import (
	"context"
	"errors"
	"strings"

	"github.com/bufbuild/connect-go"
	pb "github.com/ride-app/driver-service/api/gen/ride/driver/v1alpha1"
)

func (service *DriverServiceServer) GetLocation(ctx context.Context,
	req *connect.Request[pb.GetLocationRequest]) (*connect.Response[pb.GetLocationResponse], error) {

	if err := req.Msg.Validate(); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	driverId := strings.Split(req.Msg.Name, "/")[1]

	if driverId != req.Header().Get("Authorization") {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("permission denied"))
	}

	location, err := service.driverRepository.GetLocation(ctx, driverId)

	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if location == nil {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("location not found"))
	}

	return connect.NewResponse(&pb.GetLocationResponse{
		Location: location,
	}), nil
}
