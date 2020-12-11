package main

import (
	"fmt"
	"strconv"
	"time"
)

// Spots collect a group of spots
type Spots []Spot

// Spot is an individual observation from a listener
type Spot struct {
	SpotID       int         `json:"spotId"`
	Activator    string      `json:"activator"`
	Frequency    string      `json:"frequency"`
	Mode         string      `json:"mode"`
	Reference    string      `json:"reference"`
	ParkName     string      `json:"parkName"`
	SpotTime     string      `json:"spotTime"`
	Spotter      string      `json:"spotter"`
	Comments     string      `json:"comments"`
	Source       string      `json:"source"`
	Invalid      interface{} `json:"invalid"`
	Name         string      `json:"name"`
	LocationDesc string      `json:"locationDesc"`
}

// ToClusterFormat reformats the spot for the dx cluster
func (s *Spot) ToClusterFormat() string {
	f, err := strconv.ParseFloat(s.Frequency, 32)
	if err != nil {
		f = 4242
	}

	const longForm = "2006-01-02T15:04:05"
	t, _ := time.Parse(longForm, s.SpotTime)
	w := 17 - len(s.Spotter) // "spotter: freq" fits in 17 chars
	return fmt.Sprintf("DX de %s:%*.1f  %-13.12s%-30.30s %sZ\a\a\x0c", s.Spotter, w, f, s.Activator, s.Comments, t.Format("1504"))
}
