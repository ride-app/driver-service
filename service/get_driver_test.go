package service_test

import (
	"context"
	"errors"

	"github.com/bufbuild/connect-go"
	pb "github.com/ride-app/driver-service/api/gen/ride/driver/v1alpha1"
	"go.uber.org/mock/mockgen"
	"github.com/ride-app/driver-service/mocks"
	. "github.com/onsi/ginkgo/v2/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("GetDriver", func() {
 var (
 	ctrl *gomock.Controller
 	mockDriverRepo *MockDriverRepository
 	service *DriverServiceServer
 )
 
 BeforeEach(func() {
 	ctrl = gomock.NewController(GinkgoT())
 	mockDriverRepo = mocks.NewMockDriverRepository(ctrl)
 	service = &DriverServiceServer{
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
      mockDriverRepo.EXPECT().GetDriver(mock.Anything, mock.Anything, mock.Anything).Return(nil, nil)

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
      mockDriverRepo.EXPECT().GetDriver(mock.Anything, mock.Anything, mock.Anything).Return(nil, errors.New("invalid request"))

			_, err := service.GetDriver(context.Background(), req)
			Expect(err).To(MatchError("invalid request"))
		})
	})
})