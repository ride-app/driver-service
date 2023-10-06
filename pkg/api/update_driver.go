package api

import (
	"context"
	"errors"
	"strings"

	"connectrpc.com/connect"
	"github.com/bufbuild/protovalidate-go"
	pb "github.com/ride-app/driver-service/pkg/protos/ride/driver/v1alpha1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (service *DriverServiceServer) UpdateDriver(ctx context.Context,
	req *connect.Request[pb.UpdateDriverRequest]) (*connect.Response[pb.UpdateDriverResponse], error) {
	log := service.logger.WithFields(map[string]string{
		"method": "UpdateDriver",
	})

	validator, err := protovalidate.New()
	if err != nil {
		log.WithError(err).Info("Failed to initialize validator")

		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if err := validator.Validate(req.Msg); err != nil {
		log.WithError(err).Info("Invalid request")

		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	if req.Msg.Driver.Name == "" {
		log.Info("Driver name is empty")
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("name cannot be empty"))
	}

	uid := strings.Split(req.Msg.Driver.Name, "/")[1]

	log.Debug("uid: ", uid)
	log.Debug("Request header uid: ", req.Header().Get("uid"))

	if uid != req.Header().Get("uid") {
		log.Info("Permission denied")
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("permission denied"))
	}

	updateTime, err := service.driverRepository.UpdateDriver(ctx, log, req.Msg.Driver)

	if err != nil {
		log.WithError(err).Error("Failed to update driver")
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	req.Msg.Driver.UpdateTime = timestamppb.New(*updateTime)

	res := &pb.UpdateDriverResponse{
		Driver: req.Msg.Driver,
	}

	if err := validator.Validate(res); err != nil {
		log.WithError(err).Error("Invalid response")
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	log.Info("Driver updated")
	return connect.NewResponse(res), nil
}
