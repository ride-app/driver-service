package service

import (
	"context"
	"errors"
	"strings"

	"github.com/bufbuild/connect-go"
	pb "github.com/ride-app/driver-service/api/gen/ride/driver/v1alpha1"
)

func (service *DriverServiceServer) CreateDriver(ctx context.Context,
	req *connect.Request[pb.CreateDriverRequest]) (*connect.Response[pb.CreateDriverResponse], error) {

	if err := req.Msg.Validate(); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	if req.Msg.Driver.Name == "" {
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("invalid driver name"))
	}

	driver, err := service.driverRepository.GetDriver(ctx, strings.Split(req.Msg.Driver.Name, "/")[1])

	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if driver != nil {
		return nil, connect.NewError(connect.CodeAlreadyExists, err)
	}

	_, err = service.driverRepository.CreateDriver(ctx, req.Msg.Driver)

	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	return connect.NewResponse(&pb.CreateDriverResponse{
		Driver: req.Msg.Driver,
	}), nil
}
