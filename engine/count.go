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

var _ Processor = &CountProcessor{}

type countSample struct {
	labelValues []string
	count       uint64
}

type CountProcessor struct {
	col Collection

	samplesMu sync.RWMutex // TODO(jbd): We can switch to atomics here.
	samples   map[string]countSample

	prometheusDesc *prometheus.Desc
}

func NewCountProcessor(c Collection) *CountProcessor {
	return &CountProcessor{
		col:            c,
		samples:        make(map[string]countSample, 64),
		prometheusDesc: prometheus.NewDesc(c.Name, c.Description, c.Labels, nil),
	}
}

func (p *CountProcessor) Collection() Collection {
	return p.col
}

func (p *CountProcessor) Handle(events []event.Event) {
	p.samplesMu.Lock()
	defer p.samplesMu.Unlock()

	for _, e := range events {
		if isMatch(e, p.col.Event, p.col.Labels) {
			labelVals := make([]string, len(p.col.Labels))
			for i, label := range p.col.Labels {
				labelVals[i] = e.Labels[label]
			}
			key := mapKeyForSample(p.col.Labels, labelVals)
			_, ok := p.samples[key]
			if !ok {
				p.samples[key] = countSample{
					labelValues: labelVals,
				}
			}
			s := p.samples[key]
			s.count++
			p.samples[key] = s
		}
	}
}

func (p *CountProcessor) Describe(ch chan<- *prometheus.Desc) {
	ch <- p.prometheusDesc
}

func (p *CountProcessor) Collect(ch chan<- prometheus.Metric) {
	p.samplesMu.RLock()
	defer p.samplesMu.RUnlock()

	for _, sample := range p.samples {
		ch <- prometheus.MustNewConstMetric(
			p.prometheusDesc,
			prometheus.CounterValue,
			float64(sample.count),
			sample.labelValues...,
		)
	}
}
