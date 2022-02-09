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
	"bytes"
	"log"
	"sort"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rakyll/events-to-prom/event"
)

// TODO(jbd): We only support cumulative for now, think about deltas with a time window.

const defaultBufferSize = 32 * 1024

type Processor interface {
	prometheus.Collector
	Handle(events []event.Event)
	Collection() Collection
}

type Collection struct {
	Name        string    `json:"name,omitempty" yaml:"name,omitempty"`
	Description string    `json:"description,omitempty" yaml:"description,omitempty"`
	Aggregation string    `json:"aggregation,omitempty" yaml:"aggregation,omitempty"` // count, sum or histogram
	Event       string    `json:"event,omitempty" yaml:"event,omitempty"`
	Labels      []string  `json:"labels,omitempty" yaml:"labels,omitempty"`
	Buckets     []float64 `json:"buckets,omitempty" yaml:"buckets,omitempty"` // only if aggregation is histogram, otherwise ignored
}

type Loop struct {
	processors map[string]Processor // access only in Run
	allEvents  map[string]struct{}  // access only in Run
	// TODO(jbd): Think about organizing processors by event.

	buffer            []event.Event // access only in Run
	bufferIndex       int           // access only in Run
	BufferFlushWindow time.Duration
	maxBufferSize     int

	incomingEvents <-chan event.Event
	newCollections <-chan Collection
	removals       <-chan string

	promRegistry *prometheus.Registry
}

func NewLoop(bufferSize int, e <-chan event.Event, c <-chan Collection, r <-chan string) *Loop {
	if bufferSize == 0 {
		bufferSize = defaultBufferSize
	}
	return &Loop{
		processors: make(map[string]Processor),
		allEvents:  make(map[string]struct{}),

		buffer:            make([]event.Event, bufferSize),
		maxBufferSize:     bufferSize,
		bufferIndex:       0,
		BufferFlushWindow: 5 * time.Second,
		incomingEvents:    e,
		newCollections:    c,
		removals:          r,
		promRegistry:      prometheus.NewRegistry(),
	}
}

func (l *Loop) Run() {
	timer := time.NewTimer(l.BufferFlushWindow)
	defer timer.Stop()

	for {
		select {
		case c := <-l.newCollections:
			l.enableCollection(c)
		case name := <-l.removals:
			l.disableCollection(name)
		case e := <-l.incomingEvents:
			l.incomingEvent(e)
			if l.bufferIndex == l.maxBufferSize {
				l.flush(timer)
			}
		case <-timer.C:
			if l.bufferIndex != 0 {
				l.flush(timer)
			}
		}
	}
}

// enableCollection should only be called from Run.
func (l *Loop) enableCollection(c Collection) {
	name := c.Name
	if name == "" {
		log.Println("Failed to enable collection with empty name")
		return
	}
	if c.Event == "" {
		log.Println("Failed to enable collection with empty event")
		return
	}
	_, ok := l.processors[name]
	if ok {
		log.Printf("Failed to enable duplicated collection: %q", name)
		return
	}

	// TODO(jbd): Support sum.
	var p Processor
	switch c.Aggregation {
	case "count":
		p = NewCountProcessor(c)
	case "sum":
		p = NewSumProcessor(c)
	case "gauge":
		p = NewGaugeProcessor(c)
	case "histogram":
		if len(c.Buckets) == 0 {
			log.Printf("Failed to enable %q with no buckets", c.Name)
			return
		}
		sort.Float64s(c.Buckets)
		p = NewHistogramProcessor(c)
	default:
		log.Printf("Unknown aggregation (%q) for %q", c.Aggregation, c.Name)
	}

	l.processors[name] = p
	l.promRegistry.MustRegister(p)
	l.allEvents[c.Event] = struct{}{}
	log.Printf("Enabled collection: %q", name)
}

// disableCollection should only be called from Run.
func (l *Loop) disableCollection(name string) {
	p, ok := l.processors[name]
	if !ok {
		return
	}

	delete(l.processors, name)
	l.promRegistry.Unregister(p)
	delete(l.allEvents, p.Collection().Event)

	log.Printf("Disabled collection: %q", name)
}

// incomingEvent should only be called from Run.
func (l *Loop) incomingEvent(e event.Event) {
	// Ignore incoming events if there are no processors.
	if len(l.processors) == 0 {
		return
	}
	// Ignore incoming event if it's not currently collected.
	_, ok := l.allEvents[e.Name]
	if !ok {
		return
	}

	l.buffer[l.bufferIndex] = e
	l.bufferIndex++
}

// flush should only be called from Run.
func (l *Loop) flush(timer *time.Timer) {
	timer.Reset(l.BufferFlushWindow)
	events := l.buffer[:l.bufferIndex]
	for _, p := range l.processors {
		p.Handle(events)
	}
	log.Printf("Flushed %d events.", l.bufferIndex)
	l.bufferIndex = 0
}

func (l *Loop) Registry() *prometheus.Registry {
	return l.promRegistry
}

func mapKeyForSample(labels, labelValues []string) string {
	// Note: labels and labelValues are already sorted.
	var buf bytes.Buffer
	for i, label := range labels {
		buf.WriteString(label)
		buf.WriteByte('_')
		buf.WriteString(labelValues[i])
		buf.WriteByte('_')
	}
	return buf.String()
}

func isMatch(e event.Event, name string, labels []string) bool {
	if name != e.Name {
		return false
	}
	if len(e.Labels) < len(labels) {
		return false
	}
	for _, label := range labels {
		if _, ok := e.Labels[label]; !ok {
			return false
		}
	}
	return true
}
