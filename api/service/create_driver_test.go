//go:build unit_tests
// +build unit_tests

package service_test

import (
	"context"
	"errors"
	"time"

	"connectrpc.com/connect"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	driverService "github.com/ride-app/driver-service/api/service"
	pb "github.com/ride-app/driver-service/proto/ride/driver/v1alpha1"
	"github.com/ride-app/driver-service/testing/mocks"
	"go.uber.org/mock/gomock"
	"google.golang.org/genproto/googleapis/type/date"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var _ = Describe("CreateDriver", func() {
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

	Context("given request is valid", func() {
		var req *connect.Request[pb.CreateDriverRequest]

		BeforeEach(func() {
			req = &connect.Request[pb.CreateDriverRequest]{
				Msg: &pb.CreateDriverRequest{
					Driver: &pb.Driver{
						Name:        "drivers/valid-driver-id",
						DisplayName: "Jane Doe",
						PhotoUri:    "https://example.com/photo.jpg",
						PhoneNumber: "+911234567890",
						DateOfBirth: &date.Date{
							Year:  2000,
							Month: 1,
							Day:   1,
						},
						Gender: pb.Driver_GENDER_FEMALE,
					},
				},
			}
			req.Header().Set("uid", "valid-driver-id")
		})

		Context("and driver already exists", func() {
			BeforeEach(func() {
				mockDriverRepo.EXPECT().GetDriver(gomock.Any(), gomock.Any(), gomock.Eq("valid-driver-id")).Return(MockDriver, nil)
			})

			It("returns response with driver", func() {
				res, err := service.CreateDriver(context.Background(), req)

				Expect(err).ToNot(HaveOccurred())
				Expect(res.Msg.Driver.Name).To(Equal(req.Msg.Driver.Name))
				Expect(proto.Equal(req.Msg.Driver, res.Msg.Driver)).To(BeFalse())
			})

			When("create time is passed with the request", func() {
				var createTime *time.Time

				BeforeEach(func() {
					t := time.Now()
					createTime = &t
					req.Msg.Driver.CreateTime = timestamppb.New(t)
				})

				It("ignores createTime and returns different createTime", func() {
					res, err := service.CreateDriver(context.Background(), req)

					Expect(err).ToNot(HaveOccurred())
					Expect(res.Msg.Driver.CreateTime).To(Not(Equal(timestamppb.New(*createTime))))
				})
			})

			When("updateTime is passed with the request", func() {
				var updateTime *time.Time

				BeforeEach(func() {
					t := time.Now()
					updateTime = &t
					req.Msg.Driver.CreateTime = timestamppb.New(t)
				})

				It("ignores updateTime and returns different updateTime", func() {
					res, err := service.CreateDriver(context.Background(), req)

					Expect(err).ToNot(HaveOccurred())
					Expect(res.Msg.Driver.UpdateTime).To(Not(Equal(timestamppb.New(*updateTime))))
				})
			})
		})

		Context("and driver does not exist", func() {
			BeforeEach(func() {
				mockDriverRepo.EXPECT().GetDriver(gomock.Any(), gomock.Any(), gomock.Eq("valid-driver-id")).Return(nil, nil).AnyTimes()
			})

			It("returns valid response with same driver", func() {
				res, err := service.CreateDriver(context.Background(), req)

				Expect(err).ToNot(HaveOccurred())
				Expect(res.Msg.Driver).To(SatisfyAll(
					Not(BeNil()),
					BeAssignableToTypeOf(&pb.Driver{}),
				))
			})

			It("returns driver with createTime and updateTime", func() {
				res, err := service.CreateDriver(context.Background(), req)

				Expect(err).ToNot(HaveOccurred())
				Expect(res.Msg.Driver).To(SatisfyAll(
					Not(BeNil()),
					BeAssignableToTypeOf(&pb.Driver{}),
				))
				Expect(res.Msg.Driver.CreateTime).To(Not(BeZero()))
			})

			It("returns response where createTime and updateTime is equal", func() {
				res, err := service.CreateDriver(context.Background(), req)

				Expect(err).ToNot(HaveOccurred())
				Expect(res.Msg.Driver.CreateTime).To(Equal(res.Msg.Driver.UpdateTime))
			})

			When("create time is passed with the request", func() {
				var createTime *time.Time

				BeforeEach(func() {
					t := time.Now()
					createTime = &t
					req.Msg.Driver.CreateTime = timestamppb.New(t)
				})

				It("ignores createTime and returns different createTime", func() {
					res, err := service.CreateDriver(context.Background(), req)

					Expect(err).ToNot(HaveOccurred())
					Expect(res.Msg.Driver.CreateTime).To(Not(Equal(timestamppb.New(*createTime))))
				})
			})

			When("updateTime is passed with the request", func() {
				var updateTime *time.Time

				BeforeEach(func() {
					t := time.Now()
					updateTime = &t
					req.Msg.Driver.CreateTime = timestamppb.New(t)
				})

				It("ignores updateTime and returns different updateTime", func() {
					res, err := service.CreateDriver(context.Background(), req)

					Expect(err).ToNot(HaveOccurred())
					Expect(res.Msg.Driver.UpdateTime).To(Not(Equal(timestamppb.New(*updateTime))))
				})
			})

			When("driver repository CreateDriver returns error", func() {
				BeforeEach(func() {
					req.Msg.Driver.Name = "drivers/error-driver-id"
					req.Header().Set("uid", "error-driver-id")
				})

				It("returns internal error", func() {
					_, err := service.CreateDriver(context.Background(), req)
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

		When("driver repository GetDriver returns error", func() {
			BeforeEach(func() {
				mockDriverRepo.EXPECT().GetDriver(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
			})

			It("returns internal error", func() {
				_, err := service.CreateDriver(context.Background(), req)
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
		var req *connect.Request[pb.CreateDriverRequest]

		BeforeEach(func() {
			req = &connect.Request[pb.CreateDriverRequest]{
				Msg: &pb.CreateDriverRequest{
					Driver: &pb.Driver{
						Name:        "drivers/valid-driver-id",
						DisplayName: "John Doe",
						PhotoUri:    "https://example.com/photo.jpg",
						PhoneNumber: "+911234567890",
						DateOfBirth: &date.Date{
							Year:  2000,
							Month: 1,
							Day:   1,
						},
						Gender: pb.Driver_GENDER_MALE,
					},
				},
			}

			req.Header().Set("uid", "valid-driver-id")
		})

		When("driver name is empty", func() {
			It("returns invalid argument error", func() {
				req.Msg.Driver.Name = ""

				_, err := service.CreateDriver(context.Background(), req)
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
				req.Msg.Driver.Name = "not-drivers/invalid-driver-id"

				_, err := service.CreateDriver(context.Background(), req)
				Expect(err).To(SatisfyAll(
					HaveOccurred(),
					BeAssignableToTypeOf(&connect.Error{}),
					WithTransform(func(err error) connect.Code {
						return err.(*connect.Error).Code()
					}, Equal(connect.CodeInvalidArgument)),
				))
			})
		})

		When("driver display name is empty", func() {
			It("does not return error", func() {
				req.Msg.Driver.DisplayName = ""

				_, err := service.CreateDriver(context.Background(), req)
				Expect(err).To(
					Not(HaveOccurred()),
				)
			})
		})

		When("photo uri is empty", func() {
			It("returns invalid argument error", func() {
				req.Msg.Driver.PhotoUri = ""

				_, err := service.CreateDriver(context.Background(), req)
				Expect(err).To(SatisfyAll(
					HaveOccurred(),
					BeAssignableToTypeOf(&connect.Error{}),
					WithTransform(func(err error) connect.Code {
						return err.(*connect.Error).Code()
					}, Equal(connect.CodeInvalidArgument)),
				))
			})
		})

		When("photo uri is not a valid uri", func() {
			It("returns invalid argument error", func() {
				req.Msg.Driver.PhotoUri = "invalid-uri"

				_, err := service.CreateDriver(context.Background(), req)
				Expect(err).To(SatisfyAll(
					HaveOccurred(),
					BeAssignableToTypeOf(&connect.Error{}),
					WithTransform(func(err error) connect.Code {
						return err.(*connect.Error).Code()
					}, Equal(connect.CodeInvalidArgument)),
				))
			})
		})

		When("phone number is empty", func() {
			It("returns invalid argument error", func() {
				req.Msg.Driver.PhoneNumber = ""

				_, err := service.CreateDriver(context.Background(), req)
				Expect(err).To(SatisfyAll(
					HaveOccurred(),
					BeAssignableToTypeOf(&connect.Error{}),
					WithTransform(func(err error) connect.Code {
						return err.(*connect.Error).Code()
					}, Equal(connect.CodeInvalidArgument)),
				))
			})
		})

		When("phone number is not a valid phone number", func() {
			It("returns invalid argument error", func() {
				req.Msg.Driver.PhoneNumber = "invalid-phone-number"

				_, err := service.CreateDriver(context.Background(), req)
				Expect(err).To(SatisfyAll(
					HaveOccurred(),
					BeAssignableToTypeOf(&connect.Error{}),
					WithTransform(func(err error) connect.Code {
						return err.(*connect.Error).Code()
					}, Equal(connect.CodeInvalidArgument)),
				))
			})
		})

		When("date of birth is empty", func() {
			It("returns invalid argument error", func() {
				req.Msg.Driver.DateOfBirth = nil

				_, err := service.CreateDriver(context.Background(), req)
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

				_, err := service.CreateDriver(context.Background(), req)
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

				_, err := service.CreateDriver(context.Background(), req)
				Expect(err).To(SatisfyAll(
					HaveOccurred(),
					BeAssignableToTypeOf(&connect.Error{}),
					WithTransform(func(err error) connect.Code {
						return err.(*connect.Error).Code()
					}, Equal(connect.CodePermissionDenied)),
				))
			})
		})

		When("driver gender is unspecified", func() {
			It("returns invalid argument error", func() {
				req.Msg.Driver.Gender = pb.Driver_GENDER_UNSPECIFIED

				_, err := service.CreateDriver(context.Background(), req)
				Expect(err).To(SatisfyAll(
					HaveOccurred(),
					BeAssignableToTypeOf(&connect.Error{}),
					WithTransform(func(err error) connect.Code {
						return err.(*connect.Error).Code()
					}, Equal(connect.CodeInvalidArgument)),
				))
			})
		})
	})
})
