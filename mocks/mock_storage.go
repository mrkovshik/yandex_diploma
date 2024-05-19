// Code generated by MockGen. DO NOT EDIT.
// Source: github.com/mrkovshik/yandex_diploma/internal/service (interfaces: Storage)

// Package mock_service is a generated GoMock package.
package mock_service

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	model "github.com/mrkovshik/yandex_diploma/internal/model"
)

// MockStorage is a mock of Storage interface.
type MockStorage struct {
	ctrl     *gomock.Controller
	recorder *MockStorageMockRecorder
}

// MockStorageMockRecorder is the mock recorder for MockStorage.
type MockStorageMockRecorder struct {
	mock *MockStorage
}

// NewMockStorage creates a new mock instance.
func NewMockStorage(ctrl *gomock.Controller) *MockStorage {
	mock := &MockStorage{ctrl: ctrl}
	mock.recorder = &MockStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockStorage) EXPECT() *MockStorageMockRecorder {
	return m.recorder
}

// AddUser mocks base method.
func (m *MockStorage) AddUser(arg0 context.Context, arg1, arg2 string) (uint, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddUser", arg0, arg1, arg2)
	ret0, _ := ret[0].(uint)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// AddUser indicates an expected call of AddUser.
func (mr *MockStorageMockRecorder) AddUser(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddUser", reflect.TypeOf((*MockStorage)(nil).AddUser), arg0, arg1, arg2)
}

// FinalizeOrderAndUpdateBalance mocks base method.
func (m *MockStorage) FinalizeOrderAndUpdateBalance(arg0 context.Context, arg1 string, arg2 float64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "FinalizeOrderAndUpdateBalance", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// FinalizeOrderAndUpdateBalance indicates an expected call of FinalizeOrderAndUpdateBalance.
func (mr *MockStorageMockRecorder) FinalizeOrderAndUpdateBalance(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "FinalizeOrderAndUpdateBalance", reflect.TypeOf((*MockStorage)(nil).FinalizeOrderAndUpdateBalance), arg0, arg1, arg2)
}

// GetOrderByNumber mocks base method.
func (m *MockStorage) GetOrderByNumber(arg0 context.Context, arg1 string) (model.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrderByNumber", arg0, arg1)
	ret0, _ := ret[0].(model.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrderByNumber indicates an expected call of GetOrderByNumber.
func (mr *MockStorageMockRecorder) GetOrderByNumber(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrderByNumber", reflect.TypeOf((*MockStorage)(nil).GetOrderByNumber), arg0, arg1)
}

// GetOrdersByUserID mocks base method.
func (m *MockStorage) GetOrdersByUserID(arg0 context.Context, arg1 uint) ([]model.Order, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetOrdersByUserID", arg0, arg1)
	ret0, _ := ret[0].([]model.Order)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetOrdersByUserID indicates an expected call of GetOrdersByUserID.
func (mr *MockStorageMockRecorder) GetOrdersByUserID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetOrdersByUserID", reflect.TypeOf((*MockStorage)(nil).GetOrdersByUserID), arg0, arg1)
}

// GetPendingOrders mocks base method.
func (m *MockStorage) GetPendingOrders(arg0 context.Context) ([]string, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPendingOrders", arg0)
	ret0, _ := ret[0].([]string)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPendingOrders indicates an expected call of GetPendingOrders.
func (mr *MockStorageMockRecorder) GetPendingOrders(arg0 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPendingOrders", reflect.TypeOf((*MockStorage)(nil).GetPendingOrders), arg0)
}

// GetUserByID mocks base method.
func (m *MockStorage) GetUserByID(arg0 context.Context, arg1 uint) (model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByID", arg0, arg1)
	ret0, _ := ret[0].(model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByID indicates an expected call of GetUserByID.
func (mr *MockStorageMockRecorder) GetUserByID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByID", reflect.TypeOf((*MockStorage)(nil).GetUserByID), arg0, arg1)
}

// GetUserByLogin mocks base method.
func (m *MockStorage) GetUserByLogin(arg0 context.Context, arg1 string) (model.User, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetUserByLogin", arg0, arg1)
	ret0, _ := ret[0].(model.User)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetUserByLogin indicates an expected call of GetUserByLogin.
func (mr *MockStorageMockRecorder) GetUserByLogin(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetUserByLogin", reflect.TypeOf((*MockStorage)(nil).GetUserByLogin), arg0, arg1)
}

// GetWithdrawalsByUserID mocks base method.
func (m *MockStorage) GetWithdrawalsByUserID(arg0 context.Context, arg1 uint) ([]model.Withdrawal, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWithdrawalsByUserID", arg0, arg1)
	ret0, _ := ret[0].([]model.Withdrawal)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWithdrawalsByUserID indicates an expected call of GetWithdrawalsByUserID.
func (mr *MockStorageMockRecorder) GetWithdrawalsByUserID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWithdrawalsByUserID", reflect.TypeOf((*MockStorage)(nil).GetWithdrawalsByUserID), arg0, arg1)
}

// GetWithdrawalsSumByUserID mocks base method.
func (m *MockStorage) GetWithdrawalsSumByUserID(arg0 context.Context, arg1 uint) (float64, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetWithdrawalsSumByUserID", arg0, arg1)
	ret0, _ := ret[0].(float64)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetWithdrawalsSumByUserID indicates an expected call of GetWithdrawalsSumByUserID.
func (mr *MockStorageMockRecorder) GetWithdrawalsSumByUserID(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetWithdrawalsSumByUserID", reflect.TypeOf((*MockStorage)(nil).GetWithdrawalsSumByUserID), arg0, arg1)
}

// ProcessWithdrawal mocks base method.
func (m *MockStorage) ProcessWithdrawal(arg0 context.Context, arg1 model.Withdrawal) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ProcessWithdrawal", arg0, arg1)
	ret0, _ := ret[0].(error)
	return ret0
}

// ProcessWithdrawal indicates an expected call of ProcessWithdrawal.
func (mr *MockStorageMockRecorder) ProcessWithdrawal(arg0, arg1 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ProcessWithdrawal", reflect.TypeOf((*MockStorage)(nil).ProcessWithdrawal), arg0, arg1)
}

// SetOrderStatus mocks base method.
func (m *MockStorage) SetOrderStatus(arg0 context.Context, arg1 string, arg2 model.OrderState) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SetOrderStatus", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// SetOrderStatus indicates an expected call of SetOrderStatus.
func (mr *MockStorageMockRecorder) SetOrderStatus(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SetOrderStatus", reflect.TypeOf((*MockStorage)(nil).SetOrderStatus), arg0, arg1, arg2)
}

// UploadOrder mocks base method.
func (m *MockStorage) UploadOrder(arg0 context.Context, arg1 uint, arg2 string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "UploadOrder", arg0, arg1, arg2)
	ret0, _ := ret[0].(error)
	return ret0
}

// UploadOrder indicates an expected call of UploadOrder.
func (mr *MockStorageMockRecorder) UploadOrder(arg0, arg1, arg2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "UploadOrder", reflect.TypeOf((*MockStorage)(nil).UploadOrder), arg0, arg1, arg2)
}
