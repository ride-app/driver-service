package service

import (
	"context"
	"errors"
	"strings"

	"github.com/bufbuild/connect-go"
	pb "github.com/ride-app/driver-service/api/gen/ride/driver/v1alpha1"
	"github.com/sirupsen/logrus"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (service *DriverServiceServer) CreateDriver(ctx context.Context,
	req *connect.Request[pb.CreateDriverRequest]) (*connect.Response[pb.CreateDriverResponse], error) {

	if err := req.Msg.Validate(); err != nil {
		logrus.Info("Invalid request")

		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	uid := strings.Split(req.Msg.Driver.Name, "/")[1]

	logrus.Debug("uid: ", uid)
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

	if driver != nil {
		logrus.Info("Driver already exists")

		return connect.NewResponse(
			&pb.CreateDriverResponse{
				Driver: driver,
			},
		), nil
	}

	createTime, err := service.driverRepository.CreateDriver(ctx, req.Msg.Driver)

	if err != nil {
		logrus.WithError(err).Error("Failed to create driver")
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	req.Msg.Driver.CreateTime = timestamppb.New(*createTime)
	req.Msg.Driver.UpdateTime = timestamppb.New(*createTime)

	if err := req.Msg.Validate(); err != nil {
		logrus.WithError(err).Error("Invalid response")
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	logrus.Info("Driver created")
	return connect.NewResponse(&pb.CreateDriverResponse{
		Driver: req.Msg.Driver,
	}), nil
}
