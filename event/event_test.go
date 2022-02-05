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

package event

import (
	"encoding/json"
	"log"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func BenchmarkFastJSON(b *testing.B) {
	eventJSON := []byte(`{
		"event": "request_latency_ms",
		"labels": { 
			"foo1": "bar1",
			"foo2": "bar2",
			"foo3": "bar3",
			"foo4": "bar4",
			"foo5": "bar5"
		},
		"value": 54.7
	}`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := ParseJSON(eventJSON)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func TestParseText(t *testing.T) {
	eventText := []byte(`request_latency_ms|54.7|0|foo1:bar1|foo2:bar2|foo3:bar3|foo4:`)

	event, err := Parse(eventText)
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, event.Name, "request_latency_ms")
	assert.Equal(t, event.Value, 54.7)
	assert.Equal(t, event.Timestamp, time.Time{})
	assert.Equal(t, event.Labels["foo1"], "bar1")
	assert.Equal(t, event.Labels["foo2"], "bar2")
	assert.Equal(t, event.Labels["foo3"], "bar3")
	assert.Equal(t, event.Labels["foo4"], "")
}

func TestParseText_noLabels(t *testing.T) {
	eventText := []byte(`request_latency_ms|54.7|0`)

	event, err := Parse(eventText)
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, event.Name, "request_latency_ms")
	assert.Equal(t, event.Value, 54.7)
	assert.Equal(t, event.Timestamp, time.Time{})
}

func BenchmarkText(b *testing.B) {
	eventText := []byte(`request_latency_ms|54.7|0|foo1:bar1|foo2:bar2|foo3:bar3|foo4:bar4|foo5:bar5`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Parse(eventText)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func BenchmarkText_noLabels(b *testing.B) {
	eventText := []byte(`request_latency_ms|54.7|0`)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, err := Parse(eventText)
		if err != nil {
			log.Fatal(err)
		}
	}
}

func BenchmarkJSON(b *testing.B) {
	eventJSON := `{
		"event": "request_latency_ms",
		"labels": { 
			"foo1": "bar1",
			"foo2": "bar2",
			"foo3": "bar3",
			"foo4": "bar4",
			"foo5": "bar5"
		},
		"value": 54.7
	}`

	b.ResetTimer()
	var event Event
	for i := 0; i < b.N; i++ {
		if err := json.Unmarshal([]byte(eventJSON), &event); err != nil {
			log.Fatal(err)
		}
	}
}
