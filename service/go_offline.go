package service

import (
	"context"
	"strings"

	"github.com/bufbuild/connect-go"
	pb "github.com/ride-app/driver-service/api/gen/ride/driver/v1alpha1"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

func (service *DriverServiceServer) GoOffline(ctx context.Context,
	req *connect.Request[pb.GoOfflineRequest]) (*connect.Response[pb.GoOfflineResponse], error) {
	if err := req.Msg.Validate(); err != nil {
		return nil, status.Error(codes.InvalidArgument, err.Error())
	}

	driverId := strings.Split(req.Msg.Name, "/")[1]

	status, err := service.driverRepository.GoOffline(ctx, driverId)

	if err != nil {
		return nil, err
	}

	return connect.NewResponse(&pb.GoOfflineResponse{
		Status: status,
	}), nil
}
