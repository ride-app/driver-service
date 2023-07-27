package service

import (
	"context"
	"errors"
	"strings"

	"github.com/bufbuild/connect-go"
	pb "github.com/ride-app/driver-service/api/gen/ride/driver/v1alpha1"
	log "github.com/sirupsen/logrus"
)

func (service *DriverServiceServer) GoOffline(ctx context.Context,
	req *connect.Request[pb.GoOfflineRequest]) (*connect.Response[pb.GoOfflineResponse], error) {
	if err := req.Msg.Validate(); err != nil {
		log.Info("Invalid request: ", err)
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	uid := strings.Split(req.Msg.Name, "/")[1]

	log.Info("uid: ", uid)
	log.Debug("Request header uid: ", req.Header().Get("uid"))

	if uid != req.Header().Get("uid") {
		log.Info("Permission denied")
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("permission denied"))
	}

	status, err := service.driverRepository.GoOffline(ctx, uid)

	log.Info("Status: ", status)

	if err != nil {
		log.Error("Failed to go offline: ", err)
		return nil, err
	}

	res := &pb.GoOfflineResponse{
		Status: status,
	}

	if err := res.Validate(); err != nil {
		log.Error("Invalid response: ", err)
		return nil, err
	}

	log.Info("Driver went offline")
	return connect.NewResponse(res), nil
}
