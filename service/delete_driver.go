package service

import (
	"context"
	"errors"
	"strings"

	"github.com/bufbuild/connect-go"
	pb "github.com/ride-app/driver-service/api/gen/ride/driver/v1alpha1"
	"github.com/sirupsen/logrus"
)

func (service *DriverServiceServer) DeleteDriver(ctx context.Context,
	req *connect.Request[pb.DeleteDriverRequest]) (*connect.Response[pb.DeleteDriverResponse], error) {

	if err := req.Msg.Validate(); err != nil {
		logrus.Info("Invalid request")
		return nil, connect.NewError(connect.CodeInvalidArgument, err)
	}

	if req.Msg.Name == "" {
		logrus.Info("Name is empty")
		return nil, connect.NewError(connect.CodeInvalidArgument, errors.New("name cannot be empty"))
	}

	uid := strings.Split(req.Msg.Name, "/")[1]

	logrus.Info("uid: ", uid)
	logrus.Debug("Request header uid: ", req.Header().Get("uid"))

	if uid != req.Header().Get("uid") {
		logrus.Info("Permission denied")
		return nil, connect.NewError(connect.CodePermissionDenied, errors.New("permission denied"))
	}

	_, err := service.driverRepository.DeleteDriver(ctx, uid)

	if err != nil {
		logrus.WithError(err).Error("Failed to delete driver")
		return nil, connect.NewError(connect.CodeInternal, err)
	}

	logrus.Info("Driver deleted")
	return connect.NewResponse(&pb.DeleteDriverResponse{}), nil
}
