package service

import (
	"context"
	"errors"
	"strings"

	"github.com/bufbuild/connect-go"
	pb "github.com/ride-app/driver-service/api/gen/ride/driver/v1alpha1"
)

func (service *DriverServiceServer) GetStatus(ctx context.Context,
	req *connect.Request[pb.GetStatusRequest]) (*connect.Response[pb.GetStatusResponse], error) {

	if err := req.Msg.Validate(); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	driverId := strings.Split(req.Msg.Name, "/")[1]

	if driverId != req.Header().Get("Authorization") {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("permission denied"))
	}

	status, err := service.driverRepository.GetStatus(ctx, driverId)

	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if status == nil {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("status not found"))
	}

	return connect.NewResponse(&pb.GetStatusResponse{
		Status: status,
	}), nil
}
