package api_test

import (
	"context"

	"connectrpc.com/connect"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	driverService "github.com/ride-app/driver-service/api"
	pb "github.com/ride-app/driver-service/protos/ride/driver/v1alpha1"
	"github.com/ride-app/driver-service/testing/mocks"
	"go.uber.org/mock/gomock"
)

var _ = Describe("UpdateVehicle", func() {
	var (
		ctrl            *gomock.Controller
		mockDriverRepo  *mocks.MockDriverRepository
		mockVehicleRepo *mocks.MockVehicleRepository
		mockWalletRepo  *mocks.MockWalletRepository
		mockLogger      *mocks.MockLogger
		service         *driverService.DriverServiceServer
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockDriverRepo = mocks.NewMockDriverRepository(ctrl)
		mockVehicleRepo = mocks.NewMockVehicleRepository(ctrl)
		mockWalletRepo = mocks.NewMockWalletRepository(ctrl)
		mockLogger = &mocks.MockLogger{}
		service = driverService.New(mockDriverRepo, mockVehicleRepo, mockWalletRepo, mockLogger)
	})

	JustBeforeEach(func() {
		mockVehicleRepo.EXPECT().UpdateVehicle(gomock.Any(), gomock.Any(), gomock.Any()).Return(&pb.Vehicle{}, nil)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	It("should update the vehicle successfully", func() {
		req := connect.NewRequest(&pb.UpdateVehicleRequest{
			Vehicle: &pb.Vehicle{
				Name:         "drivers/valid-driver-id/vehicles/valid-vehicle-id",
				Type:         pb.Vehicle_TYPE_ERICKSHAW,
				DisplayName:  "Erickshaw",
				LicensePlate: "KA-01-1234",
			},
		})

		_, err := service.UpdateVehicle(context.Background(), req)
		Expect(err).To(BeNil())
	})
})
