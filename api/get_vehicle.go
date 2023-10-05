package api

import (
	"context"
	"errors"
	"strings"

	"connectrpc.com/connect"
	pb "github.com/ride-app/driver-service/proto/ride/driver/v1alpha1"
)

func (service *DriverServiceServer) GetVehicle(ctx context.Context,
	req *connect.Request[pb.GetVehicleRequest]) (*connect.Response[pb.GetVehicleResponse], error) {
	log := service.logger.WithFields(map[string]string{
		"method": "GetVehicle",
	})

	if err := req.Msg.Validate(); err != nil {
		log.WithError(err).Info("Invalid request")
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	uid := strings.Split(req.Msg.Name, "/")[1]

	log.Debug("uid: ", uid)
	log.Debug("Request header uid: ", req.Header().Get("uid"))

	if uid != req.Header().Get("uid") {
		log.Info("Permission denied")
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("permission denied"))
	}

	vehicle, err := service.vehicleRepository.GetVehicle(ctx, log, uid)

	if err != nil {
		log.WithError(err).Error("Failed to get vehicle")
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if vehicle == nil {
		log.Info("Vehicle not found")
		return nil, connect.NewError(connect.CodeNotFound, errors.New("vehicle not found"))
	}

	res := &pb.GetVehicleResponse{
		Vehicle: vehicle,
	}

	if err := res.Validate(); err != nil {
		log.WithError(err).Error("Invalid response")
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	log.Info("Vehicle found")
	return connect.NewResponse(res), nil
}
