package common

import "time"

type Entry struct {
	Value   string
	Country string
}

type LogMessage struct {
	Host            string         `json:"host"`
	RequestMethod   string         `json:"request_method"`
	RequestBody     string         `json:"request_body"`
	RequestHost     string         `json:"request_host"`
	Endpoint        string         `json:"endpoint"`
	EndpointVersion string         `json:"endpoint_version"`
	Mode            string         `json:"mode"`
	ArtifactId      int            `json:"artifact_id"`
	ItemCount       int            `json:"item_count"`
	OutputCount     int            `json:"output_count"`
	StatusCode      int            `json:"status_code"`
	Latency         float64        `json:"latency"`
	Timestamp       time.Time      `json:"timestamp"`
	InputParams     map[string]any `json:"input_params"`
	ScenarioId      int            `json:"scenario_id"`
}
