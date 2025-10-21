package domain

import "time"

// ServiceScan represents a service scan record
type ServiceScan struct {
	IP          string    `json:"ip"`
	Port        uint32    `json:"port"`
	Service     string    `json:"service"`
	Response    string    `json:"response"`
	LastScanned time.Time `json:"last_scanned"`
}

// IsNewerThan checks if this scan is newer than the given timestamp
func (ss *ServiceScan) IsNewerThan(timestamp time.Time) bool {
	return ss.LastScanned.After(timestamp)
}
