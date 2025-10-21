package domain

import (
	"encoding/json"
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
		// Try to unmarshal into V1Data struct
		var v1Data scanning.V1Data
		if err := unmarshalData(rawScan.Data, &v1Data); err != nil {
			return "", err
		}
		// V1Data.ResponseBytesUtf8 is already decoded from base64 during JSON unmarshaling
		return string(v1Data.ResponseBytesUtf8), nil
	case scanning.V2:
		// Try to unmarshal into V2Data struct
		var v2Data scanning.V2Data
		if err := unmarshalData(rawScan.Data, &v2Data); err != nil {
			return "", err
		}
		return v2Data.ResponseStr, nil
	}
	return "", nil
}

func unmarshalData(data, target interface{}) error {
	if data == nil {
		return nil
	}

	if dataMap, ok := data.(map[string]interface{}); ok {
		jsonBytes, err := json.Marshal(dataMap)
		if err != nil {
			return err
		}
		return json.Unmarshal(jsonBytes, target)
	}

	return nil
}
