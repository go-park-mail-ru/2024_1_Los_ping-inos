// Code generated by MockGen. DO NOT EDIT.
// Source: interfaces.go

// Package mocks is a generated GoMock package.
package mocks

import (
	context "context"
	reflect "reflect"

	gomock "github.com/golang/mock/gomock"
	websocket "github.com/gorilla/websocket"
	grpc "google.golang.org/grpc"
	feed "main.go/internal/feed"
	proto "main.go/internal/image/protos/gen"
	types "main.go/internal/types"
)

// MockUseCase is a mock of UseCase interface.
type MockUseCase struct {
	ctrl     *gomock.Controller
	recorder *MockUseCaseMockRecorder
}

// MockUseCaseMockRecorder is the mock recorder for MockUseCase.
type MockUseCaseMockRecorder struct {
	mock *MockUseCase
}

// NewMockUseCase creates a new mock instance.
func NewMockUseCase(ctrl *gomock.Controller) *MockUseCase {
	mock := &MockUseCase{ctrl: ctrl}
	mock.recorder = &MockUseCaseMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockUseCase) EXPECT() *MockUseCaseMockRecorder {
	return m.recorder
}

// AddConnection mocks base method.
func (m *MockUseCase) AddConnection(ctx context.Context, connection *websocket.Conn, UID types.UserID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddConnection", ctx, connection, UID)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddConnection indicates an expected call of AddConnection.
func (mr *MockUseCaseMockRecorder) AddConnection(ctx, connection, UID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddConnection", reflect.TypeOf((*MockUseCase)(nil).AddConnection), ctx, connection, UID)
}

// CreateClaim mocks base method.
func (m *MockUseCase) CreateClaim(ctx context.Context, typeID, senderID, receiverID int64) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateClaim", ctx, typeID, senderID, receiverID)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateClaim indicates an expected call of CreateClaim.
func (mr *MockUseCaseMockRecorder) CreateClaim(ctx, typeID, senderID, receiverID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateClaim", reflect.TypeOf((*MockUseCase)(nil).CreateClaim), ctx, typeID, senderID, receiverID)
}

// CreateLike mocks base method.
func (m *MockUseCase) CreateLike(profile1, profile2 types.UserID, ctx context.Context) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateLike", profile1, profile2, ctx)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateLike indicates an expected call of CreateLike.
func (mr *MockUseCaseMockRecorder) CreateLike(profile1, profile2, ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateLike", reflect.TypeOf((*MockUseCase)(nil).CreateLike), profile1, profile2, ctx)
}

// DeleteConnection mocks base method.
func (m *MockUseCase) DeleteConnection(ctx context.Context, UID types.UserID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteConnection", ctx, UID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteConnection indicates an expected call of DeleteConnection.
func (mr *MockUseCaseMockRecorder) DeleteConnection(ctx, UID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteConnection", reflect.TypeOf((*MockUseCase)(nil).DeleteConnection), ctx, UID)
}

// GetCards mocks base method.
func (m *MockUseCase) GetCards(userID types.UserID, ctx context.Context) ([]feed.Card, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetCards", userID, ctx)
	ret0, _ := ret[0].([]feed.Card)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetCards indicates an expected call of GetCards.
func (mr *MockUseCaseMockRecorder) GetCards(userID, ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetCards", reflect.TypeOf((*MockUseCase)(nil).GetCards), userID, ctx)
}

// GetChat mocks base method.
func (m *MockUseCase) GetChat(ctx context.Context, user1, user2 types.UserID) ([]feed.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetChat", ctx, user1, user2)
	ret0, _ := ret[0].([]feed.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetChat indicates an expected call of GetChat.
func (mr *MockUseCaseMockRecorder) GetChat(ctx, user1, user2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetChat", reflect.TypeOf((*MockUseCase)(nil).GetChat), ctx, user1, user2)
}

// GetClaims mocks base method.
func (m *MockUseCase) GetClaims(ctx context.Context) ([]feed.PureClaim, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetClaims", ctx)
	ret0, _ := ret[0].([]feed.PureClaim)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetClaims indicates an expected call of GetClaims.
func (mr *MockUseCaseMockRecorder) GetClaims(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetClaims", reflect.TypeOf((*MockUseCase)(nil).GetClaims), ctx)
}

// GetConnection mocks base method.
func (m *MockUseCase) GetConnection(ctx context.Context, UID types.UserID) (*websocket.Conn, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetConnection", ctx, UID)
	ret0, _ := ret[0].(*websocket.Conn)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetConnection indicates an expected call of GetConnection.
func (mr *MockUseCaseMockRecorder) GetConnection(ctx, UID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConnection", reflect.TypeOf((*MockUseCase)(nil).GetConnection), ctx, UID)
}

// GetLastMessages mocks base method.
func (m *MockUseCase) GetLastMessages(ctx context.Context, UID int64, ids []int64) ([]feed.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLastMessages", ctx, UID, ids)
	ret0, _ := ret[0].([]feed.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLastMessages indicates an expected call of GetLastMessages.
func (mr *MockUseCaseMockRecorder) GetLastMessages(ctx, UID, ids interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLastMessages", reflect.TypeOf((*MockUseCase)(nil).GetLastMessages), ctx, UID, ids)
}

// SaveMessage mocks base method.
func (m *MockUseCase) SaveMessage(ctx context.Context, message feed.MessageToReceive) (*feed.MessageToReceive, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "SaveMessage", ctx, message)
	ret0, _ := ret[0].(*feed.MessageToReceive)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// SaveMessage indicates an expected call of SaveMessage.
func (mr *MockUseCaseMockRecorder) SaveMessage(ctx, message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "SaveMessage", reflect.TypeOf((*MockUseCase)(nil).SaveMessage), ctx, message)
}

// MockPostgresStorage is a mock of PostgresStorage interface.
type MockPostgresStorage struct {
	ctrl     *gomock.Controller
	recorder *MockPostgresStorageMockRecorder
}

// MockPostgresStorageMockRecorder is the mock recorder for MockPostgresStorage.
type MockPostgresStorageMockRecorder struct {
	mock *MockPostgresStorage
}

// NewMockPostgresStorage creates a new mock instance.
func NewMockPostgresStorage(ctrl *gomock.Controller) *MockPostgresStorage {
	mock := &MockPostgresStorage{ctrl: ctrl}
	mock.recorder = &MockPostgresStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockPostgresStorage) EXPECT() *MockPostgresStorageMockRecorder {
	return m.recorder
}

// CreateClaim mocks base method.
func (m *MockPostgresStorage) CreateClaim(ctx context.Context, claim feed.Claim) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateClaim", ctx, claim)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateClaim indicates an expected call of CreateClaim.
func (mr *MockPostgresStorageMockRecorder) CreateClaim(ctx, claim interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateClaim", reflect.TypeOf((*MockPostgresStorage)(nil).CreateClaim), ctx, claim)
}

// CreateLike mocks base method.
func (m *MockPostgresStorage) CreateLike(ctx context.Context, person1ID, person2ID types.UserID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateLike", ctx, person1ID, person2ID)
	ret0, _ := ret[0].(error)
	return ret0
}

// CreateLike indicates an expected call of CreateLike.
func (mr *MockPostgresStorageMockRecorder) CreateLike(ctx, person1ID, person2ID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateLike", reflect.TypeOf((*MockPostgresStorage)(nil).CreateLike), ctx, person1ID, person2ID)
}

// CreateMessage mocks base method.
func (m *MockPostgresStorage) CreateMessage(ctx context.Context, message feed.MessageToReceive) (*feed.MessageToReceive, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "CreateMessage", ctx, message)
	ret0, _ := ret[0].(*feed.MessageToReceive)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// CreateMessage indicates an expected call of CreateMessage.
func (mr *MockPostgresStorageMockRecorder) CreateMessage(ctx, message interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "CreateMessage", reflect.TypeOf((*MockPostgresStorage)(nil).CreateMessage), ctx, message)
}

// DecreaseLikesCount mocks base method.
func (m *MockPostgresStorage) DecreaseLikesCount(ctx context.Context, personID types.UserID) (int, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DecreaseLikesCount", ctx, personID)
	ret0, _ := ret[0].(int)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// DecreaseLikesCount indicates an expected call of DecreaseLikesCount.
func (mr *MockPostgresStorageMockRecorder) DecreaseLikesCount(ctx, personID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DecreaseLikesCount", reflect.TypeOf((*MockPostgresStorage)(nil).DecreaseLikesCount), ctx, personID)
}

// GetAllClaims mocks base method.
func (m *MockPostgresStorage) GetAllClaims(ctx context.Context) ([]feed.PureClaim, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetAllClaims", ctx)
	ret0, _ := ret[0].([]feed.PureClaim)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetAllClaims indicates an expected call of GetAllClaims.
func (mr *MockPostgresStorageMockRecorder) GetAllClaims(ctx interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetAllClaims", reflect.TypeOf((*MockPostgresStorage)(nil).GetAllClaims), ctx)
}

// GetChat mocks base method.
func (m *MockPostgresStorage) GetChat(ctx context.Context, user1, user2 types.UserID) ([]feed.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetChat", ctx, user1, user2)
	ret0, _ := ret[0].([]feed.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetChat indicates an expected call of GetChat.
func (mr *MockPostgresStorageMockRecorder) GetChat(ctx, user1, user2 interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetChat", reflect.TypeOf((*MockPostgresStorage)(nil).GetChat), ctx, user1, user2)
}

// GetFeed mocks base method.
func (m *MockPostgresStorage) GetFeed(ctx context.Context, filter types.UserID) ([]*feed.Person, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetFeed", ctx, filter)
	ret0, _ := ret[0].([]*feed.Person)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetFeed indicates an expected call of GetFeed.
func (mr *MockPostgresStorageMockRecorder) GetFeed(ctx, filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetFeed", reflect.TypeOf((*MockPostgresStorage)(nil).GetFeed), ctx, filter)
}

// GetLastMessages mocks base method.
func (m *MockPostgresStorage) GetLastMessages(ctx context.Context, id int64, ids []int) ([]feed.Message, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLastMessages", ctx, id, ids)
	ret0, _ := ret[0].([]feed.Message)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLastMessages indicates an expected call of GetLastMessages.
func (mr *MockPostgresStorageMockRecorder) GetLastMessages(ctx, id, ids interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLastMessages", reflect.TypeOf((*MockPostgresStorage)(nil).GetLastMessages), ctx, id, ids)
}

// GetLike mocks base method.
func (m *MockPostgresStorage) GetLike(ctx context.Context, filter *feed.LikeGetFilter) ([]types.UserID, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetLike", ctx, filter)
	ret0, _ := ret[0].([]types.UserID)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetLike indicates an expected call of GetLike.
func (mr *MockPostgresStorageMockRecorder) GetLike(ctx, filter interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetLike", reflect.TypeOf((*MockPostgresStorage)(nil).GetLike), ctx, filter)
}

// GetPersonInterests mocks base method.
func (m *MockPostgresStorage) GetPersonInterests(ctx context.Context, personID types.UserID) ([]*feed.Interest, error) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetPersonInterests", ctx, personID)
	ret0, _ := ret[0].([]*feed.Interest)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetPersonInterests indicates an expected call of GetPersonInterests.
func (mr *MockPostgresStorageMockRecorder) GetPersonInterests(ctx, personID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetPersonInterests", reflect.TypeOf((*MockPostgresStorage)(nil).GetPersonInterests), ctx, personID)
}

// MockWebSocStorage is a mock of WebSocStorage interface.
type MockWebSocStorage struct {
	ctrl     *gomock.Controller
	recorder *MockWebSocStorageMockRecorder
}

// MockWebSocStorageMockRecorder is the mock recorder for MockWebSocStorage.
type MockWebSocStorageMockRecorder struct {
	mock *MockWebSocStorage
}

// NewMockWebSocStorage creates a new mock instance.
func NewMockWebSocStorage(ctrl *gomock.Controller) *MockWebSocStorage {
	mock := &MockWebSocStorage{ctrl: ctrl}
	mock.recorder = &MockWebSocStorageMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockWebSocStorage) EXPECT() *MockWebSocStorageMockRecorder {
	return m.recorder
}

// AddConnection mocks base method.
func (m *MockWebSocStorage) AddConnection(ctx context.Context, connection *websocket.Conn, UID types.UserID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "AddConnection", ctx, connection, UID)
	ret0, _ := ret[0].(error)
	return ret0
}

// AddConnection indicates an expected call of AddConnection.
func (mr *MockWebSocStorageMockRecorder) AddConnection(ctx, connection, UID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "AddConnection", reflect.TypeOf((*MockWebSocStorage)(nil).AddConnection), ctx, connection, UID)
}

// DeleteConnection mocks base method.
func (m *MockWebSocStorage) DeleteConnection(ctx context.Context, UID types.UserID) error {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "DeleteConnection", ctx, UID)
	ret0, _ := ret[0].(error)
	return ret0
}

// DeleteConnection indicates an expected call of DeleteConnection.
func (mr *MockWebSocStorageMockRecorder) DeleteConnection(ctx, UID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "DeleteConnection", reflect.TypeOf((*MockWebSocStorage)(nil).DeleteConnection), ctx, UID)
}

// GetConnection mocks base method.
func (m *MockWebSocStorage) GetConnection(ctx context.Context, UID types.UserID) (*websocket.Conn, bool) {
	m.ctrl.T.Helper()
	ret := m.ctrl.Call(m, "GetConnection", ctx, UID)
	ret0, _ := ret[0].(*websocket.Conn)
	ret1, _ := ret[1].(bool)
	return ret0, ret1
}

// GetConnection indicates an expected call of GetConnection.
func (mr *MockWebSocStorageMockRecorder) GetConnection(ctx, UID interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetConnection", reflect.TypeOf((*MockWebSocStorage)(nil).GetConnection), ctx, UID)
}

// MockImageClient is a mock of ImageClient interface.
type MockImageClient struct {
	ctrl     *gomock.Controller
	recorder *MockImageClientMockRecorder
}

// MockImageClientMockRecorder is the mock recorder for MockImageClient.
type MockImageClientMockRecorder struct {
	mock *MockImageClient
}

// NewMockImageClient creates a new mock instance.
func NewMockImageClient(ctrl *gomock.Controller) *MockImageClient {
	mock := &MockImageClient{ctrl: ctrl}
	mock.recorder = &MockImageClientMockRecorder{mock}
	return mock
}

// EXPECT returns an object that allows the caller to indicate expected use.
func (m *MockImageClient) EXPECT() *MockImageClientMockRecorder {
	return m.recorder
}

// GetImage mocks base method.
func (m *MockImageClient) GetImage(ctx context.Context, in *proto.GetImageRequest, opts ...grpc.CallOption) (*proto.GetImageResponce, error) {
	m.ctrl.T.Helper()
	varargs := []interface{}{ctx, in}
	for _, a := range opts {
		varargs = append(varargs, a)
	}
	ret := m.ctrl.Call(m, "GetImage", varargs...)
	ret0, _ := ret[0].(*proto.GetImageResponce)
	ret1, _ := ret[1].(error)
	return ret0, ret1
}

// GetImage indicates an expected call of GetImage.
func (mr *MockImageClientMockRecorder) GetImage(ctx, in interface{}, opts ...interface{}) *gomock.Call {
	mr.mock.ctrl.T.Helper()
	varargs := append([]interface{}{ctx, in}, opts...)
	return mr.mock.ctrl.RecordCallWithMethodType(mr.mock, "GetImage", reflect.TypeOf((*MockImageClient)(nil).GetImage), varargs...)
}
