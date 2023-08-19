package service

// import (
// 	"context"
// 	"errors"
// 	"strings"

// 	"connectrpc.com/connect"
// 	pb "github.com/ride-app/driver-service/api/gen/ride/driver/v1alpha1"
// )

// func (service *DriverServiceServer) GetLocation(ctx context.Context,
// 	req *connect.Request[pb.GetLocationRequest]) (*connect.Response[pb.GetLocationResponse], error) {
// 	log := service.logger.WithFields(map[string]string{
// 		"method": "GetLocation",
// 	})

// 	if err := req.Msg.Validate(); err != nil {
// 		log.WithError(err).Info("Invalid request")
// 		return nil, connect.NewError(connect.CodeInvalidArgument, err)
// 	}

// 	uid := strings.Split(req.Msg.Name, "/")[1]

// 	log.Debug("uid: ", uid)
// 	log.Debug("Request header uid: ", req.Header().Get("uid"))

// 	if uid != req.Header().Get("uid") {
// 		log.Info("Permission denied")
// 		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("permission denied"))
// 	}

// 	location, err := service.driverRepository.GetLocation(ctx, log, uid)

// 	if err != nil {
// 		log.WithError(err).Error("Failed to get location")
// 		return nil, connect.NewError(connect.CodeInternal, err)
// 	}

// 	if location == nil {
// 		log.Info("Location not found")
// 		return nil, connect.NewError(connect.CodeNotFound, errors.New("location not found"))
// 	}

// 	res := &pb.GetLocationResponse{
// 		Location: location,
// 	}

// 	if err := res.Validate(); err != nil {
// 		log.WithError(err).Error("Invalid response")
// 		return nil, connect.NewError(connect.CodeInternal, err)
// 	}

// 	log.Info("Location found")
// 	return connect.NewResponse(res), nil
// }
