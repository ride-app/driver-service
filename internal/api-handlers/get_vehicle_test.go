package apihandlers_test

import (
	"context"

	"connectrpc.com/connect"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	pb "github.com/ride-app/driver-service/api/ride/driver/v1alpha1"
	driverService "github.com/ride-app/driver-service/internal/api-handlers"
	"github.com/ride-app/driver-service/testing/mocks"
	"go.uber.org/mock/gomock"
)

var _ = Describe("GetVehicle", func() {
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
