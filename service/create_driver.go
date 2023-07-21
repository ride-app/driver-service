package service

import (
	"context"
	"errors"
	"strings"

	"github.com/bufbuild/connect-go"
	pb "github.com/ride-app/driver-service/api/gen/ride/driver/v1alpha1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (service *DriverServiceServer) CreateDriver(ctx context.Context,
	req *connect.Request[pb.CreateDriverRequest]) (*connect.Response[pb.CreateDriverResponse], error) {

	if err := req.Msg.Validate(); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	driverId := strings.Split(req.Msg.Driver.Name, "/")[1]

	if driverId != req.Header().Get("Authorization") {
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("permission denied"))
	}

	driver, err := service.driverRepository.GetDriver(ctx, strings.Split(req.Msg.Driver.Name, "/")[1])

	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if driver != nil {
		return connect.NewResponse(
			&pb.CreateDriverResponse{
				Driver: driver,
			},
		), nil
	}

	createTime, err := service.driverRepository.CreateDriver(ctx, req.Msg.Driver)

	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	req.Msg.Driver.CreateTime = timestamppb.New(*createTime)
	req.Msg.Driver.UpdateTime = timestamppb.New(*createTime)

	if err := req.Msg.Validate(); err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&pb.CreateDriverResponse{
		Driver: req.Msg.Driver,
	}), nil
}
