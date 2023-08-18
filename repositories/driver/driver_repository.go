//go:generate go run go.uber.org/mock/mockgen -destination ../../mocks/$GOFILE -package mocks . DriverRepository

package driverrepository

import (
	"context"
	"errors"
	"strings"
	"time"

	"cloud.google.com/go/firestore"
	firebase "firebase.google.com/go/v4"
	"firebase.google.com/go/v4/auth"
	pb "github.com/ride-app/driver-service/api/gen/ride/driver/v1alpha1"
	"github.com/ride-app/driver-service/logger"
	"google.golang.org/genproto/googleapis/type/date"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"

	"github.com/mmcloughlin/geohash"
)

type DriverRepository interface {
	CreateDriver(ctx context.Context, log logger.Logger, driver *pb.Driver) (createTime *time.Time, err error)

	GetDriver(ctx context.Context, log logger.Logger, id string) (*pb.Driver, error)

	UpdateDriver(ctx context.Context, log logger.Logger, driver *pb.Driver) (createTime *time.Time, err error)

	DeleteDriver(ctx context.Context, log logger.Logger, id string) (createTime *time.Time, err error)

	GetStatus(ctx context.Context, log logger.Logger, id string) (*pb.Status, error)

	GoOnline(ctx context.Context, log logger.Logger, id string, vehicleType *pb.Vehicle) (*pb.Status, error)

	GoOffline(ctx context.Context, log logger.Logger, id string) (*pb.Status, error)

	GetLocation(ctx context.Context, log logger.Logger, id string) (*pb.Location, error)

	UpdateLocation(ctx context.Context, log logger.Logger, id string, location *pb.Location) (updateTime *time.Time, err error)
}

type FirebaseImpl struct {
	firestore *firestore.Client
	auth      *auth.Client
}

func NewFirebaseDriverRepository(firebaseApp *firebase.App, log logger.Logger) (*FirebaseImpl, error) {
	firestore, err := firebaseApp.Firestore(context.Background())

	if err != nil {
		log.WithError(err).Error("Error initializing firestore client")
		return nil, err
	}

	auth, err := firebaseApp.Auth(context.Background())

	if err != nil {
		log.WithError(err).Error("Error initializing auth client")
		return nil, err
	}

	log.Info("Firebase driver repository initialized")
	return &FirebaseImpl{
		firestore: firestore,
		auth:      auth,
	}, nil
}

func (r *FirebaseImpl) CreateDriver(ctx context.Context, log logger.Logger, driver *pb.Driver) (createTime *time.Time, err error) {
	log.Info("Updating driver in auth")
	_, err = r.auth.UpdateUser(ctx, strings.Split(driver.Name, "/")[1], (&auth.UserToUpdate{}).DisplayName(driver.DisplayName).PhotoURL(driver.PhotoUri))

	if err != nil {
		log.WithError(err).Error("Error updating driver in auth")
		return nil, err
	}

	log.Info("Creating driver in firestore")
	writeResult, err := r.firestore.Collection("drivers").Doc(strings.Split(driver.Name, "/")[1]).Create(ctx, map[string]interface{}{
		"dateOfBirth": map[string]int32{
			"day":   driver.DateOfBirth.Day,
			"month": driver.DateOfBirth.Month,
			"year":  driver.DateOfBirth.Year,
		},
		"gender": strings.Split(pb.Driver_Gender_name[int32(driver.Gender.Number())], "_")[1],
	})

	if status.Code(err) == codes.AlreadyExists {
		log.Info("Driver already exists in firestore")

		return nil, errors.New("driver already exists in firestore")
	} else if err != nil {
		log.WithError(err).Error("Error creating driver in firestore")
		return nil, err
	}

	timestamp := writeResult.UpdateTime

	return &timestamp, nil
}

func (r *FirebaseImpl) GetDriver(ctx context.Context, log logger.Logger, id string) (*pb.Driver, error) {
	user, err := r.auth.GetUser(ctx, id)

	if auth.IsUserNotFound(err) {
		log.Info("Driver does not exist in auth")
		return nil, nil
	}

	log.Info("Getting driver from firestore")
	doc, err := r.firestore.Collection("drivers").Doc(id).Get(ctx)

	if status.Code(err) == codes.NotFound {
		log.Info("Driver does not exist in firestore")
		return nil, nil
	} else if err != nil {
		log.WithError(err).Error("Error getting driver from firestore")
		return nil, err
	}

	if !doc.Exists() {
		log.Info("Driver does not exist in firestore")
		return nil, nil
	}

	if err != nil {
		log.WithError(err).Error("Error getting driver from auth")
		return nil, err
	}

	driver := pb.Driver{
		Name:        "drivers/" + id,
		DisplayName: user.DisplayName,
		PhotoUri:    user.PhotoURL,
		PhoneNumber: user.PhoneNumber,
		DateOfBirth: &date.Date{
			Day:   int32(doc.Data()["dateOfBirth"].(map[string]interface{})["day"].(int64)),
			Month: int32(doc.Data()["dateOfBirth"].(map[string]interface{})["month"].(int64)),
			Year:  int32(doc.Data()["dateOfBirth"].(map[string]interface{})["year"].(int64)),
		},
		Gender:     pb.Driver_Gender(pb.Driver_Gender_value["GENDER_"+doc.Data()["gender"].(string)]),
		CreateTime: timestamppb.New(doc.CreateTime),
		UpdateTime: timestamppb.New(doc.UpdateTime),
	}

	return &driver, nil
}

