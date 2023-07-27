package service

import (
	"context"
	"errors"
	"strings"

	"github.com/bufbuild/connect-go"
	pb "github.com/ride-app/driver-service/api/gen/ride/driver/v1alpha1"
	"github.com/sirupsen/logrus"
)

func (service *DriverServiceServer) GetDriver(ctx context.Context,
	req *connect.Request[pb.GetDriverRequest]) (*connect.Response[pb.GetDriverResponse], error) {

	if err := req.Msg.Validate(); err != nil {
		logrus.Info("Invalid request")
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	if req.Msg.Name == "" {
		logrus.Info("Name is empty")
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("name cannot be empty"))
	}

	uid := strings.Split(req.Msg.Name, "/")[1]

	logrus.Info("uid: ", uid)
	logrus.Debug("Request header uid: ", req.Header().Get("uid"))

	if uid != req.Header().Get("uid") {
		logrus.Info("Permission denied")
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("permission denied"))
	}

	driver, err := service.driverRepository.GetDriver(ctx, uid)

	if err != nil {
		logrus.WithError(err).Error("Failed to get driver")
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if driver == nil {
		logrus.Info("Driver not found")
		return nil, connect.NewError(connect.CodeNotFound, errors.New("driver not found"))
	}

	res := &pb.GetDriverResponse{
		Driver: driver,
	}

	if err := res.Validate(); err != nil {
		logrus.WithError(err).Error("Invalid response")
		return nil, err
	}

	logrus.Info("Driver found")
	return connect.NewResponse(res), nil
}
