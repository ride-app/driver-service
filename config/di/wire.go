//go:build wireinject

package di

import (
	"github.com/google/wire"
	"github.com/ride-app/driver-service/api/service"
	"github.com/ride-app/driver-service/config"
	driverrepository "github.com/ride-app/driver-service/repositories/driver"
	vehiclerepository "github.com/ride-app/driver-service/repositories/vehicle"
	walletrepository "github.com/ride-app/driver-service/repositories/wallet"
	thirdparty "github.com/ride-app/driver-service/third-party"
	"github.com/ride-app/driver-service/utils/logger"
)

func InitializeService(logger logger.Logger, config *config.Config) (*service.DriverServiceServer, error) {
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
			service.New,
		),
	)
}
