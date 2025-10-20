package domain

import (
	"encoding/base64"
	"time"

	"github.com/censys/scan-takehome/pkg/scanning"
)

func ConvertScanToDomain(rawScan scanning.Scan) (ScanResult, error) {
	response, err := extractServiceResponse(rawScan)
	if err != nil {
		return ScanResult{}, err
	}

	return ScanResult{
		IP:        rawScan.Ip,
		Port:      rawScan.Port,
		Service:   rawScan.Service,
		Response:  response,
		Timestamp: time.Unix(rawScan.Timestamp, 0),
	}, nil
}

func extractServiceResponse(rawScan scanning.Scan) (string, error) {
	switch rawScan.DataVersion {
	case scanning.V1:
		if v1Data, ok := rawScan.Data.(*scanning.V1Data); ok {
			decoded, err := base64.StdEncoding.DecodeString(string(v1Data.ResponseBytesUtf8))
			if err != nil {
				return "", err
			}
			return string(decoded), nil
		} else if dataMap, ok := rawScan.Data.(map[string]interface{}); ok {
			if responseBytes, ok := dataMap["response_bytes_utf8"].(string); ok {
				decoded, err := base64.StdEncoding.DecodeString(responseBytes)
				if err != nil {
					return "", err
				}
				return string(decoded), nil
			}
		}
	case scanning.V2:
		if v2Data, ok := rawScan.Data.(*scanning.V2Data); ok {
			return v2Data.ResponseStr, nil
		} else if dataMap, ok := rawScan.Data.(map[string]interface{}); ok {
			if responseStr, ok := dataMap["response_str"].(string); ok {
				return responseStr, nil
			}
		}
	}
	return "", nil
}
