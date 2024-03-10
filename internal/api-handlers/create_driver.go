package apihandlers

import (
	"context"
	"errors"
	"strings"

	"connectrpc.com/connect"
	"github.com/bufbuild/protovalidate-go"
	pb "github.com/ride-app/driver-service/api/ride/driver/v1alpha1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (service *DriverServiceServer) CreateDriver(ctx context.Context,
	req *connect.Request[pb.CreateDriverRequest],
) (*connect.Response[pb.CreateDriverResponse], error) {
	log := service.logger.WithFields(map[string]string{
		"method": "CreateDriver",
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

	uid := strings.Split(req.Msg.Driver.Name, "/")[1]

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

	if driver != nil {
		log.Info("Driver already exists")

		return connect.NewResponse(
			&pb.CreateDriverResponse{
				Driver: driver,
			},
		), nil
	}

	createTime, err := service.driverRepository.CreateDriver(ctx, log, req.Msg.Driver)
	if err != nil {
		log.WithError(err).Error("Failed to create driver")
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	res := connect.NewResponse(&pb.CreateDriverResponse{
		Driver: req.Msg.Driver,
	})

	res.Msg.Driver.CreateTime = timestamppb.New(*createTime)
	res.Msg.Driver.UpdateTime = timestamppb.New(*createTime)

	if err := validator.Validate(res.Msg); err != nil {
		log.WithError(err).Error("Invalid response")
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	log.Info("Driver created")
	return res, nil
}
