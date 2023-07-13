package service

import (
	"context"
	"errors"
	"strings"

	"github.com/bufbuild/connect-go"
	pb "github.com/ride-app/driver-service/api/gen/ride/driver/v1alpha1"
)

func (service *DriverServiceServer) GoOffline(ctx context.Context,
	req *connect.Request[pb.GoOfflineRequest]) (*connect.Response[pb.GoOfflineResponse], error) {
	if err := req.Msg.Validate(); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	driverId := strings.Split(req.Msg.Name, "/")[1]

	if driverId != req.Header().Get("Authorization") {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("permission denied"))
	}

	status, err := service.driverRepository.GoOffline(ctx, driverId)

	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&pb.GoOfflineResponse{
		Status: status,
	}), nil
}
