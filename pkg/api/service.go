package api

import (
	dr "github.com/ride-app/driver-service/pkg/repositories/driver"
	vr "github.com/ride-app/driver-service/pkg/repositories/vehicle"
	wr "github.com/ride-app/driver-service/pkg/repositories/wallet"
	"github.com/ride-app/driver-service/pkg/utils/logger"
)

type DriverServiceServer struct {
	driverRepository  dr.DriverRepository
	vehicleRepository vr.VehicleRepository
	walletrepository  wr.WalletRepository
	logger            logger.Logger
}

func New(
	driverRepository dr.DriverRepository,
	vehicleRepository vr.VehicleRepository,
	walletrepository wr.WalletRepository,
	logger logger.Logger,
) *DriverServiceServer {
	return &DriverServiceServer{
		driverRepository:  driverRepository,
		vehicleRepository: vehicleRepository,
		walletrepository:  walletrepository,
		logger:            logger,
	}
}
