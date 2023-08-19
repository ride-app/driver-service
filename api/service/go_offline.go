package service

import (
	"context"
	"errors"
	"strings"

	"github.com/bufbuild/connect-go"
	pb "github.com/ride-app/driver-service/api/gen/ride/driver/v1alpha1"
)

func (service *DriverServiceServer) GoOffline(ctx context.Context,
	req *connect.Request[pb.GoOfflineRequest]) (*connect.Response[pb.GoOfflineResponse], error) {
	log := service.logger.WithFields(map[string]string{
		"method": "GoOffline",
	})

	if err := req.Msg.Validate(); err != nil {
		log.WithError(err).Info("Invalid request")
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	uid := strings.Split(req.Msg.Name, "/")[1]

	log.Info("uid: ", uid)
	log.Debug("Request header uid: ", req.Header().Get("uid"))

	if uid != req.Header().Get("uid") {
		log.Info("Permission denied")
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("permission denied"))
	}

	status, err := service.driverRepository.GoOffline(ctx, log, uid)

	log.Info("Status: ", status)

	if err != nil {
		log.Error("Failed to go offline")
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := &pb.GoOfflineResponse{
		Status: status,
	}

	if err := res.Validate(); err != nil {
		log.Error("Invalid response")
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	log.Info("Driver went offline")
	return connect.NewResponse(res), nil
}
