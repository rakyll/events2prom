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

func TestCount(t *testing.T) {
	p := NewCountProcessor(Collection{
		Name:        "request_latency_ms",
		Description: "Request latency in ms",
		Event:       "request_latency_ms",
		Labels:      []string{"region", "az"},
	})
	p.Handle([]event.Event{
		{
			Name: "request_latency_ms",
			Labels: map[string]string{
				"region":  "us-east-1",
				"az":      "us-east-1c",
				"service": "logging",
			},
			Timestamp: time.Now(),
			Value:     100,
		},
		{
			Name: "request_latency_ms",
			Labels: map[string]string{
				"region":  "us-east-1",
				"az":      "us-east-1c",
				"service": "logging",
			},
			Timestamp: time.Now(),
			Value:     200,
		},
		{
			Name: "request_latency_ms",
			Labels: map[string]string{
				"region": "us-west-1",
				"az":     "us-west-1c",
			},
			Timestamp: time.Now(),
			Value:     5,
		},
	})

	assert.Equal(t,
		p.samples["region_us-east-1_az_us-east-1c_"].count,
		uint64(2))

	assert.Equal(t,
		p.samples["region_us-west-1_az_us-west-1c_"].count,
		uint64(1))
}

func BenchmarkCount(b *testing.B) {
	p := NewCountProcessor(Collection{
		Name:        "request_latency_ms_by_region",
		Description: "Request latency in ms by region",
		Event:       "request_latency_ms",
		Labels:      []string{"region"},
	})
	now := time.Now()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		p.Handle([]event.Event{
			{
				Name: "request_latency_ms",
				Labels: map[string]string{
					"region":  "us-east-1",
					"az":      "us-east-1c",
					"service": "logging",
				},
				Timestamp: now,
				Value:     100,
			},
			{
				Name: "request_latency_ms",
				Labels: map[string]string{
					"az":     "us-east-1c",
					"region": "us-west-1",
				},
				Timestamp: now,
				Value:     5,
			},
		})
	}
}
