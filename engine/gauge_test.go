// Copyright 2015 The Prometheus Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package engine

import (
	"testing"
	"time"

	"github.com/rakyll/events2prom/event"
	"github.com/stretchr/testify/assert"
)

func TestGauge(t *testing.T) {
	p := NewGaugeProcessor(Collection{
		Name:        "cpu_total",
		Description: "CPU time",
		Event:       "cpu",
		Labels:      []string{"region", "az"},
	})
	p.Handle([]event.Event{
		{
			Name: "cpu",
			Labels: map[string]string{
				"region":  "us-east-1",
				"az":      "us-east-1c",
				"service": "logging",
			},
			Timestamp: time.Now(),
			Value:     300.5,
		},
		{
			Name: "cpu",
			Labels: map[string]string{
				"region":  "us-east-1",
				"az":      "us-east-1c",
				"service": "logging",
			},
			Timestamp: time.Now(),
			Value:     312,
		},
		{
			Name: "cpu",
			Labels: map[string]string{
				"region": "us-west-1",
				"az":     "us-west-1c",
			},
			Timestamp: time.Now(),
			Value:     23.0,
		},
	})

	assert.Equal(t,
		p.samples["region_us-east-1_az_us-east-1c_"].value, 312.0)

	assert.Equal(t,
		p.samples["region_us-west-1_az_us-west-1c_"].value, 23.0)
}
