package main

type snapshotContents struct {
	Data []snapshot `json:"data"`
}

type snapshot struct {
	Voltage     float64 `json:"voltage"`
	Motion      int     `json:"motion"`
	Temperature float64 `json:"temperature"`
	Sound       float64 `json:"sound"`
	Voc         int     `json:"voc"`
	Illuminance float64 `json:"illuminace"`
	Humidity    float64 `json:"humidity"`
	Timestamp   string  `json:"timestamp"`
}