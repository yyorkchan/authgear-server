// Code generated by MockGen. DO NOT EDIT.
// Source: service.go

// Package verification is a generated GoMock package.
package verification

import (
	config "github.com/authgear/authgear-server/pkg/auth/config"
	authenticator "github.com/authgear/authgear-server/pkg/auth/dependency/authenticator"
	identity "github.com/authgear/authgear-server/pkg/auth/dependency/identity"
	otp "github.com/authgear/authgear-server/pkg/otp"
	gomock "github.com/golang/mock/gomock"
	reflect "reflect"
)

// MockIdentityProvider is a mock of IdentityProvider interface
type MockIdentityProvider struct {
	ctrl     *gomock.Controller
	recorder *MockIdentityProviderMockRecorder
}

// MockIdentityProviderMockRecorder is the mock recorder for MockIdentityProvider
type MockIdentityProviderMockRecorder struct {
	mock *MockIdentityProvider
}

// NewMockIdentityProvider creates a new mock instance
func NewMockIdentityProvider(ctrl *gomock.Controller) *MockIdentityProvider {
	mock := &MockIdentityProvider{ctrl: ctrl}
	mock.recorder = &MockIdentityProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockIdentityProvider) EXPECT() *MockIdentityProviderMockRecorder {
	return m.recorder
}

// ListByUser mocks base method
func (m *MockIdentityProvider) ListByUser(userID string) ([]*identity.Info, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "ListByUser", userID)
	ret0, _ := ret[0].([]*identity.Info)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// ListByUser indicates an expected call of ListByUser
func (mr *MockIdentityProviderMockRecorder) ListByUser(userID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "ListByUser", reflect.TypeOf((*MockIdentityProvider)(nil).ListByUser), userID)
}

// RelateIdentityToAuthenticator mocks base method
func (m *MockIdentityProvider) RelateIdentityToAuthenticator(ii *identity.Info, as *authenticator.Spec) *authenticator.Spec {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "RelateIdentityToAuthenticator", ii, as)
	ret0, _ := ret[0].(*authenticator.Spec)
	return ret0
}

// RelateIdentityToAuthenticator indicates an expected call of RelateIdentityToAuthenticator
func (mr *MockIdentityProviderMockRecorder) RelateIdentityToAuthenticator(ii, as interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "RelateIdentityToAuthenticator", reflect.TypeOf((*MockIdentityProvider)(nil).RelateIdentityToAuthenticator), ii, as)
}

// MockAuthenticatorProvider is a mock of AuthenticatorProvider interface
type MockAuthenticatorProvider struct {
	ctrl     *gomock.Controller
	recorder *MockAuthenticatorProviderMockRecorder
}

// MockAuthenticatorProviderMockRecorder is the mock recorder for MockAuthenticatorProvider
type MockAuthenticatorProviderMockRecorder struct {
	mock *MockAuthenticatorProvider
}

// NewMockAuthenticatorProvider creates a new mock instance
func NewMockAuthenticatorProvider(ctrl *gomock.Controller) *MockAuthenticatorProvider {
	mock := &MockAuthenticatorProvider{ctrl: ctrl}
	mock.recorder = &MockAuthenticatorProviderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockAuthenticatorProvider) EXPECT() *MockAuthenticatorProviderMockRecorder {
	return m.recorder
}

// List mocks base method
func (m *MockAuthenticatorProvider) List(userID string, filters ...authenticator.Filter) ([]*authenticator.Info, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{userID}
	for _, a := range filters {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "List", varargs...)
	ret0, _ := ret[0].([]*authenticator.Info)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// List indicates an expected call of List
func (mr *MockAuthenticatorProviderMockRecorder) List(userID interface{}, filters ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{userID}, filters...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "List", reflect.TypeOf((*MockAuthenticatorProvider)(nil).List), varargs...)
}

// MockOTPMessageSender is a mock of OTPMessageSender interface
type MockOTPMessageSender struct {
	ctrl     *gomock.Controller
	recorder *MockOTPMessageSenderMockRecorder
}

// MockOTPMessageSenderMockRecorder is the mock recorder for MockOTPMessageSender
type MockOTPMessageSenderMockRecorder struct {
	mock *MockOTPMessageSender
}

// NewMockOTPMessageSender creates a new mock instance
func NewMockOTPMessageSender(ctrl *gomock.Controller) *MockOTPMessageSender {
	mock := &MockOTPMessageSender{ctrl: ctrl}
	mock.recorder = &MockOTPMessageSenderMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockOTPMessageSender) EXPECT() *MockOTPMessageSenderMockRecorder {
	return m.recorder
}

// SendEmail mocks base method
func (m *MockOTPMessageSender) SendEmail(email string, opts otp.SendOptions, message config.EmailMessageConfig) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendEmail", email, opts, message)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendEmail indicates an expected call of SendEmail
func (mr *MockOTPMessageSenderMockRecorder) SendEmail(email, opts, message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendEmail", reflect.TypeOf((*MockOTPMessageSender)(nil).SendEmail), email, opts, message)
}

// SendSMS mocks base method
func (m *MockOTPMessageSender) SendSMS(phone string, opts otp.SendOptions, message config.SMSMessageConfig) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SendSMS", phone, opts, message)
	ret0, _ := ret[0].(error)
	return ret0
}

// SendSMS indicates an expected call of SendSMS
func (mr *MockOTPMessageSenderMockRecorder) SendSMS(phone, opts, message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SendSMS", reflect.TypeOf((*MockOTPMessageSender)(nil).SendSMS), phone, opts, message)
}

// MockStore is a mock of Store interface
type MockStore struct {
	ctrl     *gomock.Controller
	recorder *MockStoreMockRecorder
}

// MockStoreMockRecorder is the mock recorder for MockStore
type MockStoreMockRecorder struct {
	mock *MockStore
}

// NewMockStore creates a new mock instance
func NewMockStore(ctrl *gomock.Controller) *MockStore {
	mock := &MockStore{ctrl: ctrl}
	mock.recorder = &MockStoreMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use
func (m *MockStore) EXPECT() *MockStoreMockRecorder {
	return m.recorder
}

// Create mocks base method
func (m *MockStore) Create(code *Code) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Create", code)
	ret0, _ := ret[0].(error)
	return ret0
}

// Create indicates an expected call of Create
func (mr *MockStoreMockRecorder) Create(code interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Create", reflect.TypeOf((*MockStore)(nil).Create), code)
}

// Get mocks base method
func (m *MockStore) Get(code string) (*Code, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Get", code)
	ret0, _ := ret[0].(*Code)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// Get indicates an expected call of Get
func (mr *MockStoreMockRecorder) Get(code interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Get", reflect.TypeOf((*MockStore)(nil).Get), code)
}

// Delete mocks base method
func (m *MockStore) Delete(code string) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "Delete", code)
	ret0, _ := ret[0].(error)
	return ret0
}

// Delete indicates an expected call of Delete
func (mr *MockStoreMockRecorder) Delete(code interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "Delete", reflect.TypeOf((*MockStore)(nil).Delete), code)
}
