package main

import "fmt"

type snapshotDataset struct {
	Data []snapshotDatasetEntry `json:"data"`
}

type snapshotDatasetEntry struct {
	Voltage     float64 `json:"voltage"`
	Motion      int     `json:"motion"`
	Temperature float64 `json:"temperature"`
	Sound       float64 `json:"sound"`
	Voc         int     `json:"voc"`
	Illuminance float64 `json:"illuminace"`
	Humidity    float64 `json:"humidity"`
	Timestamp   string  `json:"timestamp"`
}

type snapshot struct {
	Data []snapshotEntry `json:"data" xml:"data"`
}

type snapshotEntry struct {
	Motion      int     `json:"motion" xml:"motion"`
	Sound       float64 `json:"sound" xml:"sound"`
	Illuminance float64 `json:"illuminance" xml:"illuminance"`
	Temperature float64 `json:"temperature" xml:"temperature"`
}

func (s *snapshotEntry) String() string {
	return fmt.Sprintf("Motion: %d, Sound: %.2f, Illuminance: %.2f, Temperature: %.2f", s.Motion, s.Sound, s.Illuminance, s.Temperature)
}
