//go:build unit_tests
// +build unit_tests

package service_test

import (
	"context"
	"errors"

	"github.com/bufbuild/connect-go"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	pb "github.com/ride-app/driver-service/api/gen/ride/driver/v1alpha1"
	driverService "github.com/ride-app/driver-service/api/service"
	"github.com/ride-app/driver-service/testing/mocks"
	"go.uber.org/mock/gomock"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ = Describe("GoOnline", func() {
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
		SetupStubs(mockDriverRepo, mockVehicleRepo, mockWalletRepo, mockLogger)
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("given the request is valid", func() {
		var (
			req *connect.Request[pb.GoOnlineRequest]
		)

		BeforeEach(func() {
			req = &connect.Request[pb.GoOnlineRequest]{
				Msg: &pb.GoOnlineRequest{
					Name: "drivers/valid-driver-id",
					Location: &pb.Location{
						Name:      "drivers/valid-driver-id/location",
						Latitude:  12.9716,
						Longitude: 77.5946,
						Timestamp: timestamppb.Now(),
					},
				},
			}

			req.Header().Set("uid", "valid-driver-id")
		})

		When("when the driver is not found", func() {
			BeforeEach(func() {
				req.Msg.Name = "drivers/invalid-driver-id"
				req.Msg.Location.Name = "drivers/invalid-driver-id/location"

				req.Header().Set("uid", "invalid-driver-id")
			})

			It("returns failed precondition error", func() {
				_, err := service.GoOnline(context.Background(), req)

				Expect(err).To(SatisfyAll(
					HaveOccurred(),
					BeAssignableToTypeOf(&connect.Error{}),
					WithTransform(func(err error) connect.Code {
						return err.(*connect.Error).Code()
					}, Equal(connect.CodeFailedPrecondition)),
				))
			})
		})

		When("the driver repository returns error", func() {
			BeforeEach(func() {
				mockDriverRepo.EXPECT().GetDriver(gomock.Any(), gomock.Any(), gomock.Eq("valid-driver-id")).Return(nil, errors.New("error"))
			})

			It("returns internal error", func() {
				_, err := service.GoOnline(context.Background(), req)

				Expect(err).To(SatisfyAll(
					HaveOccurred(),
					BeAssignableToTypeOf(&connect.Error{}),
					WithTransform(func(err error) connect.Code {
						return err.(*connect.Error).Code()
					}, Equal(connect.CodeInternal)),
				))
			})
		})

		// when driver's wallet is not found return failed precondition error
		When("the wallet is not found", func() {
			BeforeEach(func() {
				mockDriverRepo.EXPECT().GetDriver(gomock.Any(), gomock.Any(), gomock.Eq("valid-driver-id")).Return(&pb.Driver{}, nil)
				mockWalletRepo.EXPECT().GetWallet(gomock.Any(), gomock.Any(), gomock.Eq("valid-driver-id"), gomock.Any()).Return(nil, nil)
			})

			It("returns failed precondition error", func() {
				_, err := service.GoOnline(context.Background(), req)

				Expect(err).To(SatisfyAll(
					HaveOccurred(),
					BeAssignableToTypeOf(&connect.Error{}),
					WithTransform(func(err error) connect.Code {
						return err.(*connect.Error).Code()
					}, Equal(connect.CodeFailedPrecondition)),
				))
			})
		})
	})

	Context("given the request is invalid", func() {
	})
})
