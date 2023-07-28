package service

import (
	"context"
	"errors"
	"strings"

	"github.com/bufbuild/connect-go"
	pb "github.com/ride-app/driver-service/api/gen/ride/driver/v1alpha1"
	"github.com/sirupsen/logrus"
)

func (service *DriverServiceServer) GetLocation(ctx context.Context,
	req *connect.Request[pb.GetLocationRequest]) (*connect.Response[pb.GetLocationResponse], error) {

	if err := req.Msg.Validate(); err != nil {
		logrus.Info("Invalid request")
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	uid := strings.Split(req.Msg.Name, "/")[1]

	logrus.Debug("uid: ", uid)
	logrus.Debug("Request header uid: ", req.Header().Get("uid"))

	if uid != req.Header().Get("uid") {
		logrus.Info("Permission denied")
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("permission denied"))
	}

	location, err := service.driverRepository.GetLocation(ctx, uid)

	if err != nil {
		logrus.WithError(err).Error("Failed to get location")
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if location == nil {
		logrus.Info("Location not found")
		return nil, connect.NewError(connect.CodeNotFound, errors.New("location not found"))
	}

	res := &pb.GetLocationResponse{
		Location: location,
	}

	if err := res.Validate(); err != nil {
		logrus.WithError(err).Error("Invalid response")
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	logrus.Info("Location found")
	return connect.NewResponse(res), nil
}
