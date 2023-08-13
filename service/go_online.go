package service

import (
	"context"
	"errors"
	"strings"

	"github.com/bufbuild/connect-go"
	pb "github.com/ride-app/driver-service/api/gen/ride/driver/v1alpha1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

func (service *DriverServiceServer) GoOnline(ctx context.Context,
	req *connect.Request[pb.GoOnlineRequest]) (*connect.Response[pb.GoOnlineResponse], error) {
	log := service.logger.WithFields(map[string]string{
		"method": "GoOnline",
	})

	if err := req.Msg.Validate(); err != nil {
  		log.Info("invalid request")
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	uid := strings.Split(req.Msg.Name, "/")[1]

	log.Debug("uid: ", uid)
	log.Debug("Request header uid: ", req.Header().Get("uid"))

	if uid != req.Header().Get("uid") {
  		log.Info("permission denied")
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("permission denied"))
	}

	driver, err := service.driverRepository.GetDriver(ctx, log, uid)

	if err != nil {
  		log.WithError(err).Error("failed to get driver")
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if driver == nil {
  		log.Info("driver not found")
		return nil, connect.NewError(connect.CodeFailedPrecondition, errors.New("driver not found"))
	}

	wallet, err := service.walletrepository.GetWallet(ctx, log, uid, req.Header().Get("Authorization"))

	if err != nil {
  		log.WithError(err).Error("failed to get wallet")
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if wallet.Balance <= 0 {
  		log.Info("insufficient wallet balance: ", wallet.Balance)

		return nil, connect.NewError(connect.CodeFailedPrecondition, errors.New("insufficient wallet balance"))
	}

	vehicle, err := service.vehicleRepository.GetVehicle(ctx, log, uid)

	if err != nil {
  		log.WithError(err).Error("failed to get vehicle")
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	if vehicle == nil {
  		log.Info("vehicle not found")
		return nil, connect.NewError(connect.CodeFailedPrecondition, errors.New("vehicle not found"))
	}

	status, err := service.driverRepository.GoOnline(ctx, log, uid, vehicle)

	if err != nil {
  		log.WithError(err).Error("failed to go online")
		return nil, connect.NewError(connect.CodeInternal, err)
	}

  		log.Info("status: ", status.Online)

	updateTime, err := service.driverRepository.UpdateLocation(ctx, log, uid, req.Msg.Location)

	if err != nil {
  		log.WithError(err).Error("failed to update location")
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	status.UpdateTime = timestamppb.New(*updateTime)

	res := &pb.GoOnlineResponse{
		Status: status,
	}

	if err := res.Validate(); err != nil {
  		log.WithError(err).Error("invalid response")
		return nil, connect.NewError(connect.CodeInternal, err)
	}

  		log.Info("driver is online")
	return connect.NewResponse(res), nil
}
