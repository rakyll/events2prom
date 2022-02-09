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
	"sync"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/rakyll/events-to-prom/event"
)

var _ Processor = &GaugeProcessor{}

type gaugeSample struct {
	labelValues []string
	value       float64
}

type GaugeProcessor struct {
	col Collection

	samplesMu sync.RWMutex
	samples   map[string]gaugeSample

	prometheusDesc *prometheus.Desc
}

func NewGaugeProcessor(c Collection) *GaugeProcessor {
	return &GaugeProcessor{
		col:            c,
		samples:        make(map[string]gaugeSample, 64),
		prometheusDesc: prometheus.NewDesc(c.Name, c.Description, c.Labels, nil),
	}
}

func (p *GaugeProcessor) Collection() Collection {
	return p.col
}

func (p *GaugeProcessor) Handle(events []event.Event) {
	p.samplesMu.Lock()
	defer p.samplesMu.Unlock()

	for _, e := range events {
		if isMatch(e, p.col.Event, p.col.Labels) {
			labelValse := make([]string, len(p.col.Labels))
			for i, label := range p.col.Labels {
				labelValse[i] = e.Labels[label]
			}
			key := mapKeyForSample(p.col.Labels, labelValse)
			p.samples[key] = gaugeSample{
				labelValues: labelValse,
				value:       e.Value,
			}
		}
	}
}

func (p *GaugeProcessor) Describe(ch chan<- *prometheus.Desc) {
	ch <- p.prometheusDesc
}

func (p *GaugeProcessor) Collect(ch chan<- prometheus.Metric) {
	p.samplesMu.RLock()
	defer p.samplesMu.RUnlock()

	for _, sample := range p.samples {
		ch <- prometheus.MustNewConstMetric(
			p.prometheusDesc,
			prometheus.CounterValue,
			sample.value,
			sample.labelValues...,
		)
	}
}
