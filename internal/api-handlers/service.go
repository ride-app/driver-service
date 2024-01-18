package apihandlers

import (
	dr "github.com/ride-app/driver-service/internal/repositories/driver"
	vr "github.com/ride-app/driver-service/internal/repositories/vehicle"
	wr "github.com/ride-app/driver-service/internal/repositories/wallet"
	"github.com/ride-app/go/pkg/logger"
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
