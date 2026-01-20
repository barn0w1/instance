package vast

import (
	"encoding/json"
)

// SearchOrder represents sorting direction
type SearchOrder string

const (
	Asc  SearchOrder = "asc"
	Desc SearchOrder = "desc"
)

// InstanceType represents the pricing model
type InstanceType string

const (
	OnDemand InstanceType = "on-demand"
	Reserved InstanceType = "reserved"
	Bid      InstanceType = "bid"
)

// SearchBuilder helps construct the complex JSON query for Vast.ai
type SearchBuilder struct {
	filters map[string]interface{}
	params  map[string]interface{}
}

func NewSearch() *SearchBuilder {
	return &SearchBuilder{
		filters: make(map[string]interface{}),
		params: map[string]interface{}{
			"verified": map[string]bool{"eq": true}, // default
			"rentable": map[string]bool{"eq": true}, // default
			"type":     OnDemand,                    // default
			"limit":    100,
		},
	}
}

// --- Basic Parameters ---

func (b *SearchBuilder) Limit(n int) *SearchBuilder {
	b.params["limit"] = n
	return b
}

func (b *SearchBuilder) Type(t InstanceType) *SearchBuilder {
	b.params["type"] = t
	return b
}

func (b *SearchBuilder) Order(field string, order SearchOrder) *SearchBuilder {
	// API expects: "order": [["dph_total", "asc"]]
	b.params["order"] = [][]string{{field, string(order)}}
	return b
}

// --- Filters (Operators) ---

// addFilter is a helper to add operator-based filters
func (b *SearchBuilder) addFilter(field string, op string, value interface{}) *SearchBuilder {
	if _, ok := b.filters[field]; !ok {
		b.filters[field] = make(map[string]interface{})
	}
	// merge into existing map for this field
	f := b.filters[field].(map[string]interface{})
	f[op] = value
	return b
}

// GpuName filters by GPU model (e.g. "RTX 4090")
func (b *SearchBuilder) GpuName(name string) *SearchBuilder {
	return b.addFilter("gpu_name", "eq", name)
}

// MinGpus sets minimum number of GPUs
func (b *SearchBuilder) MinGpus(n int) *SearchBuilder {
	return b.addFilter("num_gpus", "gte", n)
}

// MaxPrice sets maximum price per hour ($)
func (b *SearchBuilder) MaxPrice(p float64) *SearchBuilder {
	return b.addFilter("dph_total", "lte", p)
}

// MinVRAM sets minimum GPU RAM in MB
func (b *SearchBuilder) MinVRAM(mb int) *SearchBuilder {
	return b.addFilter("gpu_ram", "gte", mb)
}

// MinReliability sets minimum reliability score (0.0 - 1.0)
func (b *SearchBuilder) MinReliability(score float64) *SearchBuilder {
	return b.addFilter("reliability2", "gte", score)
}

// Custom allows adding any arbitrary filter supported by the API
func (b *SearchBuilder) Custom(field, op string, value interface{}) *SearchBuilder {
	return b.addFilter(field, op, value)
}

// Build creates the final JSON body
func (b *SearchBuilder) Build() ([]byte, error) {
	// Merge params and filters into one map for JSON marshaling
	finalMap := make(map[string]interface{})

	for k, v := range b.params {
		finalMap[k] = v
	}
	for k, v := range b.filters {
		finalMap[k] = v
	}

	return json.Marshal(finalMap)
}
