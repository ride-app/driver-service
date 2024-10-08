// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/ride-app/driver-service/internal/repositories/vehicle (interfaces: VehicleRepository)
//
// Generated by this command:
//
//	mockgen -destination ./mock/vehicle_repository.go . VehicleRepository
//

// Package mock_vehicle is a generated GoMock package.
package mock_vehicle

import (
	context "context"
	reflect "reflect"

	logger "github.com/dragonfish/go/v2/pkg/logger"
	v1alpha1 "github.com/ride-app/driver-service/api/ride/driver/v1alpha1"
	gomock "go.uber.org/mock/gomock"
	timestamppb "google.golang.org/protobuf/types/known/timestamppb"
)

// MockVehicleRepository is a mock of VehicleRepository interface.
type MockVehicleRepository struct {
	ctrl     *gomock.Controller
	recorder *MockVehicleRepositoryMockRecorder
}

// MockVehicleRepositoryMockRecorder is the mock recorder for MockVehicleRepository.
type MockVehicleRepositoryMockRecorder struct {
	mock *MockVehicleRepository
}

// NewMockVehicleRepository creates a new mock instance.
func NewMockVehicleRepository(ctrl *gomock.Controller) *MockVehicleRepository {
	mock := &MockVehicleRepository{ctrl: ctrl}
	mock.recorder = &MockVehicleRepositoryMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockVehicleRepository) EXPECT() *MockVehicleRepositoryMockRecorder {
	return m.recorder
}

// GetVehicle mocks base method.
func (m *MockVehicleRepository) GetVehicle(arg0 context.Context, arg1 logger.Logger, arg2 string) (*v1alpha1.Vehicle, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetVehicle", arg0, arg1, arg2)
	ret0, _ := ret[0].(*v1alpha1.Vehicle)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetVehicle indicates an expected call of GetVehicle.
func (mr *MockVehicleRepositoryMockRecorder) GetVehicle(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetVehicle", reflect.TypeOf((*MockVehicleRepository)(nil).GetVehicle), arg0, arg1, arg2)
}

// UpdateVehicle mocks base method.
func (m *MockVehicleRepository) UpdateVehicle(arg0 context.Context, arg1 logger.Logger, arg2 *v1alpha1.Vehicle) (*timestamppb.Timestamp, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UpdateVehicle", arg0, arg1, arg2)
	ret0, _ := ret[0].(*timestamppb.Timestamp)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// UpdateVehicle indicates an expected call of UpdateVehicle.
func (mr *MockVehicleRepositoryMockRecorder) UpdateVehicle(arg0, arg1, arg2 any) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UpdateVehicle", reflect.TypeOf((*MockVehicleRepository)(nil).UpdateVehicle), arg0, arg1, arg2)
}
