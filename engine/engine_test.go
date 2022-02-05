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

	"github.com/rakyll/events-to-prom/event"
)

func BenchmarkIsMatch(b *testing.B) {
	event := event.Event{
		Name: "request_latency_ms",
		Labels: map[string]string{
			"region":  "us-east-1",
			"az":      "us-east-1c",
			"service": "logging",
		},
		Timestamp: time.Now(),
		Value:     100,
	}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		isMatch(event, "request_latency_ms", []string{"region", "az"})
	}
}

func BenchmarkMapKeyForSample(b *testing.B) {
	labels := []string{"region", "az"}
	values := []string{"us-east-1", "us-east-1c"}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		mapKeyForSample(labels, values)
	}
}

func TestIsMatch(t *testing.T) {
	tests := []struct {
		name      string
		event     event.Event
		eventName string
		labels    []string
		want      bool
	}{
		{
			name: "unmatching name",
			event: event.Event{
				Name:   "request_latency_ms",
				Labels: map[string]string{},
				Value:  100,
			},
			eventName: "request_latency_ms1",
			want:      false,
		},
		{
			name: "no labels to match",
			event: event.Event{
				Name: "request_latency_ms",
				Labels: map[string]string{
					"foo": "bar",
				},
				Value: 100,
			},
			eventName: "request_latency_ms",
			labels:    nil,
			want:      true,
		},
		{
			name: "no labels",
			event: event.Event{
				Name:   "request_latency_ms",
				Labels: map[string]string{},
				Value:  100,
			},
			eventName: "request_latency_ms",
			labels:    []string{"region", "az"},
			want:      false,
		},
		{
			name: "not all labels",
			event: event.Event{
				Name: "request_latency_ms",
				Labels: map[string]string{
					"region": "us-east-1",
				},
				Value: 100,
			},
			eventName: "request_latency_ms",
			labels:    []string{"region", "az"},
			want:      false,
		},
		{
			name: "all labels",
			event: event.Event{
				Name: "request_latency_ms",
				Labels: map[string]string{
					"region": "us-east-1",
					"az":     "us-east-1c",
				},
				Value: 100,
			},
			eventName: "request_latency_ms",
			labels:    []string{"region", "az"},
			want:      true,
		},
		{
			name: "more labels",
			event: event.Event{
				Name: "request_latency_ms",
				Labels: map[string]string{
					"region": "us-east-1",
					"az":     "us-east-1c",
					"foo":    "bar",
				},
				Value: 100,
			},
			eventName: "request_latency_ms",
			labels:    []string{"region", "az"},
			want:      true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := isMatch(tt.event, tt.eventName, tt.labels); got != tt.want {
				t.Errorf("isMatch() = %v, want %v", got, tt.want)
			}
		})
	}
}
