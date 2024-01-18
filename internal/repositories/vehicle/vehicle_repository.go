//go:generate go run go.uber.org/mock/mockgen -destination ../../../pkg/testing/mocks/$GOFILE -package mocks . VehicleRepository

package vehiclerepository

import (
	"context"
	"errors"
	"strings"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"connectrpc.com/connect"
	pb "github.com/ride-app/driver-service/api/ride/driver/v1alpha1"
	"github.com/ride-app/go/pkg/logger"
)

type VehicleRepository interface {
	GetVehicle(ctx context.Context, log logger.Logger, id string) (*pb.Vehicle, error)
	UpdateVehicle(ctx context.Context, log logger.Logger, vehicle *pb.Vehicle) (updateTime *timestamppb.Timestamp, err error)
}

type FirebaseImpl struct {
	firestore *firestore.Client
}

func NewFirebaseVehicleRepository(firebaseApp *firebase.App, log logger.Logger) (*FirebaseImpl, error) {
	firestore, err := firebaseApp.Firestore(context.Background())

	if err != nil {
		log.WithError(err).Error("Error initializing firestore client")
		return nil, err
	}

	log.Info("Firebase vehicle repository initialized")
	return &FirebaseImpl{
		firestore: firestore,
	}, nil
}

func (r *FirebaseImpl) GetVehicle(ctx context.Context, log logger.Logger, id string) (*pb.Vehicle, error) {
	log.Info("Getting vehicle from firestore")
	doc, err := r.firestore.Collection("vehicles").Doc(id).Get(ctx)

	if status.Code(err) == codes.NotFound {
		log.Info("Vehicle not found in firestore")
		return nil, nil
	} else if err != nil {
		log.WithError(err).Error("Error getting vehicle from firestore")
		return nil, err
	}

	if !doc.Exists() {
		log.WithError(err).Error("Vehicle does not exist in firestore")
		return nil, nil
	}

	var vehicleType pb.Vehicle_Type

	switch doc.Data()["type"] {
	case strings.Split(pb.Vehicle_TYPE_AUTORICKSHAW.String(), "_")[1]:
		vehicleType = pb.Vehicle_TYPE_AUTORICKSHAW
	case strings.Split(pb.Vehicle_TYPE_ERICKSHAW.String(), "_")[1]:
		vehicleType = pb.Vehicle_TYPE_ERICKSHAW
	case strings.Split(pb.Vehicle_TYPE_MOTORCYCLE.String(), "_")[1]:
		vehicleType = pb.Vehicle_TYPE_MOTORCYCLE
	default:
		vehicleType = pb.Vehicle_TYPE_UNSPECIFIED
	}

	if vehicleType == pb.Vehicle_TYPE_UNSPECIFIED {
		log.WithError(err).Error("Unknown vehicle type")
		return nil, connect.NewError(connect.CodeFailedPrecondition, errors.New("unknown vehicle type"))
	}
	// Hardcode e-rickshaw for now
	vehicle := pb.Vehicle{
		Name:         "drivers/" + id + "/vehicle",
		Type:         vehicleType,
		DisplayName:  doc.Data()["display_name"].(string),
		CreateTime:   timestamppb.New(doc.CreateTime),
		UpdateTime:   timestamppb.New(doc.UpdateTime),
		LicensePlate: doc.Data()["license_plate"].(string),
	}

	return &vehicle, nil
}

func (r *FirebaseImpl) UpdateVehicle(ctx context.Context, log logger.Logger, vehicle *pb.Vehicle) (updateTime *timestamppb.Timestamp, err error) {
	log.Info("Updating vehicle in firestore")
	x, err := r.firestore.Collection("vehicles").Doc(strings.Split(vehicle.Name, "/")[1]).Set(ctx, map[string]interface{}{
		"license_plate": vehicle.LicensePlate,
		"type":          strings.Split(vehicle.Type.String(), "_")[1],
		"display_name":  vehicle.DisplayName,
	})

	if err != nil {
		log.WithError(err).Error("Error updating vehicle in firestore")
		return nil, err
	}

	return timestamppb.New(x.UpdateTime), nil
}
