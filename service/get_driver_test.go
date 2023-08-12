package service_test

import (
	"context"
	"errors"

	"github.com/bufbuild/connect-go"
	. "github.com/onsi/ginkgo/v2"
	. "github.com/onsi/gomega"
	pb "github.com/ride-app/driver-service/api/gen/ride/driver/v1alpha1"
	"github.com/ride-app/driver-service/mocks"
	driverService "github.com/ride-app/driver-service/service"
	"go.uber.org/mock/gomock"
)

var _ = Describe("GetDriver", func() {
	var (
		ctrl           *gomock.Controller
		mockDriverRepo *mocks.MockDriverRepository
		service        *driverService.DriverServiceServer
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockDriverRepo = mocks.NewMockDriverRepository(ctrl)
		service = &driverService.DriverServiceServer{
			driverRepository: mockDriverRepo,
		}
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("when the request is valid", func() {
		It("returns the driver", func() {
			req := &connect.Request[pb.GetDriverRequest]{
				Msg: &pb.GetDriverRequest{
					Name: "drivers/1",
				},
			}
			mockDriverRepo.EXPECT().GetDriver(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)

			_, err := service.GetDriver(context.Background(), req)
			Expect(err).To(BeNil())
		})
	})

	Context("when the request is invalid", func() {
		It("returns an error", func() {
			req := &connect.Request[pb.GetDriverRequest]{
				Msg: &pb.GetDriverRequest{
					Name: "",
				},
			}
			mockDriverRepo.EXPECT().GetDriver(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("invalid request"))

			_, err := service.GetDriver(context.Background(), req)
			Expect(err).To(MatchError("invalid request"))
		})
	})
})
