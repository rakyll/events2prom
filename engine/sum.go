// Copyright 2022 The Prometheus Authors
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
	"github.com/rakyll/events2prom/event"
)

var _ Processor = &SumProcessor{}

type sumSample struct {
	labelValues []string
	sum         float64
}

type SumProcessor struct {
	col Collection

	samplesMu sync.RWMutex
	samples   map[string]sumSample

	prometheusDesc *prometheus.Desc
}

func NewSumProcessor(c Collection) *SumProcessor {
	return &SumProcessor{
		col:            c,
		samples:        make(map[string]sumSample, 64),
		prometheusDesc: prometheus.NewDesc(c.Name, c.Description, c.Labels, nil),
	}
}

func (p *SumProcessor) Collection() Collection {
	return p.col
}

func (p *SumProcessor) Handle(events []event.Event) {
	p.samplesMu.Lock()
	defer p.samplesMu.Unlock()

	col := p.col
	for _, e := range events {
		if isMatch(e, col.Event, col.Labels) {
			key, labelVals := generateKeyLabelVals(p.col, e)
			_, ok := p.samples[key]
			if !ok {
				p.samples[key] = sumSample{
					labelValues: labelVals,
				}
			}
			s := p.samples[key]
			s.sum += e.Value
			p.samples[key] = s
		}
	}
}

func (p *SumProcessor) Describe(ch chan<- *prometheus.Desc) {
	ch <- p.prometheusDesc
}

func (p *SumProcessor) Collect(ch chan<- prometheus.Metric) {
	p.samplesMu.RLock()
	defer p.samplesMu.RUnlock()

	for _, sample := range p.samples {
		ch <- prometheus.MustNewConstMetric(
			p.prometheusDesc,
			// sum is not a first class data type in Promehteus,
			// use a gauge instead.
			prometheus.GaugeValue,
			sample.sum,
			sample.labelValues...,
		)
	}
}
