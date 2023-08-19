package main

import (
	"context"
	"fmt"
	"net/http"

	"github.com/bufbuild/connect-go"
	"github.com/ride-app/driver-service/api/gen/ride/driver/v1alpha1/driverv1alpha1connect"
	"github.com/ride-app/driver-service/api/interceptors"
	"github.com/ride-app/driver-service/config"
	"github.com/ride-app/driver-service/config/di"
	"github.com/ride-app/driver-service/utils/logger"
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

func main() {
	config, err := config.New()

	log := logger.New(config)

	if err != nil {
		log.WithError(err).Fatal("Failed to read environment variables")
	}

	service, err := di.InitializeService(log, config)

	if err != nil {
		log.Fatalf("Failed to initialize service: %v", err)
	}

	log.Info("Service Initialized")

	// Create a context that, when cancelled, ends the JWKS background refresh goroutine.
	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	authInterceptor, err := interceptors.NewAuthInterceptor(ctx, log)

	if err != nil {
		log.Fatalf("Failed to initialize auth interceptor: %v", err)
	}

	connectInterceptors := connect.WithInterceptors(authInterceptor)

	path, handler := driverv1alpha1connect.NewDriverServiceHandler(service, connectInterceptors)
	mux := http.NewServeMux()
	mux.Handle(path, handler)

	// trunk-ignore(semgrep/go.lang.security.audit.net.use-tls.use-tls)
	panic(http.ListenAndServe(
		fmt.Sprintf("0.0.0.0:%d", config.Port),
		// Use h2c so we can serve HTTP/2 without TLS.
		h2c.NewHandler(mux, &http2.Server{}),
	))

}
