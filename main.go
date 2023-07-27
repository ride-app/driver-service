package main

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/bufbuild/connect-go"
	"github.com/ilyakaznacheev/cleanenv"
	"github.com/ride-app/driver-service/api/gen/ride/driver/v1alpha1/driverv1alpha1connect"
	"github.com/ride-app/driver-service/config"
	"github.com/ride-app/driver-service/di"
	"github.com/ride-app/driver-service/interceptors"
	"github.com/sirupsen/logrus"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	service, err := di.InitializeService()

	if err != nil {
		logrus.Fatalf("Failed to initialize service: %v", err)
	}

	logrus.Info("Service Initialized")

	// Create a context that, when cancelled, ends the JWKS background refresh goroutine.
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	authInterceptor, err := interceptors.NewAuthInterceptor(ctx)

	if err != nil {
		logrus.Fatalf("Failed to initialize auth interceptor: %v", err)
	}

	connectInterceptors := connect.WithInterceptors(authInterceptor)

	path, handler := driverv1alpha1connect.NewDriverServiceHandler(service, connectInterceptors)
	mux := http.NewServeMux()
	mux.Handle(path, handler)

	panic(http.ListenAndServe(
		fmt.Sprintf("0.0.0.0:%d", config.Env.Port),
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	))

}

func init() {
	logrus.SetReportCaller(true)

	logrus.SetFormatter(&logrus.JSONFormatter{
		FieldMap: logrus.FieldMap{
			logrus.FieldKeyTime:  "timestamp",
			logrus.FieldKeyLevel: "severity",
			logrus.FieldKeyMsg:   "message",
		},
		TimestampFormat: time.RFC3339Nano,
	})

	logrus.SetLevel(logrus.InfoLevel)

	err := cleanenv.ReadEnv(&config.Env)

	if config.Env.Debug {
		logrus.SetLevel(logrus.DebugLevel)
	}

	if err != nil {
		logrus.Warnf("Could not load config: %v", err)
	}
}
