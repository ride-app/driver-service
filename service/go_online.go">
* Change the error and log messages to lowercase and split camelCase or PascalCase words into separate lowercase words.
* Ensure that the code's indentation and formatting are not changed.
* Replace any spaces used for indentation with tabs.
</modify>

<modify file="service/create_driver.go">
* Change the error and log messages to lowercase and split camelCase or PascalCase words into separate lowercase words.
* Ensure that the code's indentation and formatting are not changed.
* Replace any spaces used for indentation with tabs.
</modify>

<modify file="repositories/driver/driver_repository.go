service/go_online.go:

package main

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
	"google.golang.org/genproto/googleapis/type/phone_number"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/timestamppb"
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
		log.WithError(err).Error("Error initializing Firestore client")
		return nil, err
	}

	auth, err := firebaseApp.Auth(context.Background())

	if err != nil {
		log.WithError(err).Error("Error initializing Auth client")
		return nil, err
	}

 log.Info("firebase driver repository initialized")
	return &FirebaseImpl{
		firestore: firestore,
		auth:      auth,
	}, nil
}

func (r *FirebaseImpl) CreateDriver(ctx context.Context, log logger.Logger, driver *pb.Driver) (createTime *time.Time, err error) {
	log.Info("Updating driver in Auth")
	_, err = r.auth.UpdateUser(ctx, strings.Split(driver.Name, "/")[1], (&auth.UserToUpdate{}).DisplayName(driver.DisplayName).PhotoURL(driver.PhotoUri))

	if err != nil {
  log.WithError(err).Error("error updating driver in Auth")
		return nil, err
	}

	log.Info("Creating driver in Firestore")
	writeResult, err := r.firestore.Collection("drivers").Doc(strings.Split(driver.Name, "/")[1]).Create(ctx, map[string]interface{}{
		"dateOfBirth": map[string]int32{
			"day":   driver.DateOfBirth.Day,
			"month": driver.DateOfBirth.Month,
			"year":  driver.DateOfBirth.Year,
		},
		"gender": driver.Gender,
		"phone": map[string]interface{}{
			"number": driver.PhoneNumber.Number,
			"type":   driver.PhoneNumber.Type,
		},
		"createTime": timestamppb.Now(),
		"updateTime": timestamppb.Now(),
	})

	if err != nil {
  log.WithError(err).Error("error creating driver in Firestore")
		return nil, err
	}

	createTime = writeResult.UpdateTime

 log.Info("driver created successfully")
	return createTime, nil
}

