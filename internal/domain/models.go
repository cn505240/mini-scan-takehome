package domain

import "time"

// ServiceScan represents a unique service scan record
type ServiceScan struct {
	IP          string    `json:"ip"`
	Port        uint32    `json:"port"`
	Service     string    `json:"service"`
	Response    string    `json:"response"`
	LastScanned time.Time `json:"last_scanned"`
}

// ScanResult represents an incoming scan result from the message queue
type ScanResult struct {
	IP        string    `json:"ip"`
	Port      uint32    `json:"port"`
	Service   string    `json:"service"`
	Response  string    `json:"response"`
	Timestamp time.Time `json:"timestamp"`
}

// IsNewerThan checks if this scan result is newer than the given timestamp
func (sr *ScanResult) IsNewerThan(timestamp time.Time) bool {
	return sr.Timestamp.After(timestamp)
}

// ToServiceScan converts a ScanResult to a ServiceScan
func (sr *ScanResult) ToServiceScan() ServiceScan {
	return ServiceScan{
		IP:          sr.IP,
		Port:        sr.Port,
		Service:     sr.Service,
		Response:    sr.Response,
		LastScanned: sr.Timestamp,
	}
}
