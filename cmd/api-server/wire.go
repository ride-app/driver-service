//go:build wireinject

package main

import (
	"github.com/google/wire"
	"github.com/ride-app/driver-service/config"
	"github.com/ride-app/driver-service/internal/utils/logger"
	driverrepository "github.com/ride-app/driver-service/pkg/repositories/driver"
	vehiclerepository "github.com/ride-app/driver-service/pkg/repositories/vehicle"
	walletrepository "github.com/ride-app/driver-service/pkg/repositories/wallet"
	thirdparty "github.com/ride-app/driver-service/third-party"
)

func InitializeService(logger logger.Logger, config *config.Config) (*api_handlers.DriverServiceServer, error) {
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
			api_handlers.New,
		),
	)
}
