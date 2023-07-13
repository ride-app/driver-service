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
	if err := req.Msg.Validate(); err != nil {
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	driverId := strings.Split(req.Msg.Name, "/")[1]

	notificationToken := req.Msg.NotificationToken

	wallet, err := service.walletrepository.GetWallet(ctx, driverId, req.Header().Get("Authorization"))

	if err != nil {
		return nil, err
	}

	if wallet.Balance <= 0 {
		return nil, connect.NewError(connect.CodeFailedPrecondition, errors.New("insufficient wallet balance"))
	}

	vehicle, err := service.vehicleRepository.GetVehicle(ctx, driverId)

	if err != nil {
		return nil, err
	}

	if vehicle == nil {
		return nil, connect.NewError(connect.CodeFailedPrecondition, errors.New("vehicle not found"))
	}

	status, err := service.driverRepository.GoOnline(ctx, driverId, vehicle, notificationToken)

	if err != nil {
		return nil, err
	}

	updateTime, err := service.driverRepository.UpdateLocation(ctx, driverId, req.Msg.Location)

	if err != nil {
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	status.UpdateTime = timestamppb.New(*updateTime)

	return connect.NewResponse(&pb.GoOnlineResponse{
		Status: status,
	}), nil
}
