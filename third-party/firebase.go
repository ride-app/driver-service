package thirdparty

import (
	"context"

	firebase "firebase.google.com/go/v4"
	"github.com/dragonfish/go/v2/pkg/logger"
	"github.com/ride-app/driver-service/config"
)

func NewFirebaseApp(log logger.Logger, config *config.Config) (*firebase.App, error) {
	ctx := context.Background()
	conf := &firebase.Config{ProjectID: config.ProjectID}
	app, err := firebase.NewApp(ctx, conf)

	if err != nil {
		log.WithError(err).Fatal("Cannot initialize firebase app")
		return nil, err
	}

	return app, nil
}
