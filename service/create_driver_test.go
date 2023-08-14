//go:build unit_tests
// +build unit_tests

package service_test

import (
	"context"
	"errors"
	"time"

	"github.com/bufbuild/connect-go"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	pb "github.com/ride-app/driver-service/api/gen/ride/driver/v1alpha1"
	"github.com/ride-app/driver-service/mocks"
	driverService "github.com/ride-app/driver-service/service"
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
		mockDriverRepo.EXPECT().CreateDriver(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
		mockDriverRepo.EXPECT().GetDriver(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
		mockDriverRepo.EXPECT().UpdateDriver(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
		mockDriverRepo.EXPECT().DeleteDriver(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()

		mockVehicleRepo.EXPECT().GetVehicle(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
		mockWalletRepo.EXPECT().GetWallet(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
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

		When("driver already exists", func() {
			BeforeEach(func() {
				expectedDriver := &pb.Driver{
					Name:        "drivers/valid-driver-id",
					DisplayName: "Jane Doe",
					PhotoUri:    "https://example.com/photo.jpg",
					PhoneNumber: "+910987654321",
					DateOfBirth: &date.Date{
						Year:  2001,
						Month: 1,
						Day:   1,
					},
					Gender:     pb.Driver_GENDER_FEMALE,
					CreateTime: timestamppb.Now(),
					UpdateTime: timestamppb.New(time.Now().Add(-(time.Hour * 10))),
				}
				mockDriverRepo.EXPECT().GetDriver(gomock.Any(), gomock.Any(), gomock.Any()).Return(expectedDriver, nil)
			})

			It("returns response with driver", func() {
				res, err := service.CreateDriver(context.Background(), req)

				Expect(err).ToNot(HaveOccurred())
				Expect(res.Msg.Driver.Name).To(Equal(req.Msg.Driver.Name))
				Expect(proto.Equal(req.Msg.Driver, res.Msg.Driver)).To(BeFalse())
			})
		})

		When("driver does not exist", func() {
			var createTime *time.Time

			BeforeEach(OncePerOrdered, func() {
				t := time.Now().Add(time.Second * 10)
				createTime = &t

				mockDriverRepo.EXPECT().CreateDriver(gomock.Any(), gomock.Any(), gomock.Eq(req.Msg.Driver)).Return(createTime, nil)
			})

			AfterEach(func() {
				createTime = nil
			})

			It("returns valid response with same driver", func() {
				res, err := service.CreateDriver(context.Background(), req)

				Expect(err).ToNot(HaveOccurred())
				Expect(res.Msg.Driver).To(SatisfyAll(
					Not(BeNil()),
					BeAssignableToTypeOf(&pb.Driver{}),
				))
				Expect(res.Msg.Validate()).To(Succeed())
			})

			It("returns driver with updated createTime and updateTime", func() {
				res, err := service.CreateDriver(context.Background(), req)

				req.Msg.Driver.CreateTime = timestamppb.New(*createTime)
				req.Msg.Driver.UpdateTime = timestamppb.New(*createTime)

				Expect(err).ToNot(HaveOccurred())
				Expect(res.Msg.Driver).To(SatisfyAll(
					Not(BeNil()),
					BeAssignableToTypeOf(&pb.Driver{}),
				))
				Expect(res.Msg.Validate()).To(Succeed())
				Expect(proto.Equal(req.Msg.Driver, res.Msg.Driver)).To(BeTrue())
			})

			It("returns response where createTime and updateTime is equal", func() {
				res, err := service.CreateDriver(context.Background(), req)

				req.Msg.Driver.CreateTime = timestamppb.New(*createTime)
				req.Msg.Driver.UpdateTime = timestamppb.New(*createTime)

				Expect(err).ToNot(HaveOccurred())
				Expect(res.Msg.Driver.CreateTime).To(Equal(res.Msg.Driver.UpdateTime))
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

		When("driver repository CreateDriver returns error", func() {
			BeforeEach(func() {
				mockDriverRepo.EXPECT().CreateDriver(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("error"))
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
			BeforeEach(func() {
				createtime := time.Now()
				mockDriverRepo.EXPECT().CreateDriver(gomock.Any(), gomock.Any(), gomock.Any()).Return(&createtime, nil)
			})

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

		When("create time is nil", func() {
			var createTime *time.Time

			BeforeEach(func() {
				t := time.Now()
				createTime = &t
				mockDriverRepo.EXPECT().CreateDriver(gomock.Any(), gomock.Any(), gomock.Any()).Return(createTime, nil)
			})

			AfterEach(func() {
				createTime = nil
			})

			It("returns response with a createTime field", func() {
				res, err := service.CreateDriver(context.Background(), req)

				Expect(err).ToNot(HaveOccurred())
				Expect(res.Msg.Driver.CreateTime).To(Equal(timestamppb.New(*createTime)))
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
