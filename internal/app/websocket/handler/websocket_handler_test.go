package handler

import (
	"errors"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/Beretta350/gochat/mocks"
	"github.com/stretchr/testify/mock"
)

func TestHandleConnection_SuccessfulConnection(t *testing.T) {
	// Arrange
	mockService := mocks.NewWebsocketService(t)
	mockUpgrader := mocks.NewUpgrader(t)
	mockConn := mocks.NewConn(t)
	mockKafkaClient := mocks.NewKafkaClient(t)

	// Create a request with a token
	req := httptest.NewRequest(http.MethodGet, "/?token=test-token", nil)
	w := httptest.NewRecorder()

	// Mock the upgrade call
	mockUpgrader.On("Upgrade", w, req, mock.Anything).Return(mockConn, nil)

	// Mock the service calls
	mockService.On("SetupSession", mock.Anything, mockConn, "test-token").Return(mockKafkaClient, nil)
	mockService.On("HandleSession", mock.Anything, mockConn, mockKafkaClient).Return()

	// Mock the client close connection
	mockKafkaClient.On("CloseConnection").Return()

	// Mock the connection close
	mockConn.On("Close").Return(nil)

	// Create the handler
	handler := NewWebsocketHandler(mockService, mockUpgrader)

	// Act
	handler.HandleConnection(w, req)

	// Assert
	mockUpgrader.AssertExpectations(t)
	mockService.AssertExpectations(t)
	mockKafkaClient.AssertExpectations(t)
	mockConn.AssertExpectations(t)
}

func TestHandleConnection_NoToken(t *testing.T) {
	// Arrange
	mockService := mocks.NewWebsocketService(t)
	mockUpgrader := mocks.NewUpgrader(t)
	mockConn := mocks.NewConn(t)

	// Create a request without a token
	req := httptest.NewRequest(http.MethodGet, "/", nil)
	w := httptest.NewRecorder()

	// Mock the upgrade call
	mockUpgrader.On("Upgrade", w, req, mock.Anything).Return(mockConn, nil)

	// Mock the connection close
	mockConn.On("Close").Return(nil)

	// Create the handler
	handler := NewWebsocketHandler(mockService, mockUpgrader)

	// Act
	handler.HandleConnection(w, req)

	// Assert
	mockUpgrader.AssertExpectations(t)
	mockConn.AssertExpectations(t)
	// Service methods should not be called
	mockService.AssertNotCalled(t, "SetupSession")
	mockService.AssertNotCalled(t, "HandleSession")
}

func TestHandleConnection_UpgradeError(t *testing.T) {
	// Arrange
	mockService := mocks.NewWebsocketService(t)
	mockUpgrader := mocks.NewUpgrader(t)

	// Create a request with a token
	req := httptest.NewRequest(http.MethodGet, "/?token=test-token", nil)
	w := httptest.NewRecorder()

	// Mock the upgrade call to fail
	upgradeErr := errors.New("upgrade error")
	mockUpgrader.On("Upgrade", w, req, mock.Anything).Return(nil, upgradeErr)

	// Create the handler
	handler := NewWebsocketHandler(mockService, mockUpgrader)

	// Act
	handler.HandleConnection(w, req)

	// Assert
	mockUpgrader.AssertExpectations(t)
	// Service methods should not be called
	mockService.AssertNotCalled(t, "SetupSession")
	mockService.AssertNotCalled(t, "HandleSession")
}

func TestHandleConnection_SetupSessionError(t *testing.T) {
	// Arrange
	mockService := mocks.NewWebsocketService(t)
	mockUpgrader := mocks.NewUpgrader(t)
	mockConn := mocks.NewConn(t)

	// Create a request with a token
	req := httptest.NewRequest(http.MethodGet, "/?token=test-token", nil)
	w := httptest.NewRecorder()

	// Mock the upgrade call
	mockUpgrader.On("Upgrade", w, req, mock.Anything).Return(mockConn, nil)

	// Mock the service call to fail
	setupErr := errors.New("setup error")
	mockService.On("SetupSession", mock.Anything, mockConn, "test-token").Return(nil, setupErr)

	// Mock the connection close
	mockConn.On("Close").Return(nil)

	// Create the handler
	handler := NewWebsocketHandler(mockService, mockUpgrader)

	// Act
	handler.HandleConnection(w, req)

	// Assert
	mockUpgrader.AssertExpectations(t)
	mockService.AssertExpectations(t)
	mockConn.AssertExpectations(t)
	// HandleSession should not be called
	mockService.AssertNotCalled(t, "HandleSession")
}
