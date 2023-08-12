package service

import (
	"context"
	"errors"

	"github.com/bufbuild/connect-go"
	pb "github.com/ride-app/driver-service/api/gen/ride/driver/v1alpha1"
	"github.com/golang/mock/gomock"
	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GoOnline", func() {
	var (
		ctrl *gomock.Controller
		mockDriverRepo *MockDriverRepository
		service *DriverServiceServer
	)

	BeforeEach(func() {
		ctrl = gomock.NewController(GinkgoT())
		mockDriverRepo = NewMockDriverRepository(ctrl)
		service = &DriverServiceServer{
			driverRepository: mockDriverRepo,
		}
	})

	AfterEach(func() {
		ctrl.Finish()
	})

	Context("when the request is valid", func() {
		It("returns the driver", func() {
			req := &connect.Request[pb.GoOnlineRequest]{
				Msg: &pb.GoOnlineRequest{
					Name: "drivers/1",
				},
			}
			mockDriverRepo.EXPECT().GoOnline(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)

			_, err := service.GoOnline(context.Background(), req)
			Expect(err).To(BeNil())
		})
	})

	Context("when the request is invalid", func() {
		It("returns an error", func() {
			req := &connect.Request[pb.GoOnlineRequest]{
				Msg: &pb.GoOnlineRequest{
					Name: "",
				},
			}
			mockDriverRepo.EXPECT().GoOnline(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("invalid request"))

			_, err := service.GoOnline(context.Background(), req)
			Expect(err).To(MatchError("invalid request"))
		})
	})
})