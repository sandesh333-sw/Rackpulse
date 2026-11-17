package main

type MetricStatus string

const (
	MetricStatusOk            MetricStatus = "ok"
	MetricStatusUnavailable   MetricStatus = "unavailable"
	MetricStatusNotApplicable MetricStatus = "not_applicable"
)

type MetricReading struct {
	Value  *float64     `json:"value"`
	Unit   string       `json:"unit"`
	Status MetricStatus `json:"status"`
}

type Reading struct {
	Metrics map[string]MetricReading `json:"metrics"`
}
