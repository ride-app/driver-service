package service

import (
	dr "github.com/ride-app/driver-service/repositories/driver"
	vr "github.com/ride-app/driver-service/repositories/vehicle"
	wr "github.com/ride-app/driver-service/repositories/wallet"
	"github.com/ride-app/driver-service/utils/logger"
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
