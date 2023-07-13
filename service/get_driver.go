package service

import (
	"context"
	"errors"
	"strings"

	"github.com/bufbuild/connect-go"
	pb "github.com/ride-app/driver-service/api/gen/ride/driver/v1alpha1"
)

func (service *DriverServiceServer) GetDriver(ctx context.Context,
	req *connect.Request[pb.GetDriverRequest]) (*connect.Response[pb.GetDriverResponse], error) {

	if err := req.Msg.Validate(); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	if req.Msg.Name == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("name cannot be empty"))
	}

	driverId := strings.Split(req.Msg.Name, "/")[1]

	if driverId != req.Header().Get("Authorization") {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("permission denied"))
	}

	driver, err := service.driverRepository.GetDriver(ctx, driverId)

	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if driver == nil {
		return nil, connect.NewError(connect.CodeNotFound, errors.New("driver not found"))
	}

	return connect.NewResponse(&pb.GetDriverResponse{
		Driver: driver,
	}), nil
}
