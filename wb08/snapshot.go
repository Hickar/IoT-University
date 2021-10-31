package main

import (
	"fmt"
	"strconv"
)

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

func (s *snapshot) ToCSVRow() []string {
	return []string{
		fmt.Sprintf("%f", s.Voltage),
		strconv.Itoa(s.Motion),
		fmt.Sprintf("%f", s.Temperature),
		fmt.Sprintf("%f", s.Sound),
		strconv.Itoa(s.Voc),
		fmt.Sprintf("%f", s.Illuminance),
		fmt.Sprintf("%f", s.Humidity),
		s.Timestamp,
	}
}