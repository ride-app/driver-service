package service

import (
	dr "github.com/ride-app/driver-service/repositories/driver"
	vr "github.com/ride-app/driver-service/repositories/vehicle"
	wr "github.com/ride-app/driver-service/repositories/wallet"
)

type DriverServiceServer struct {
	driverRepository  dr.DriverRepository
	vehicleRepository vr.VehicleRepository
	walletrepository  wr.WalletRepository
}

func New(
	driverRepository dr.DriverRepository,
	vehicleRepository vr.VehicleRepository,
	walletrepository wr.WalletRepository,
) *DriverServiceServer {
	return &DriverServiceServer{
		driverRepository:  driverRepository,
		vehicleRepository: vehicleRepository,
		walletrepository:  walletrepository,
	}
}
