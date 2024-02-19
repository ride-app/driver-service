package apihandlers_test

import (
	"context"

	"connectrpc.com/connect"
	mock_logger "github.com/dragonfish/go/v2/pkg/logger/mock"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	pb "github.com/ride-app/driver-service/api/ride/driver/v1alpha1"
	apihandlers "github.com/ride-app/driver-service/internal/api-handlers"
	mock_driver "github.com/ride-app/driver-service/internal/repositories/driver/mock"
	mock_vehicle "github.com/ride-app/driver-service/internal/repositories/vehicle/mock"
	mock_wallet "github.com/ride-app/driver-service/internal/repositories/wallet/mock"
	"go.uber.org/mock/gomock"
)

var _ = Describe("GetVehicle", func() {
	var (
		ctrl            *gomock.Controller
		mockDriverRepo  *mock_driver.MockDriverRepository
		mockVehicleRepo *mock_vehicle.MockVehicleRepository
		mockWalletRepo  *mock_wallet.MockWalletRepository
		mockLogger      *mock_logger.MockLogger
		service         *apihandlers.DriverServiceServer
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockDriverRepo = mock_driver.NewMockDriverRepository(ctrl)
		mockVehicleRepo = mock_vehicle.NewMockVehicleRepository(ctrl)
		mockWalletRepo = mock_wallet.NewMockWalletRepository(ctrl)
		mockLogger = &mock_logger.MockLogger{}
		service = apihandlers.New(mockDriverRepo, mockVehicleRepo, mockWalletRepo, mockLogger)
	})

	JustBeforeEach(func() {
		mockVehicleRepo.EXPECT().GetVehicle(gomock.Any(), gomock.Any(), gomock.Any()).Return(&pb.Vehicle{}, nil)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	It("should get the vehicle successfully", func() {
		req := connect.NewRequest(&pb.GetVehicleRequest{
			Name: "drivers/valid-driver-id/vehicles/valid-vehicle-id",
		})

		_, err := service.GetVehicle(context.Background(), req)
		Expect(err).To(BeNil())
	})
})
