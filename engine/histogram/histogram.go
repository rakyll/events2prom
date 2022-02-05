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

package histogram

type Histogram struct {
	buckets []float64
	counts  []uint64
	sum     float64
}

func NewHistogram(b []float64) *Histogram {
	return &Histogram{
		buckets: b,
		counts:  make([]uint64, len(b)),
	}
}

func (h *Histogram) Add(v float64) {
	for i, b := range h.buckets {
		if v <= b {
			h.counts[i]++
			break
		}
	}
	h.sum += v
}

func (h *Histogram) Buckets() map[float64]uint64 {
	m := make(map[float64]uint64)
	var le uint64
	for i, b := range h.buckets {
		le += h.counts[i]
		m[b] = le
	}
	return m
}

func (h *Histogram) Total() uint64 {
	var total uint64
	for _, c := range h.counts {
		total += c
	}
	return total
}

func (h *Histogram) Sum() float64 {
	return h.sum
}
