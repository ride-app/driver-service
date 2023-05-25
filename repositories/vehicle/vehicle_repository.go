//go:generate go run github.com/golang/mock/mockgen -destination ../../mocks/$GOFILE -package mocks . VehicleRepository

package vehiclerepository

import (
	"context"
	"errors"
	"fmt"
	"strings"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/bufbuild/connect-go"
	pb "github.com/ride-app/driver-service/api/gen/ride/driver/v1alpha1"
)

type VehicleRepository interface {
	GetVehicle(ctx context.Context, id string) (*pb.Vehicle, error)
	UpdateVehicle(ctx context.Context, vehicle *pb.Vehicle) (updateTime *timestamppb.Timestamp, err error)
}

type FirebaseImpl struct {
	firestore *firestore.Client
}

func NewFirebaseVehicleRepository(firebaseApp *firebase.App) (*FirebaseImpl, error) {
	firestore, err := firebaseApp.Firestore(context.Background())

	if err != nil {
		return nil, err
	}

	return &FirebaseImpl{
		firestore: firestore,
	}, nil
}

func (r *FirebaseImpl) GetVehicle(ctx context.Context, id string) (*pb.Vehicle, error) {
	doc, err := r.firestore.Collection("vehicles").Doc(id).Get(ctx)

	if err != nil {
		return nil, err
	}

	var vehicleType pb.Vehicle_Type

	switch doc.Data()["type"] {
	case strings.ToLower(strings.Split(pb.Vehicle_TYPE_AUTORICKSHAW.String(), "_")[2]):
		vehicleType = pb.Vehicle_TYPE_AUTORICKSHAW
	case strings.ToLower(strings.Split(pb.Vehicle_TYPE_ERICKSHAW.String(), "_")[2]):
		vehicleType = pb.Vehicle_TYPE_ERICKSHAW
	case strings.ToLower(strings.Split(pb.Vehicle_TYPE_MOTORCYCLE.String(), "_")[2]):
		vehicleType = pb.Vehicle_TYPE_MOTORCYCLE
	default:
		vehicleType = pb.Vehicle_TYPE_UNSPECIFIED
	}

	if vehicleType == pb.Vehicle_TYPE_UNSPECIFIED {
		return nil, connect.NewError(connect.CodeFailedPrecondition, errors.New("unknown vehicle type"))
	}
	// Hardcode e-rickshaw for now
	vehicle := pb.Vehicle{
		Name:         "drivers/" + id + "/vehicle",
		Type:         vehicleType,
		DisplayName:  "E-Rickshaw",
		CreateTime:   timestamppb.New(doc.CreateTime),
		UpdateTime:   timestamppb.New(doc.UpdateTime),
		LicensePlate: fmt.Sprintf("%v", doc.Data()["license_plate"]),
	}

	return &vehicle, nil
}

func (r *FirebaseImpl) UpdateVehicle(ctx context.Context, vehicle *pb.Vehicle) (updateTime *timestamppb.Timestamp, err error) {
	x, err := r.firestore.Collection("vehicles").Doc(strings.Split(vehicle.Name, "/")[1]).Set(ctx, map[string]interface{}{
		"license_plate": vehicle.LicensePlate,
		"type":          strings.ToLower(strings.Split(vehicle.Type.String(), "_")[1]),
	})

	if err != nil {
		return nil, err
	}

	return timestamppb.New(x.UpdateTime), nil
}
