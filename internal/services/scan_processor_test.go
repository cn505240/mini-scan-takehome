package services

import (
	"context"
	"testing"
	"time"

	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"

	"github.com/censys/scan-takehome/internal/domain"
	"github.com/censys/scan-takehome/internal/mocks"
)

func TestScanProcessor_ProcessScanResult_NewScan(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockScanRepository(ctrl)
	processor := NewScanProcessor(mockRepo)

	scan := &domain.ServiceScan{
		IP:          "192.168.1.1",
		Port:        8080,
		Service:     "HTTP",
		Response:    "Hello World",
		LastScanned: time.Now(),
	}

	mockRepo.EXPECT().
		GetLatestScan(gomock.Any(), "192.168.1.1", uint32(8080), "HTTP").
		Return(nil, nil)

	mockRepo.EXPECT().
		UpsertScan(gomock.Any(), scan).
		Return(nil)

	err := processor.ProcessScanResult(context.Background(), scan)

	assert.NoError(t, err)
}

func TestScanProcessor_ProcessScanResult_NewerScan(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockScanRepository(ctrl)
	processor := NewScanProcessor(mockRepo)

	now := time.Now()
	olderTime := now.Add(-time.Hour)

	scan := &domain.ServiceScan{
		IP:          "192.168.1.1",
		Port:        8080,
		Service:     "HTTP",
		Response:    "New Response",
		LastScanned: now,
	}

	existingScan := &domain.ServiceScan{
		IP:          "192.168.1.1",
		Port:        8080,
		Service:     "HTTP",
		Response:    "Old Response",
		LastScanned: olderTime,
	}

	mockRepo.EXPECT().
		GetLatestScan(gomock.Any(), "192.168.1.1", uint32(8080), "HTTP").
		Return(existingScan, nil)

	mockRepo.EXPECT().
		UpsertScan(gomock.Any(), scan).
		Return(nil)

	err := processor.ProcessScanResult(context.Background(), scan)

	assert.NoError(t, err)
}

func TestScanProcessor_ProcessScanResult_OlderScan(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := mocks.NewMockScanRepository(ctrl)
	processor := NewScanProcessor(mockRepo)

	now := time.Now()
	newerTime := now.Add(time.Hour)

	scan := &domain.ServiceScan{
		IP:          "192.168.1.1",
		Port:        8080,
		Service:     "HTTP",
		Response:    "Old Response",
		LastScanned: now,
	}

	existingScan := &domain.ServiceScan{
		IP:          "192.168.1.1",
		Port:        8080,
		Service:     "HTTP",
		Response:    "Newer Response",
		LastScanned: newerTime,
	}

	mockRepo.EXPECT().
		GetLatestScan(gomock.Any(), "192.168.1.1", uint32(8080), "HTTP").
		Return(existingScan, nil)

	err := processor.ProcessScanResult(context.Background(), scan)

	assert.NoError(t, err)
}
