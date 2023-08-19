package service_test

import (
	"context"

	"connectrpc.com/connect"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	pb "github.com/ride-app/driver-service/api/gen/ride/driver/v1alpha1"
	driverService "github.com/ride-app/driver-service/api/service"
	"github.com/ride-app/driver-service/testing/mocks"
	"go.uber.org/mock/gomock"
)

var _ = Describe("GoOffline", func() {
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
		mockDriverRepo.EXPECT().GoOffline(gomock.Any(), gomock.Any(), gomock.Any()).Return(&pb.Status{}, nil)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	It("should go offline successfully", func() {
		req := connect.NewRequest(&pb.GoOfflineRequest{
			Name: "drivers/valid-driver-id",
		})

		_, err := service.GoOffline(context.Background(), req)
		Expect(err).To(BeNil())
	})
})
