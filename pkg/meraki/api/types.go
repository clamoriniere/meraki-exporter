package api

import (
	"fmt"
)

// Version represent the API version
type Version string

// Organization represent an organization
type Organization struct {
	Version Version

	Id   int    `json:"id,omitempty"`
	Name string `json:"name,omitempty"`
}

func (o Organization) String() string {
	return fmt.Sprintf("Organization: Id:[%d] Name:[%s]", o.Id, o.Name)
}

// OrganizationList Organization list
type OrganizationList struct {
	Version Version

	Items []Organization
}

// Network represent a network API object
type Network struct {
	Version Version

	Id       string `json:"id,omitempty"`
	Name     string `json:"name,omitempty"`
	TimeZone string `json:"timeZone,omitempty"`
	Type     string `json:"type,omitempty"`
}

func (n Network) String() string {
	return fmt.Sprintf("Network: Id:[%s] Name:[%s] TimeZone:[%s] Type:[%s]", n.Id, n.Name, n.TimeZone, n.Type)
}

// NetworkList Network list
type NetworkList struct {
	Version Version

	Items []Network
}

// Device represent a device API object
type Device struct {
	Version Version

	Wan1IP    string  `json:"wan1Ip,omitempty"`
	Wan2IP    string  `json:"wan2Ip,omitempty"`
	Serial    string  `json:"serial,omitempty"`
	Mac       string  `json:"mac,omitempty"`
	Lat       float64 `json:"lat,omitempty"`
	Lng       float64 `json:"lng,omitempty"`
	Address   string  `json:"address,omitempty"`
	Tags      string  `json:"tags,omitempty"`
	Name      string  `json:"name,omitempty"`
	Model     string  `json:"model,omitempty"`
	NetworkID string  `json:"networkId,omitempty"`
}

func (d Device) String() string {
	return fmt.Sprintf("Device: Name:[%s] Wan1IP:[%s] Wan2IP:[%s] Serial:[%s]", d.Name, d.Wan1IP, d.Wan2IP, d.Serial)
}

// DeviceList Device list
type DeviceList struct {
	Version Version

	Items []Device
}

// Performance represent a performance API object
type Performance struct {
	Version Version

	PerfScore int `json:"perfScore,omitempty"`
}

func (p Performance) String() string {
	return fmt.Sprintf("Performance: %d", p.PerfScore)
}

// PerformanceList Performance list
type PerformanceList struct {
	Version Version

	Items []Performance
}
