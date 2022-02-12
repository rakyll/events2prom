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
	"github.com/rakyll/events2prom/engine/histogram"
	"github.com/rakyll/events2prom/event"
)

var _ Processor = &HistogramProcessor{}

type histogramSample struct {
	histogram   *histogram.Histogram
	labelValues []string
}

type HistogramProcessor struct {
	col Collection

	samplesMu sync.RWMutex
	samples   map[string]histogramSample

	prometheusDesc *prometheus.Desc
}

func NewHistogramProcessor(c Collection) *HistogramProcessor {
	return &HistogramProcessor{
		col:            c,
		samples:        make(map[string]histogramSample, 64),
		prometheusDesc: prometheus.NewDesc(c.Name, c.Description, c.Labels, nil),
	}
}

func (p *HistogramProcessor) Collection() Collection {
	return p.col
}

func (p *HistogramProcessor) Handle(events []event.Event) {
	p.samplesMu.Lock()
	defer p.samplesMu.Unlock()

	for _, e := range events {
		if isMatch(e, p.col.Event, p.col.Labels) {
			key, labelVals := generateKeyLabelVals(p.col, e)
			_, ok := p.samples[key]
			if !ok {
				p.samples[key] = histogramSample{
					histogram:   histogram.NewHistogram(p.col.Buckets),
					labelValues: labelVals,
				}
			}
			hist := p.samples[key].histogram
			hist.Add(e.Value)
			p.samples[key] = histogramSample{
				histogram:   hist,
				labelValues: labelVals,
			}
		}
	}
}

func (p *HistogramProcessor) Describe(ch chan<- *prometheus.Desc) {
	ch <- p.prometheusDesc
}

func (p *HistogramProcessor) Collect(ch chan<- prometheus.Metric) {
	p.samplesMu.RLock()
	defer p.samplesMu.RUnlock()

	for _, sample := range p.samples {
		ch <- prometheus.MustNewConstHistogram(
			p.prometheusDesc,
			sample.histogram.Total(),
			sample.histogram.Sum(),
			sample.histogram.Buckets(),
			sample.labelValues...,
		)
	}
}
