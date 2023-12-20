package api_handlers

import (
	"context"
	"errors"
	"strings"

	"connectrpc.com/connect"
	"github.com/bufbuild/protovalidate-go"
	pb "github.com/ride-app/driver-service/api/ride/driver/v1alpha1"
)

func (service *DriverServiceServer) GetDriver(ctx context.Context,
	req *connect.Request[pb.GetDriverRequest]) (*connect.Response[pb.GetDriverResponse], error) {
	log := service.logger.WithFields(map[string]string{
		"method": "GetDriver",
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

	uid := strings.Split(req.Msg.Name, "/")[1]

	log.Debug("uid: ", uid)
	log.Debug("Request header uid: ", req.Header().Get("uid"))

	if uid != req.Header().Get("uid") {
		log.Info("Permission denied")
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("permission denied"))
	}

	driver, err := service.driverRepository.GetDriver(ctx, log, uid)

	if err != nil {
		log.WithError(err).Error("Failed to get driver")
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if driver == nil {
		log.Info("Driver not found")
		return nil, connect.NewError(connect.CodeNotFound, errors.New("driver not found"))
	}

	res := &pb.GetDriverResponse{
		Driver: driver,
	}

	if err := validator.Validate(res); err != nil {
		log.WithError(err).Error("Invalid response")
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	log.Info("Driver found")
	return connect.NewResponse(res), nil
}
