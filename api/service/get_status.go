package service

// import (
// 	"context"
// 	"errors"
// 	"strings"

// 	"connectrpc.com/connect"
// 	pb "github.com/ride-app/driver-service/api/gen/ride/driver/v1alpha1"
// )

// func (service *DriverServiceServer) GetStatus(ctx context.Context,
// 	req *connect.Request[pb.GetStatusRequest]) (*connect.Response[pb.GetStatusResponse], error) {
// 	log := service.logger.WithFields(map[string]string{
// 		"method": "GetStatus",
// 	})

// 	if err := req.Msg.Validate(); err != nil {
// 		log.WithError(err).Info("Invalid request")
// 		return nil, connect.NewError(connect.CodeInvalidArgument, err)
// 	}

// 	uid := strings.Split(req.Msg.Name, "/")[1]

// 	log.Debug("uid: ", uid)
// 	log.Debug("Request header uid: ", req.Header().Get("uid"))

// 	if uid != req.Header().Get("uid") {
// 		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("permission denied"))
// 	}

// 	status, err := service.driverRepository.GetStatus(ctx, log, uid)

// 	if err != nil {
// 		log.WithError(err).Error("Failed to get status")
// 		return nil, connect.NewError(connect.CodeInternal, err)
// 	}

// 	if status == nil {
// 		log.Info("Status not found")
// 		return nil, connect.NewError(connect.CodeNotFound, errors.New("status not found"))
// 	}

// 	res := &pb.GetStatusResponse{
// 		Status: status,
// 	}

// 	if err := res.Validate(); err != nil {
// 		log.WithError(err).Error("Invalid response")
// 		return nil, connect.NewError(connect.CodeInternal, err)
// 	}

// 	log.Info("Status found")
// 	return connect.NewResponse(res), nil
// }
