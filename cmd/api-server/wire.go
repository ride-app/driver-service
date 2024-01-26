//go:build wireinject

package main

import (
	"github.com/deb-tech-n-sol/go/pkg/logger"
	"github.com/google/wire"
	"github.com/ride-app/driver-service/config"
	apihandlers "github.com/ride-app/driver-service/internal/api-handlers"
	driverrepository "github.com/ride-app/driver-service/internal/repositories/driver"
	vehiclerepository "github.com/ride-app/driver-service/internal/repositories/vehicle"
	walletrepository "github.com/ride-app/driver-service/internal/repositories/wallet"
	thirdparty "github.com/ride-app/driver-service/third-party"
)

func InitializeService(logger logger.Logger, config *config.Config) (*apihandlers.DriverServiceServer, error) {
	panic(
		wire.Build(
			thirdparty.NewFirebaseApp,
			driverrepository.NewFirebaseDriverRepository,
			walletrepository.New,
			vehiclerepository.NewFirebaseVehicleRepository,
			wire.Bind(
				new(driverrepository.DriverRepository),
				new(*driverrepository.FirebaseImpl),
			),
			wire.Bind(
				new(walletrepository.WalletRepository),
				new(*walletrepository.Impl),
			),
			wire.Bind(
				new(vehiclerepository.VehicleRepository),
				new(*vehiclerepository.FirebaseImpl),
			),
			apihandlers.New,
		),
	)
}
