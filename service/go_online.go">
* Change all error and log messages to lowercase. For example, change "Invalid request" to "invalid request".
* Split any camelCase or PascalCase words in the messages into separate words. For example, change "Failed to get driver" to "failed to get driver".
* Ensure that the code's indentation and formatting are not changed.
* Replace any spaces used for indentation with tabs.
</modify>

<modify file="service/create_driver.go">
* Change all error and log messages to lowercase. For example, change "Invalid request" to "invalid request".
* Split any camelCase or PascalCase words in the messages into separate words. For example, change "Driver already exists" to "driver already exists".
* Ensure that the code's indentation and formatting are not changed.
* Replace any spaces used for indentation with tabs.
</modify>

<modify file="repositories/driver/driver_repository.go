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
	"google.golang.org/genproto/googleapis/type/phone_number"
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
 	log.witherror(err).error("error initializing firestore client")
		return nil, err
	}
 
	auth, err := firebaseApp.Auth(context.Background())
 
	if err != nil {
 	log.witherror(err).error("error initializing auth client")
		return nil, err
	}
 
 	log.info("firebase driver repository initialized")
	return &FirebaseImpl{
		firestore: firestore,
		auth:      auth,
	}, nil
}
 
func (r *FirebaseImpl) CreateDriver(ctx context.Context, log logger.Logger, driver *pb.Driver) (createTime *time.Time, err error) {
 	log.info("updating driver in auth")
 	_, err = r.auth.updateUser(ctx, strings.split(driver.name, "/")[1], (&auth.userToUpdate{}).displayName(driver.displayName).photoURL(driver.photoUri))
 
 	if err != nil {
 	log.witherror(err).error("error updating driver in auth")
 		return nil, err
 	}
 
 	log.info("creating driver in firestore")
	writeResult, err := r.firestore.Collection("drivers").Doc(strings.Split(driver.Name, "/")[1]).Create(ctx, map[string]interface{}{
		"dateOfBirth": map[string]int32{
			"day":   driver.DateOfBirth.Day,
...

