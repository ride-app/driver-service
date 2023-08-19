package service

import (
	"context"
	"errors"
	"strings"

	"connectrpc.com/connect"
	pb "github.com/ride-app/driver-service/api/gen/ride/driver/v1alpha1"
)

func (service *DriverServiceServer) UpdateLocation(ctx context.Context,
	req *connect.Request[pb.UpdateLocationRequest]) (*connect.Response[pb.UpdateLocationResponse], error) {
	log := service.logger.WithFields(map[string]string{
		"method": "UpdateLocation",
	})

	if err := req.Msg.Validate(); err != nil {
		log.WithError(err).Info("Invalid request")
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	uid := strings.Split(req.Msg.Parent, "/")[1]

	log.Debug("uid: ", uid)
	log.Debug("Request header uid: ", req.Header().Get("uid"))

	if uid != req.Header().Get("uid") {
		log.Info("Permission denied")
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("permission denied"))
	}

	status, err := service.driverRepository.GetStatus(ctx, log, uid)

	if err != nil {
		log.WithError(err).Error("Failed to get status")
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if !status.Online {
		log.Info("Driver is offline")
		return nil, connect.NewError(connect.CodeFailedPrecondition, err)
	}

	_, err = service.driverRepository.UpdateLocation(ctx, log, uid, req.Msg.Location)

	if err != nil {
		log.WithError(err).Error("Failed to update location")
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	log.Info("Location updated")
	return connect.NewResponse(&pb.UpdateLocationResponse{}), nil
}
