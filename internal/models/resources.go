package models

// HostResources representa os recursos dispon?veis no host
type HostResources struct {
	TotalCPU        float64 `json:"total_cpu"`          // Total de CPUs dispon?veis
	TotalMemory     int64   `json:"total_memory"`        // Total de mem?ria em bytes
	TotalDisk       int64   `json:"total_disk"`          // Total de disco em bytes
	UsedCPU         float64 `json:"used_cpu"`            // CPU em uso
	UsedMemory      int64   `json:"used_memory"`         // Mem?ria em uso
	UsedDisk        int64   `json:"used_disk"`           // Disco em uso
	AvailableCPU    float64 `json:"available_cpu"`       // CPU dispon?vel
	AvailableMemory int64   `json:"available_memory"`    // Mem?ria dispon?vel
	AvailableDisk   int64   `json:"available_disk"`     // Disco dispon?vel
}

// ServiceMinRequirements define requisitos m?nimos por tipo de servi?o
var ServiceMinRequirements = map[string]Resources{
	"n8n": {
		CPU: CPUResources{
			Cores:    0.5,
			MaxCores: 1.0,
			Priority: 512,
		},
		Memory: MemoryResources{
			Limit:       "512m",
			Reservation: "256m",
			Swap:        "256m",
		},
		Disk: DiskResources{
			Size:       "5g",
			VolumeType: "standard",
		},
	},
	"wordpress": {
		CPU: CPUResources{
			Cores:    0.5,
			MaxCores: 1.0,
			Priority: 512,
		},
		Memory: MemoryResources{
			Limit:       "512m",
			Reservation: "256m",
			Swap:        "256m",
		},
		Disk: DiskResources{
			Size:       "10g",
			VolumeType: "standard",
		},
	},
	"custom": {
		CPU: CPUResources{
			Cores:    0.25,
			MaxCores: 0.5,
			Priority: 256,
		},
		Memory: MemoryResources{
			Limit:       "256m",
			Reservation: "128m",
			Swap:        "128m",
		},
		Disk: DiskResources{
			Size:       "1g",
			VolumeType: "standard",
		},
	},
}
