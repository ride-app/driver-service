package service

import (
	"context"
	"errors"
	"strings"

	"github.com/bufbuild/connect-go"
	pb "github.com/ride-app/driver-service/api/gen/ride/driver/v1alpha1"
	"github.com/sirupsen/logrus"
)

func (service *DriverServiceServer) UpdateVehicle(ctx context.Context,
	req *connect.Request[pb.UpdateVehicleRequest]) (*connect.Response[pb.UpdateVehicleResponse], error) {
	if err := req.Msg.Validate(); err != nil {
		logrus.Info("Invalid request")
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	uid := strings.Split(req.Msg.Vehicle.Name, "/")[1]

	logrus.Info("uid: ", uid)
	logrus.Debug("Request header uid: ", req.Header().Get("uid"))

	if uid != req.Header().Get("uid") {
		logrus.Info("Permission denied")
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("permission denied"))
	}

	_, err := service.vehicleRepository.UpdateVehicle(ctx, req.Msg.Vehicle)

	if err != nil {
		logrus.WithError(err).Error("Failed to update vehicle")
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := &pb.UpdateVehicleResponse{
		Vehicle: req.Msg.Vehicle,
	}

	if err := res.Validate(); err != nil {
		logrus.WithError(err).Error("Invalid response")
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	logrus.Info("Vehicle updated")
	return connect.NewResponse(res), nil
}
