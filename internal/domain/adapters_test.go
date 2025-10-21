package domain

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"

	"github.com/censys/scan-takehome/pkg/scanning"
)

func TestConvertScanToDomain_V2Data(t *testing.T) {
	rawScan := scanning.Scan{
		Ip:          "192.168.1.1",
		Port:        8080,
		Service:     "HTTP",
		Timestamp:   1640995200,
		DataVersion: scanning.V2,
		Data: &scanning.V2Data{
			ResponseStr: "Hello World",
		},
	}

	result, err := ConvertScanToDomain(rawScan)

	assert.NoError(t, err)
	assert.Equal(t, "192.168.1.1", result.IP)
	assert.Equal(t, uint32(8080), result.Port)
	assert.Equal(t, "HTTP", result.Service)
	assert.Equal(t, "Hello World", result.Response)
	assert.Equal(t, time.Unix(1640995200, 0), result.LastScanned)
}

func TestConvertScanToDomain_V1Data(t *testing.T) {
	rawScan := scanning.Scan{
		Ip:          "10.0.0.1",
		Port:        22,
		Service:     "SSH",
		Timestamp:   1640995200,
		DataVersion: scanning.V1,
		Data: &scanning.V1Data{
			ResponseBytesUtf8: []byte("U1NILTIuMC1PcGVuU1NIXzguMg=="), // base64 encoded "SSH-2.0-OpenSSH_8.2"
		},
	}

	result, err := ConvertScanToDomain(rawScan)

	assert.NoError(t, err)
	assert.Equal(t, "10.0.0.1", result.IP)
	assert.Equal(t, uint32(22), result.Port)
	assert.Equal(t, "SSH", result.Service)
	assert.Equal(t, "SSH-2.0-OpenSSH_8.2", result.Response)
	assert.Equal(t, time.Unix(1640995200, 0), result.LastScanned)
}

func TestConvertScanToDomain_V1DataFromJSON(t *testing.T) {
	// Simulate JSON unmarshaling where Data becomes map[string]interface{}
	rawScan := scanning.Scan{
		Ip:          "10.0.0.1",
		Port:        22,
		Service:     "SSH",
		Timestamp:   1640995200,
		DataVersion: scanning.V1,
		Data: map[string]interface{}{
			"response_bytes_utf8": "U1NILTIuMC1PcGVuU1NIXzguMg==",
		},
	}

	result, err := ConvertScanToDomain(rawScan)

	assert.NoError(t, err)
	assert.Equal(t, "10.0.0.1", result.IP)
	assert.Equal(t, uint32(22), result.Port)
	assert.Equal(t, "SSH", result.Service)
	assert.Equal(t, "SSH-2.0-OpenSSH_8.2", result.Response)
	assert.Equal(t, time.Unix(1640995200, 0), result.LastScanned)
}

func TestConvertScanToDomain_V2DataFromJSON(t *testing.T) {
	// Simulate JSON unmarshaling where Data becomes map[string]interface{}
	rawScan := scanning.Scan{
		Ip:          "192.168.1.1",
		Port:        8080,
		Service:     "HTTP",
		Timestamp:   1640995200,
		DataVersion: scanning.V2,
		Data: map[string]interface{}{
			"response_str": "Hello World",
		},
	}

	result, err := ConvertScanToDomain(rawScan)

	assert.NoError(t, err)
	assert.Equal(t, "192.168.1.1", result.IP)
	assert.Equal(t, uint32(8080), result.Port)
	assert.Equal(t, "HTTP", result.Service)
	assert.Equal(t, "Hello World", result.Response)
	assert.Equal(t, time.Unix(1640995200, 0), result.LastScanned)
}

func TestConvertScanToDomain_UnknownDataVersion(t *testing.T) {
	rawScan := scanning.Scan{
		Ip:          "192.168.1.1",
		Port:        8080,
		Service:     "HTTP",
		Timestamp:   1640995200,
		DataVersion: 999, // Unknown version
		Data:        "some data",
	}

	result, err := ConvertScanToDomain(rawScan)

	assert.NoError(t, err)
	assert.Equal(t, "192.168.1.1", result.IP)
	assert.Equal(t, uint32(8080), result.Port)
	assert.Equal(t, "HTTP", result.Service)
	assert.Equal(t, "", result.Response) // Empty response for unknown version
	assert.Equal(t, time.Unix(1640995200, 0), result.LastScanned)
}

func TestConvertScanToDomain_InvalidBase64(t *testing.T) {
	rawScan := scanning.Scan{
		Ip:          "10.0.0.1",
		Port:        22,
		Service:     "SSH",
		Timestamp:   1640995200,
		DataVersion: scanning.V1,
		Data: &scanning.V1Data{
			ResponseBytesUtf8: []byte("invalid-base64!"),
		},
	}

	result, err := ConvertScanToDomain(rawScan)

	assert.Error(t, err)
	assert.Equal(t, ServiceScan{}, result)
}

func TestConvertScanToDomain_V1DataWithMapButInvalidBase64(t *testing.T) {
	rawScan := scanning.Scan{
		Ip:          "10.0.0.1",
		Port:        22,
		Service:     "SSH",
		Timestamp:   1640995200,
		DataVersion: scanning.V1,
		Data: map[string]interface{}{
			"response_bytes_utf8": "invalid-base64!",
		},
	}

	result, err := ConvertScanToDomain(rawScan)

	assert.Error(t, err)
	assert.Equal(t, ServiceScan{}, result)
}

func TestConvertScanToDomain_EmptyData(t *testing.T) {
	rawScan := scanning.Scan{
		Ip:          "192.168.1.1",
		Port:        8080,
		Service:     "HTTP",
		Timestamp:   1640995200,
		DataVersion: scanning.V2,
		Data:        nil,
	}

	result, err := ConvertScanToDomain(rawScan)

	assert.NoError(t, err)
	assert.Equal(t, "192.168.1.1", result.IP)
	assert.Equal(t, uint32(8080), result.Port)
	assert.Equal(t, "HTTP", result.Service)
	assert.Equal(t, "", result.Response) // Empty response for nil data
	assert.Equal(t, time.Unix(1640995200, 0), result.LastScanned)
}
