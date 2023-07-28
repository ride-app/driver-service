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

func (service *DriverServiceServer) UpdateDriver(ctx context.Context,
	req *connect.Request[pb.UpdateDriverRequest]) (*connect.Response[pb.UpdateDriverResponse], error) {
	if err := req.Msg.Validate(); err != nil {
		logrus.Info("Invalid request")
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	if req.Msg.Driver.Name == "" {
		logrus.Info("Driver name is empty")
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("name cannot be empty"))
	}

	uid := strings.Split(req.Msg.Driver.Name, "/")[1]

	logrus.Debug("uid: ", uid)
	logrus.Debug("Request header uid: ", req.Header().Get("uid"))

	if uid != req.Header().Get("uid") {
		logrus.Info("Permission denied")
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("permission denied"))
	}

	updateTime, err := service.driverRepository.UpdateDriver(ctx, req.Msg.Driver)

	if err != nil {
		logrus.WithError(err).Error("Failed to update driver")
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	req.Msg.Driver.UpdateTime = timestamppb.New(*updateTime)

	res := &pb.UpdateDriverResponse{
		Driver: req.Msg.Driver,
	}

	if err := res.Validate(); err != nil {
		logrus.WithError(err).Error("Invalid response")
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	logrus.Info("Driver updated")
	return connect.NewResponse(res), nil
}
