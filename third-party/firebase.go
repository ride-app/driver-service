package thirdparty

import (
	"context"

	firebase "firebase.google.com/go/v4"
	"github.com/ride-app/driver-service/config"
	"github.com/ride-app/driver-service/logger"
)

func NewFirebaseApp(log logger.Logger) (*firebase.App, error) {
	ctx := context.Background()
	conf := &firebase.Config{ProjectID: config.Env.ProjectID}
	app, err := firebase.NewApp(ctx, conf)

	if err != nil {
		log.WithError(err).Fatal("Cannot initialize firebase app")
		return nil, err
	}

	return app, nil
}