func (r *FirebaseImpl) UpdateDriver(ctx context.Context, log logger.Logger, driver *pb.Driver) (updateTime *time.Time, err error) {
	log.Info("Updating driver in auth")
	_, err = r.auth.UpdateUser(ctx, strings.Split(driver.Name, "/")[1], (&auth.UserToUpdate{}).DisplayName(driver.DisplayName).PhotoURL(driver.PhotoUri).PhoneNumber(driver.PhoneNumber))

	if auth.IsUserNotFound(err) {
		log.Info("Driver does not exist in auth")
		return nil, errors.New("driver does not exist in auth")
	} else if err != nil {
		log.WithError(err).Error("Error updating driver in auth")
		return nil, err
	}

	timestamp := time.Now()

	return &timestamp, nil
}

func (r *FirebaseImpl) DeleteDriver(ctx context.Context, log logger.Logger, id string) (deleteTime *time.Time, err error) {
	log.Info("Deleting driver from firestore")
	writeResult, err := r.firestore.Collection("drivers").Doc(id).Delete(ctx)

	if status.Code(err) == codes.NotFound {
		log.Info("Driver does not exist in firestore")
		return nil, nil
	} else if err != nil {
		log.WithError(err).Error("Error deleting driver from firestore")
		return nil, err
	}

	timestamp := writeResult.UpdateTime

	return &timestamp, nil
}

func (r *FirebaseImpl) GetStatus(ctx context.Context, log logger.Logger, id string) (*pb.Status, error) {
	log.Info("Getting status from firestore")
	doc, err := r.firestore.Collection("activeDrivers").Doc(id).Get(ctx)

	if status.Code(err) == codes.NotFound {
		log.Info("Driver does not exist in firestore")
		return nil, nil
	} else if err != nil {
		log.WithError(err).Error("Error getting status from firestore")
		return nil, err
	}

	if !doc.Exists() {
		log.Info("Driver does not exist in firestore")
		return nil, nil
	}

	status := pb.Status{
		Name:       "drivers/" + id + "/status",
		Online:     doc.Exists(),
		UpdateTime: timestamppb.New(doc.UpdateTime),
	}

	return &status, nil
}

func (r *FirebaseImpl) GoOnline(ctx context.Context, log logger.Logger, id string, vehicle *pb.Vehicle) (*pb.Status, error) {
	log.Info("Updating active driver in firestore")

	ref := r.firestore.Collection("activeDrivers").Doc(id)
	err := r.firestore.RunTransaction(ctx, func(ctx context.Context, tx *firestore.Transaction) error {
		doc, err := tx.Get(ref)

		if err != nil && !(status.Code(err) == codes.NotFound) {
			log.WithError(err).Error("Error getting active driver from firestore")
			return err
		}

		if doc.Exists() {
			return nil
		}

		err = tx.Set(ref, map[string]interface{}{
			"vehicleId":    strings.Split(vehicle.Name, "/")[1],
			"licensePlate": vehicle.LicensePlate,
			"vehicleType":  strings.ToLower(vehicle.Type.String()),
			"capacity":     4,
		})

		if err != nil {
			log.WithError(err).Error("Error setting active driver in firestore")
			return err
		}

		return nil
	})

	if err != nil {
		log.WithError(err).Error("Error updating active driver in firestore")
		return nil, err
	}

	return &pb.Status{
		Name:       "drivers/" + id + "/status",
		Online:     true,
		UpdateTime: timestamppb.Now(),
	}, nil
}

func (r *FirebaseImpl) GoOffline(ctx context.Context, log logger.Logger, id string) (*pb.Status, error) {
	log.Info("Deleting active driver from firestore")
	_, err := r.firestore.Collection("activeDrivers").Doc(id).Delete(ctx)

	if status.Code(err) == codes.NotFound {
		log.Info("Driver does not exist in active drivers in firestore")
	} else if err != nil {
		log.WithError(err).Error("Error deleting active driver from firestore")
		return nil, err
	}

	return &pb.Status{
		Name:       "drivers/" + id + "/status",
		Online:     false,
		UpdateTime: timestamppb.Now(),
	}, nil
}

func (r *FirebaseImpl) GetLocation(ctx context.Context, log logger.Logger, id string) (*pb.Location, error) {
	log.Info("Checking if driver is active in firestore")
	doc, err := r.firestore.Collection("activeDrivers").Doc(id).Get(ctx)

	if status.Code(err) == codes.NotFound {
		log.Info("Driver does not exist in active drivers in firestore")
		return nil, nil
	} else if err != nil {
		log.WithError(err).Error("Error checking if driver is active in firestore")
		return nil, err
	}

	if !doc.Exists() {
		log.Info("Driver is not active in firestore")
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

func (r *FirebaseImpl) UpdateLocation(ctx context.Context, log logger.Logger, id string, location *pb.Location) (updateTime *time.Time, err error) {

	log.Info("Calculating geohash")
	hash := geohash.Encode(location.Latitude, location.Longitude)

	log.Info("Updating driver location in firestore")
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
		log.WithError(err).Error("Error updating driver location in firestore")
		return nil, err
	}

	return &res.UpdateTime, nil
}
