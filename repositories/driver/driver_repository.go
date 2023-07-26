//go:generate go run go.uber.org/mock/mockgen -destination ../../mocks/$GOFILE -package mocks . DriverRepository

package driverrepository

import (
	"context"
	"errors"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go"
	"firebase.google.com/go/auth"
	"github.com/bufbuild/connect-go"
	pb "github.com/ride-app/driver-service/api/gen/ride/driver/v1alpha1"
	"google.golang.org/genproto/googleapis/type/phone_number"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/mmcloughlin/geohash"
)

type DriverRepository interface {
	CreateDriver(ctx context.Context, driver *pb.Driver) (createTime *time.Time, err error)

	GetDriver(ctx context.Context, id string) (*pb.Driver, error)

	UpdateDriver(ctx context.Context, driver *pb.Driver) (createTime *time.Time, err error)

	DeleteDriver(ctx context.Context, id string) (createTime *time.Time, err error)

	GetStatus(ctx context.Context, id string) (*pb.Status, error)

	GoOnline(ctx context.Context, id string, vehicleType *pb.Vehicle) (*pb.Status, error)

	GoOffline(ctx context.Context, id string) (*pb.Status, error)

	GetLocation(ctx context.Context, id string) (*pb.Location, error)

	UpdateLocation(ctx context.Context, id string, location *pb.Location) (updateTime *time.Time, err error)
}

type FirebaseImpl struct {
	firestore *firestore.Client
	auth      *auth.Client
}

func NewFirebaseDriverRepository(firebaseApp *firebase.App) (*FirebaseImpl, error) {
	firestore, err := firebaseApp.Firestore(context.Background())

	if err != nil {
		return nil, err
	}

	auth, err := firebaseApp.Auth(context.Background())

	if err != nil {
		return nil, err
	}

	return &FirebaseImpl{
		firestore: firestore,
		auth:      auth,
	}, nil
}

func (r *FirebaseImpl) CreateDriver(ctx context.Context, driver *pb.Driver) (createTime *time.Time, err error) {
	_, err = r.auth.UpdateUser(ctx, strings.Split(driver.Name, "/")[1], (&auth.UserToUpdate{}).DisplayName(driver.DisplayName).PhotoURL(driver.PhotoUri))

	if err != nil {
		return nil, err
	}

	writeResult, err := r.firestore.Collection("drivers").Doc(strings.Split(driver.Name, "/")[1]).Create(ctx, map[string]interface{}{
		"dateOfBirth": map[string]int32{
			"day":   driver.DateOfBirth.Day,
			"month": driver.DateOfBirth.Month,
			"year":  driver.DateOfBirth.Year,
		},
		"gender": pb.Driver_Gender_name[int32(driver.Gender.Number())],
	})

	if err != nil {
		return nil, err
	}

	timestamp := writeResult.UpdateTime

	return &timestamp, nil
}

func (r *FirebaseImpl) GetDriver(ctx context.Context, id string) (*pb.Driver, error) {
	doc, err := r.firestore.Collection("drivers").Doc(id).Get(ctx)

	if err != nil {
		return nil, err
	}

	if !doc.Exists() {
		return nil, nil
	}

	user, err := r.auth.GetUser(ctx, id)

	if err != nil {
		return nil, err
	}

	driver := pb.Driver{
		Name:        "drivers/" + id,
		DisplayName: user.DisplayName,
		PhotoUri:    user.PhotoURL,
		PhoneNumber: &phone_number.PhoneNumber{
			Kind: &phone_number.PhoneNumber_E164Number{
				E164Number: user.PhoneNumber,
			},
		},
		CreateTime: timestamppb.New(doc.CreateTime),
		UpdateTime: timestamppb.New(doc.UpdateTime),
	}

	return &driver, nil
}

func (r *FirebaseImpl) UpdateDriver(ctx context.Context, driver *pb.Driver) (updateTime *time.Time, err error) {
	_, err = r.auth.UpdateUser(ctx, strings.Split(driver.Name, "/")[1], (&auth.UserToUpdate{}).DisplayName(driver.DisplayName).PhotoURL(driver.PhotoUri).PhoneNumber(driver.PhoneNumber.GetE164Number()))

	if err != nil {
		return nil, err
	}

	timestamp := time.Now()

	return &timestamp, nil
}

func (r *FirebaseImpl) DeleteDriver(ctx context.Context, id string) (deleteTime *time.Time, err error) {
	status, err := r.GetStatus(ctx, id)

	if err != nil {
		return nil, err
	}

	if status.Online {
		return nil, connect.NewError(connect.CodeFailedPrecondition, errors.New("driver is online"))
	}

	writeResult, err := r.firestore.Collection("drivers").Doc(id).Delete(ctx)

	if err != nil {
		return nil, err
	}

	timestamp := writeResult.UpdateTime

	return &timestamp, nil
}

func (r *FirebaseImpl) GetStatus(ctx context.Context, id string) (*pb.Status, error) {
	doc, err := r.firestore.Collection("activeDrivers").Doc(id).Get(ctx)

	if err != nil {
		return nil, err
	}

	if !doc.Exists() {
		return nil, nil
	}

	status := pb.Status{
		Name:       "drivers/" + id + "/status",
		Online:     doc.Exists(),
		UpdateTime: timestamppb.New(doc.UpdateTime),
	}

	return &status, nil
}

func (r *FirebaseImpl) GoOnline(ctx context.Context, id string, vehicle *pb.Vehicle) (*pb.Status, error) {
	_, err := r.firestore.Collection("activeDrivers").Doc(id).Set(ctx, map[string]interface{}{
		"vehicleId":    strings.Split(vehicle.Name, "/")[1],
		"licensePlate": vehicle.LicensePlate,
		"vehicleType":  strings.ToLower(vehicle.Type.String()),
		"capacity":     4,
	})

	if err != nil {
		return nil, err
	}

	return &pb.Status{
		Name:       "drivers/" + id + "/status",
		Online:     true,
		UpdateTime: timestamppb.Now(),
	}, nil
}

func (r *FirebaseImpl) GoOffline(ctx context.Context, id string) (*pb.Status, error) {
	_, err := r.firestore.Collection("activeDrivers").Doc(id).Delete(ctx)

	if err != nil {
		return nil, err
	}

	return &pb.Status{
		Name:       "drivers/" + id + "/status",
		Online:     false,
		UpdateTime: timestamppb.Now(),
	}, nil
}

func (r *FirebaseImpl) GetLocation(ctx context.Context, id string) (*pb.Location, error) {
	doc, err := r.firestore.Collection("activeDrivers").Doc(id).Get(ctx)

	if err != nil {
		return nil, err
	}

	if !doc.Exists() {
		return nil, nil
	}

	data := doc.Data()

	location := data["location"].(map[string]interface{})
	latitude := location["latitude"].(float64)
	longitude := location["longitude"].(float64)

	return &pb.Location{
		Latitude:  latitude,
		Longitude: longitude,
	}, nil
}

func (r *FirebaseImpl) UpdateLocation(ctx context.Context, id string, location *pb.Location) (updateTime *time.Time, err error) {

	hash := geohash.Encode(location.Latitude, location.Longitude)

	res, err := r.firestore.Collection("activeDrivers").Doc(id).Update(ctx, []firestore.Update{
		{
			Path:  "location.latitude",
			Value: location.Latitude,
		},
		{
			Path:  "location.longitude",
			Value: location.Longitude,
		},
		{
			Path:  "geohash",
			Value: hash,
		},
	})

	if err != nil {
		return nil, err
	}

	return &res.UpdateTime, nil
}
