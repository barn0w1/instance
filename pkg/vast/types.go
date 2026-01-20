package vast

// Offer represents a GPU instance offer from Vast.ai
type Offer struct {
	ID         int     `json:"id"`
	MachineID  int     `json:"machine_id"`
	HostID     int     `json:"host_id"`
	GpuName    string  `json:"gpu_name"`
	NumGpus    int     `json:"num_gpus"`
	GpuRam     int     `json:"gpu_ram"` // MB
	TotalFlops float64 `json:"total_flops"`

	// Price & Costs
	DphTotal    float64 `json:"dph_total"` // $/hr
	DphBase     float64 `json:"dph_base"`
	StorageCost float64 `json:"storage_cost"` // $/GB/month

	// Performance & Specs
	DlPerf      float64 `json:"dlperf"`
	Reliability float64 `json:"reliability2"`
	InetDown    float64 `json:"inet_down"` // MB/s
	InetUp      float64 `json:"inet_up"`   // MB/s
	CpuCores    int     `json:"cpu_cores"`
	CpuRam      int     `json:"cpu_ram"`    // MB
	DiskSpace   float64 `json:"disk_space"` // GB

	// Status
	Verified bool   `json:"verification,string"` // "verified" comes as string often, strictly checking API response type is safer if generic
	Rented   bool   `json:"rented"`
	Region   string `json:"geolocation"`
}

// searchResponse wraps the list of offers
type searchResponse struct {
	Offers []Offer `json:"offers"`
}
