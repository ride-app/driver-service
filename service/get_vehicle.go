package service

import (
	"context"
	"errors"
	"strings"

	"github.com/bufbuild/connect-go"
	pb "github.com/ride-app/driver-service/api/gen/ride/driver/v1alpha1"
	"github.com/sirupsen/logrus"
)

func (service *DriverServiceServer) GetVehicle(ctx context.Context,
	req *connect.Request[pb.GetVehicleRequest]) (*connect.Response[pb.GetVehicleResponse], error) {
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

	vehicle, err := service.vehicleRepository.GetVehicle(ctx, uid)

	if err != nil {
		logrus.WithError(err).Error("Failed to get vehicle")
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if vehicle == nil {
		logrus.Info("Vehicle not found")
		return nil, connect.NewError(connect.CodeNotFound, errors.New("vehicle not found"))
	}

	res := &pb.GetVehicleResponse{
		Vehicle: vehicle,
	}

	if err := res.Validate(); err != nil {
		logrus.WithError(err).Error("Invalid response")
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	logrus.Info("Vehicle found")
	return connect.NewResponse(res), nil
}
