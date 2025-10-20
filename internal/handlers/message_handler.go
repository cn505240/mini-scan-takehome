package handlers

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"log"
	"time"

	"github.com/censys/scan-takehome/internal/domain"
	"github.com/censys/scan-takehome/pkg/scanning"
)

type ScanProcessor interface {
	ProcessScanResult(ctx context.Context, scanResult domain.ScanResult) error
}

type MessageHandler struct {
	processor ScanProcessor
}

func NewMessageHandler(processor ScanProcessor) *MessageHandler {
	return &MessageHandler{
		processor: processor,
	}
}

func (mh *MessageHandler) HandleMessage(ctx context.Context, msgData []byte) error {
	var rawScan scanning.Scan
	if err := json.Unmarshal(msgData, &rawScan); err != nil {
		log.Printf("Failed to parse message: %v", err)
		return err
	}

	scanResult, err := mh.convertToDomainModel(rawScan)
	if err != nil {
		log.Printf("Failed to convert message to domain model: %v", err)
		return err
	}

	return mh.processor.ProcessScanResult(ctx, scanResult)
}

func (mh *MessageHandler) convertToDomainModel(rawScan scanning.Scan) (domain.ScanResult, error) {
	response, err := mh.extractServiceResponse(rawScan)
	if err != nil {
		return domain.ScanResult{}, err
	}

	return domain.ScanResult{
		IP:        rawScan.Ip,
		Port:      rawScan.Port,
		Service:   rawScan.Service,
		Response:  response,
		Timestamp: time.Unix(rawScan.Timestamp, 0),
	}, nil
}

func (mh *MessageHandler) extractServiceResponse(rawScan scanning.Scan) (string, error) {
	switch rawScan.DataVersion {
	case scanning.V1:
		if v1Data, ok := rawScan.Data.(*scanning.V1Data); ok {
			decoded, err := base64.StdEncoding.DecodeString(string(v1Data.ResponseBytesUtf8))
			if err != nil {
				return "", err
			}
			return string(decoded), nil
		}
	case scanning.V2:
		if v2Data, ok := rawScan.Data.(*scanning.V2Data); ok {
			return v2Data.ResponseStr, nil
		}
	}
	return "", nil
}
