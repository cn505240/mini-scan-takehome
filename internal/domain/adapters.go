package domain

import (
	"encoding/base64"
	"time"

	"github.com/censys/scan-takehome/pkg/scanning"
)

func ConvertScanToDomain(rawScan scanning.Scan) (ServiceScan, error) {
	response, err := extractServiceResponse(rawScan)
	if err != nil {
		return ServiceScan{}, err
	}

	return ServiceScan{
		IP:          rawScan.Ip,
		Port:        rawScan.Port,
		Service:     rawScan.Service,
		Response:    response,
		LastScanned: time.Unix(rawScan.Timestamp, 0),
	}, nil
}

func extractServiceResponse(rawScan scanning.Scan) (string, error) {
	switch rawScan.DataVersion {
	case scanning.V1:
		switch data := rawScan.Data.(type) {
		case *scanning.V1Data:
			decoded, err := base64.StdEncoding.DecodeString(string(data.ResponseBytesUtf8))
			if err != nil {
				return "", err
			}
			return string(decoded), nil
		case map[string]interface{}:
			if responseBytes, ok := data["response_bytes_utf8"].(string); ok {
				decoded, err := base64.StdEncoding.DecodeString(responseBytes)
				if err != nil {
					return "", err
				}
				return string(decoded), nil
			}
		}
	case scanning.V2:
		switch data := rawScan.Data.(type) {
		case *scanning.V2Data:
			return data.ResponseStr, nil
		case map[string]interface{}:
			if responseStr, ok := data["response_str"].(string); ok {
				return responseStr, nil
			}
		}
	}
	return "", nil
}
