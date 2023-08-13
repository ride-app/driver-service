import (
	"context"
	"errors"
	"strings"

	"github.com/bufbuild/connect-go"
	pb "github.com/ride-app/driver-service/api/gen/ride/driver/v1alpha1"
	"google.golang.org/protobuf/types/known/timestamppb"
)

import (
	"context"
	"errors"
	"strings"

	"cloud.google.com/go/firestore"
	"github.com/ride-app/driver-service/internal/models"
	"github.com/ride-app/driver-service/internal/utils"
	"github.com/sirupsen/logrus"
)

type DriverRepository interface {
	GetDriver(ctx context.Context, log *logrus.Entry, uid string) (*models.Driver, error)
	CreateDriver(ctx context.Context, log *logrus.Entry, driver *models.Driver) (*string, error)
	GoOnline(ctx context.Context, log *logrus.Entry, uid string, vehicle *models.Vehicle) (*models.Status, error)
	UpdateLocation(ctx context.Context, log *logrus.Entry, uid string, location *models.Location) (*string, error)
}

type FirebaseDriverRepository struct {
	client *firestore.Client
}

func NewFirebaseDriverRepository(client *firestore.Client) *FirebaseDriverRepository {
	return &FirebaseDriverRepository{
		client: client,
	}
}

func (repo *FirebaseDriverRepository) GetDriver(ctx context.Context, log *logrus.Entry, uid string) (*models.Driver, error) {
 	log.Info("getting driver")

	dsnap, err := repo.client.Collection("drivers").Doc(uid).Get(ctx)

	if err != nil {
-		log.WithError(err).Error("failed to get driver")
+		log.WithError(err).Error(modifyErrorMessage("failed to get driver"))

		if strings.Contains(err.Error(), "rpc error: code = NotFound") {
			return nil, nil
		}

		return nil, err
	}

	driver := &models.Driver{}
	err = dsnap.DataTo(driver)

	if err != nil {
-		log.WithError(err).Error("failed to convert driver data")
+		log.WithError(err).Error(modifyErrorMessage("failed to convert driver data"))

		return nil, err
	}

	return driver, nil
}

func (repo *FirebaseDriverRepository) CreateDriver(ctx context.Context, log *logrus.Entry, driver *models.Driver) (*string, error) {
 	log.Info("creating driver")

	_, err := repo.client.Collection("drivers").Doc(driver.UID).Set(ctx, driver)

	if err != nil {
-		log.WithError(err).Error("failed to create driver")
+		log.WithError(err).Error(modifyErrorMessage("failed to create driver"))

		return nil, err
	}

	createTime := utils.GetTimestampString()

	return &createTime, nil
}

func (repo *FirebaseDriverRepository) GoOnline(ctx context.Context, log *logrus.Entry, uid string, vehicle *models.Vehicle) (*models.Status, error) {
 	log.Info("going online")

	status := &models.Status{
		Online: true,
	}

	_, err := repo.client.Collection("drivers").Doc(uid).Set(ctx, map[string]interface{}{
		"status": status,
	}, firestore.Merge([]string{"status"}))

	if err != nil {
-		log.WithError(err).Error("failed to go online")
+		log.WithError(err).Error(modifyErrorMessage("failed to go online"))

		return nil, err
	}

	return status, nil
}

func (repo *FirebaseDriverRepository) UpdateLocation(ctx context.Context, log *logrus.Entry, uid string, location *models.Location) (*string, error) {
 	log.Info("updating location")

	_, err := repo.client.Collection("drivers").Doc(uid).Set(ctx, map[string]interface{}{
		"location": location,
	}, firestore.Merge([]string{"location"}))

	if err != nil {
-		log.WithError(err).Error("failed to update location")
+		log.WithError(err).Error(modifyErrorMessage("failed to update location"))

		return nil, err
	}

	updateTime := utils.GetTimestampString()

	return &updateTime, nil
}
