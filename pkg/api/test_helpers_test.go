package api_test

import (
	"errors"
	"time"

	pb "github.com/ride-app/driver-service/pkg/protos/ride/driver/v1alpha1"
	"github.com/ride-app/driver-service/pkg/testing/mocks"
	"go.uber.org/mock/gomock"
	"google.golang.org/genproto/googleapis/type/date"
	"google.golang.org/protobuf/reflect/protoreflect"
	"google.golang.org/protobuf/types/known/timestamppb"
)

var MockDriver = &pb.Driver{
	Name:        "drivers/valid-driver-id",
	DisplayName: "John Doe",
	PhotoUri:    "https://example.com/photo.jpg",
	PhoneNumber: "+911234567890",
	DateOfBirth: &date.Date{
		Year:  2000,
		Month: 1,
		Day:   1,
	},
	Gender:     pb.Driver_GENDER_MALE,
	CreateTime: timestamppb.New(time.Now().Add(-(time.Hour * 10))),
	UpdateTime: timestamppb.Now(),
}

var MockVehicle = &pb.Vehicle{
	Type:         pb.Vehicle_TYPE_ERICKSHAW,
	LicensePlate: "WB281234",
	CreateTime:   timestamppb.New(time.Now().Add(-(time.Hour * 10))),
	UpdateTime:   timestamppb.Now(),
}

var MockStatus = &pb.Status{
	Name:       "drivers/valid-driver-id/status",
	Online:     true,
	UpdateTime: timestamppb.Now(),
}

func SetupStubs(
	mockDriverRepo *mocks.MockDriverRepository,
	mockVehicleRepo *mocks.MockVehicleRepository,
	mockWalletRepo *mocks.MockWalletRepository,
	mockLogger *mocks.MockLogger,
) {
	t := time.Now().UTC()
	mockDriverRepo.EXPECT().CreateDriver(gomock.Any(), gomock.Any(), ProtoFieldMatcher[*pb.Driver]("name", "drivers/error-driver-id")).Return(nil, errors.New("error")).AnyTimes()
	mockDriverRepo.EXPECT().CreateDriver(gomock.Any(), gomock.Any(), gomock.Not(ProtoFieldMatcher[*pb.Driver]("name", "drivers/valid-driver-id"))).Return(nil, nil).AnyTimes()
	mockDriverRepo.EXPECT().CreateDriver(gomock.Any(), gomock.Any(), ProtoFieldMatcher[*pb.Driver]("name", "drivers/valid-driver-id")).Return(&t, nil).AnyTimes()

	mockDriverRepo.EXPECT().GetDriver(gomock.Any(), gomock.Any(), gomock.Eq("error-driver-id")).Return(nil, errors.New("error")).AnyTimes()
	mockDriverRepo.EXPECT().GetDriver(gomock.Any(), gomock.Any(), gomock.Not(gomock.Eq("valid-driver-id"))).Return(nil, nil).AnyTimes()
	mockDriverRepo.EXPECT().GetDriver(gomock.Any(), gomock.Any(), gomock.Eq("valid-driver-id")).Return(MockDriver, nil).AnyTimes()

	mockDriverRepo.EXPECT().GoOnline(gomock.Any(), gomock.Any(), gomock.Eq("error-driver-id"), gomock.Any()).Return(nil, errors.New("error")).AnyTimes()
	mockDriverRepo.EXPECT().GoOnline(gomock.Any(), gomock.Any(), gomock.Not(gomock.Eq("valid-driver-id")), gomock.Any()).Return(nil, nil).AnyTimes()
	mockDriverRepo.EXPECT().GoOnline(gomock.Any(), gomock.Any(), gomock.Eq("valid-driver-id"), gomock.Any()).Return(MockStatus, nil).AnyTimes()

	mockDriverRepo.EXPECT().UpdateDriver(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
	mockDriverRepo.EXPECT().DeleteDriver(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()

	mockVehicleRepo.EXPECT().GetVehicle(gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()

	mockWalletRepo.EXPECT().GetWallet(gomock.Any(), gomock.Any(), gomock.Any(), gomock.Any()).Return(nil, nil).AnyTimes()
}

type protoFieldMatcher[T protoreflect.ProtoMessage] struct {
	field string
	value string
}

func (p protoFieldMatcher[T]) Matches(x interface{}) bool {
	return p.value == x.(T).ProtoReflect().Get(x.(T).ProtoReflect().Descriptor().Fields().ByTextName("name")).String()
}

func (p protoFieldMatcher[T]) String() string {
	return "has the same field value"
}

func ProtoFieldMatcher[T protoreflect.ProtoMessage](field string, value string) gomock.Matcher {
	return protoFieldMatcher[T]{field: field, value: value}
}
