package thirdparty

import (
	"context"
	"log"

	firebase "firebase.google.com/go"
	"github.com/ride-app/driver-service/config"
)

func NewFirebaseApp() (*firebase.App, error) {
	ctx := context.Background()
	conf := &firebase.Config{ProjectID: config.Env.Firebase_Project_Id}
	app, err := firebase.NewApp(ctx, conf)

	if err != nil {
		log.Fatalln(err)
		return nil, err
	}

	return app, nil
}
