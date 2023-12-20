package apihandlers

import (
	"context"
	"errors"
	"strings"

	"connectrpc.com/connect"
	"github.com/bufbuild/protovalidate-go"
	pb "github.com/ride-app/driver-service/api/ride/driver/v1alpha1"
)

func (service *DriverServiceServer) DeleteDriver(ctx context.Context,
	req *connect.Request[pb.DeleteDriverRequest]) (*connect.Response[pb.DeleteDriverResponse], error) {
	log := service.logger.WithFields(map[string]string{
		"method": "DeleteDriver",
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

	status, err := service.driverRepository.GetStatus(ctx, log, uid)

	if err != nil {
		log.WithError(err).Error("Failed to get driver status")
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if status == nil {
		log.Info("Status not found")
		return nil, connect.NewError(connect.CodeFailedPrecondition, errors.New("driver status unkown"))
	}

	if status.Online {
		log.Info("Driver is online")
		return nil, connect.NewError(connect.CodeFailedPrecondition, errors.New("driver is online"))
	}

	_, err = service.driverRepository.DeleteDriver(ctx, log, uid)

	if err != nil {
		log.WithError(err).Error("Failed to delete driver")
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	log.Info("Driver deleted")
	return connect.NewResponse(&pb.DeleteDriverResponse{}), nil
}
