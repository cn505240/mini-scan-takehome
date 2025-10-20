package handlers

import (
	"context"
	"encoding/json"
	"testing"
	"time"

	"github.com/censys/scan-takehome/internal/domain"
	"github.com/censys/scan-takehome/internal/mocks"
	"github.com/censys/scan-takehome/pkg/scanning"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
)

func TestMessageHandler_HandleMessage_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProcessor := mocks.NewMockScanProcessor(ctrl)
	handler := NewMessageHandler(mockProcessor)

	scan := scanning.Scan{
		Ip:          "192.168.1.1",
		Port:        8080,
		Service:     "HTTP",
		Timestamp:   1640995200, // 2022-01-01 00:00:00
		DataVersion: scanning.V2,
		Data: &scanning.V2Data{
			ResponseStr: "Hello World",
		},
	}

	msgData, err := json.Marshal(scan)
	if err != nil {
		t.Fatalf("Failed to marshal scan: %v", err)
	}

	expectedScanResult := domain.ScanResult{
		IP:        "192.168.1.1",
		Port:      8080,
		Service:   "HTTP",
		Response:  "Hello World",
		Timestamp: time.Unix(scan.Timestamp, 0),
	}

	mockProcessor.EXPECT().
		ProcessScanResult(gomock.Any(), gomock.Any()).
		Do(func(ctx context.Context, scanResult domain.ScanResult) {
			if scanResult.IP != expectedScanResult.IP {
				t.Errorf("Expected IP %s, got %s", expectedScanResult.IP, scanResult.IP)
			}
			if scanResult.Port != expectedScanResult.Port {
				t.Errorf("Expected Port %d, got %d", expectedScanResult.Port, scanResult.Port)
			}
			if scanResult.Service != expectedScanResult.Service {
				t.Errorf("Expected Service %s, got %s", expectedScanResult.Service, scanResult.Service)
			}
			if scanResult.Response != expectedScanResult.Response {
				t.Errorf("Expected Response %s, got %s", expectedScanResult.Response, scanResult.Response)
			}
		}).
		Return(nil)

	err = handler.HandleMessage(context.Background(), msgData)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestMessageHandler_HandleMessage_V1Data_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProcessor := mocks.NewMockScanProcessor(ctrl)
	handler := NewMessageHandler(mockProcessor)

	msgData := []byte(`{
		"ip": "10.0.0.1",
		"port": 22,
		"service": "SSH",
		"timestamp": 1640995200,
		"data_version": 1,
		"data": {
			"response_bytes_utf8": "U1NILTIuMC1PcGVuU1NIXzguMg=="
		}
	}`)

	mockProcessor.EXPECT().
		ProcessScanResult(gomock.Any(), gomock.Any()).
		Do(func(ctx context.Context, scanResult domain.ScanResult) {
			if scanResult.IP != "10.0.0.1" {
				t.Errorf("Expected IP 10.0.0.1, got %s", scanResult.IP)
			}
			if scanResult.Port != 22 {
				t.Errorf("Expected Port 22, got %d", scanResult.Port)
			}
			if scanResult.Service != "SSH" {
				t.Errorf("Expected Service SSH, got %s", scanResult.Service)
			}
			if scanResult.Response != "SSH-2.0-OpenSSH_8.2" {
				t.Errorf("Expected Response SSH-2.0-OpenSSH_8.2, got %s", scanResult.Response)
			}
		}).
		Return(nil)

	err := handler.HandleMessage(context.Background(), msgData)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}

func TestMessageHandler_HandleMessage_InvalidJSON(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProcessor := mocks.NewMockScanProcessor(ctrl)
	handler := NewMessageHandler(mockProcessor)

	msgData := []byte(`{"invalid": json}`)

	err := handler.HandleMessage(context.Background(), msgData)
	if err == nil {
		t.Error("Expected error for invalid JSON, got nil")
	}
}

func TestMessageHandler_HandleMessage_ProcessorError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProcessor := mocks.NewMockScanProcessor(ctrl)
	handler := NewMessageHandler(mockProcessor)

	scan := scanning.Scan{
		Ip:          "192.168.1.1",
		Port:        8080,
		Service:     "HTTP",
		Timestamp:   1640995200,
		DataVersion: scanning.V2,
		Data: &scanning.V2Data{
			ResponseStr: "Hello World",
		},
	}

	msgData, err := json.Marshal(scan)
	if err != nil {
		t.Fatalf("Failed to marshal scan: %v", err)
	}

	mockProcessor.EXPECT().
		ProcessScanResult(gomock.Any(), gomock.Any()).
		Return(assert.AnError)

	err = handler.HandleMessage(context.Background(), msgData)
	if err == nil {
		t.Error("Expected error from processor, got nil")
	}
}

func TestMessageHandler_HandleMessage_UnknownDataVersion(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockProcessor := mocks.NewMockScanProcessor(ctrl)
	handler := NewMessageHandler(mockProcessor)

	scan := scanning.Scan{
		Ip:          "192.168.1.1",
		Port:        8080,
		Service:     "HTTP",
		Timestamp:   1640995200,
		DataVersion: 999, // Unknown version
		Data:        "some data",
	}

	msgData, err := json.Marshal(scan)
	if err != nil {
		t.Fatalf("Failed to marshal scan: %v", err)
	}

	mockProcessor.EXPECT().
		ProcessScanResult(gomock.Any(), gomock.Any()).
		Do(func(ctx context.Context, scanResult domain.ScanResult) {
			if scanResult.Response != "" {
				t.Errorf("Expected empty response for unknown data version, got %s", scanResult.Response)
			}
		}).
		Return(nil)

	err = handler.HandleMessage(context.Background(), msgData)
	if err != nil {
		t.Errorf("Expected no error, got %v", err)
	}
}
