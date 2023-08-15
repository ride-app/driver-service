package service_test

import (
	"context"

	"github.com/bufbuild/connect-go"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	pb "github.com/ride-app/driver-service/api/gen/ride/driver/v1alpha1"
	"github.com/ride-app/driver-service/mocks"
	driverService "github.com/ride-app/driver-service/service"
	"go.uber.org/mock/gomock"
)

var _ = Describe("GetLocation", func() {
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
 		mockDriverRepo.EXPECT().GetLocation(gomock.Any(), gomock.Any(), gomock.Any()).Return(&pb.Location{}, nil)
 	})

	AfterEach(func() {
		ctrl.Finish()
	})

 	It("should get the location successfully", func() {
 		location, err := service.GetLocation(context.Background(), &pb.GetLocationRequest{Id: "test-id"})
 		Expect(err).To(BeNil())
 		Expect(location).To(Equal(&pb.Location{}))
 	})
})

