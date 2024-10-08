//go:build unit_tests
// +build unit_tests

package apihandlers_test

import (
	"context"
	"errors"

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
	"google.golang.org/genproto/googleapis/type/date"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ = Describe("GetDriver", func() {
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
		mockLogger = &mocks.MockLogger{}
		service = apihandlers.New(mockDriverRepo, mockVehicleRepo, mockWalletRepo, mockLogger)
	})

	JustBeforeEach(func() {
		mockDriverRepo.EXPECT().
			GetDriver(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, nil).
			AnyTimes()
		mockDriverRepo.EXPECT().
			GetDriver(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, nil).
			AnyTimes()
		mockDriverRepo.EXPECT().
			UpdateDriver(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, nil).
			AnyTimes()
		mockDriverRepo.EXPECT().
			DeleteDriver(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, nil).
			AnyTimes()
		mockVehicleRepo.EXPECT().
			GetVehicle(gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, nil).
			AnyTimes()
		mockWalletRepo.EXPECT().
			GetWallet(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).
			Return(nil, nil).
			AnyTimes()
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("given the request is valid", func() {
		var req *connect.Request[pb.GetDriverRequest]

		BeforeEach(func() {
			req = &connect.Request[pb.GetDriverRequest]{
				Msg: &pb.GetDriverRequest{
					Name: "drivers/valid-driver-id",
				},
			}

			req.Header().Set("uid", "valid-driver-id")
		})

		When("the driver does not exist", func() {
			BeforeEach(func() {
				mockDriverRepo.EXPECT().
					GetDriver(gomock.Any(), gomock.Any(), gomock.Eq("valid-driver-id")).
					Return(nil, nil)
			})
			It("returns not found error", func() {
				_, err := service.GetDriver(context.Background(), req)
				Expect(err).To(SatisfyAll(
					HaveOccurred(),
					BeAssignableToTypeOf(&connect.Error{}),
					WithTransform(func(err error) connect.Code {
						return err.(*connect.Error).Code()
					}, Equal(connect.CodeNotFound)),
				))
			})
		})

		When("the driver does exist", func() {
			var driver *pb.Driver

			BeforeEach(func() {
				driver = &pb.Driver{
					Name:        "drivers/valid-driver-id",
					DisplayName: "John Doe",
					PhotoUri:    "https://example.com/photo.jpg",
					PhoneNumber: "+911234567890",
					DateOfBirth: &date.Date{
						Year:  1990,
						Month: 1,
						Day:   1,
					},
					Gender:     pb.Driver_GENDER_MALE,
					CreateTime: timestamppb.Now(),
					UpdateTime: timestamppb.Now(),
				}

				mockDriverRepo.EXPECT().
					GetDriver(gomock.Any(), gomock.Any(), gomock.Eq("valid-driver-id")).
					Return(driver, nil)
			})

			AfterEach(func() {
				driver = nil
			})

			It("returns the driver", func() {
				res, err := service.GetDriver(context.Background(), req)
				Expect(err).To(BeNil())
				Expect(protos.Equal(driver, res.Msg.Driver)).To(BeTrue())
			})
		})

		When("driver repository GetDriver returns error", func() {
			BeforeEach(func() {
				mockDriverRepo.EXPECT().
					GetDriver(gomock.Any(), gomock.Any(), gomock.Eq("valid-driver-id")).
					Return(nil, errors.New("error"))
			})

			It("returns internal error", func() {
				_, err := service.GetDriver(context.Background(), req)
				Expect(err).To(SatisfyAll(
					HaveOccurred(),
					BeAssignableToTypeOf(&connect.Error{}),
					WithTransform(func(err error) connect.Code {
						return err.(*connect.Error).Code()
					}, Equal(connect.CodeInternal)),
				))
			})
		})
	})

	Context("given all other request parameters are valid", func() {
		var req *connect.Request[pb.GetDriverRequest]

		BeforeEach(func() {
			req = &connect.Request[pb.GetDriverRequest]{
				Msg: &pb.GetDriverRequest{
					Name: "drivers/valid-driver-id",
				},
			}
			req.Header().Set("uid", "valid-driver-id")
		})

		When("driver name is empty", func() {
			It("returns invalid argument error", func() {
				req.Msg.Name = ""

				_, err := service.GetDriver(context.Background(), req)
				Expect(err).To(SatisfyAll(
					HaveOccurred(),
					BeAssignableToTypeOf(&connect.Error{}),
					WithTransform(func(err error) connect.Code {
						return err.(*connect.Error).Code()
					}, Equal(connect.CodeInvalidArgument)),
				))
			})
		})

		When("driver name does not match drivers/driverId pattern", func() {
			It("returns invalid argument error", func() {
				req.Msg.Name = "not-drivers/invalid-driver-id"

				_, err := service.GetDriver(context.Background(), req)
				Expect(err).To(SatisfyAll(
					HaveOccurred(),
					BeAssignableToTypeOf(&connect.Error{}),
					WithTransform(func(err error) connect.Code {
						return err.(*connect.Error).Code()
					}, Equal(connect.CodeInvalidArgument)),
				))
			})
		})

		When("uid header is missing in request", func() {
			It("returns permission denied error", func() {
				req.Header().Del("uid")

				_, err := service.GetDriver(context.Background(), req)
				Expect(err).To(SatisfyAll(
					HaveOccurred(),
					BeAssignableToTypeOf(&connect.Error{}),
					WithTransform(func(err error) connect.Code {
						return err.(*connect.Error).Code()
					}, Equal(connect.CodePermissionDenied)),
				))
			})
		})

		When("uid header is not equal to driver id", func() {
			It("returns permission denied error", func() {
				req.Header().Set("uid", "invalid-driver-id")

				_, err := service.GetDriver(context.Background(), req)
				Expect(err).To(SatisfyAll(
					HaveOccurred(),
					BeAssignableToTypeOf(&connect.Error{}),
					WithTransform(func(err error) connect.Code {
						return err.(*connect.Error).Code()
					}, Equal(connect.CodePermissionDenied)),
				))
			})
		})
	})
})
