//go:build unit_tests
// +build unit_tests

package service_test

// var _ = Describe("GoOnline", func() {
// 	var (
// 		ctrl            *gomock.Controller
// 		mockDriverRepo  *mocks.MockDriverRepository
// 		mockVehicleRepo *mocks.MockVehicleRepository
// 		mockWalletRepo  *mocks.MockWalletRepository
// 		mockLogger      *mocks.MockLogger
// 		service         *driverService.DriverServiceServer
// 	)

// 	BeforeEach(func() {
// 		ctrl = gomock.NewController(GinkgoT())
// 		mockDriverRepo = mocks.NewMockDriverRepository(ctrl)
// 		mockVehicleRepo = mocks.NewMockVehicleRepository(ctrl)
// 		mockWalletRepo = mocks.NewMockWalletRepository(ctrl)
// 		mockLogger = &mocks.MockLogger{}
// 		service = driverService.New(mockDriverRepo, mockVehicleRepo, mockWalletRepo, mockLogger)
// 	})

// 	AfterEach(func() {
// 		ctrl.Finish()
// 	})

// 	Context("when the request is valid", func() {
// 		It("returns the driver", func() {
// 			req := &connect.Request[pb.GoOnlineRequest]{
// 				Msg: &pb.GoOnlineRequest{
// 					Name: "drivers/1",
// 				},
// 			}
// 			mockDriverRepo.EXPECT().GoOnline(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil)

// 			_, err := service.GoOnline(context.Background(), req)
// 			Expect(err).To(BeNil())
// 		})
// 	})

// 	Context("when the request is invalid", func() {
// 		It("returns an error", func() {
// 			req := &connect.Request[pb.GoOnlineRequest]{
// 				Msg: &pb.GoOnlineRequest{
// 					Name: "",
// 				},
// 			}
// 			mockDriverRepo.EXPECT().GoOnline(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, errors.New("invalid request"))

// 			_, err := service.GoOnline(context.Background(), req)
// 			Expect(err).To(MatchError("invalid request"))
// 		})
// 	})
// })
