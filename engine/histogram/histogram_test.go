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

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestEmptyHistogram(t *testing.T) {
	h := NewHistogram([]float64{100, 200, 300, 400, 500, 600, 700, 800, 900, 1000})
	for _, v := range h.Buckets() {
		assert.Equal(t, v, uint64(0))
	}
	assert.Equal(t, h.Total(), uint64(0))
}

func TestHistogram(t *testing.T) {
	h := NewHistogram([]float64{100, 200, 300, 400, 500, 600, 700, 800, 900, 1000})
	for i := float64(1); i <= 1000; i++ {
		h.Add(i)
	}
	assert.Equal(t, h.Buckets(), map[float64]uint64{
		100.0:  100,
		200.0:  200,
		300.0:  300,
		400.0:  400,
		500.0:  500,
		600.0:  600,
		700.0:  700,
		800.0:  800,
		900.0:  900,
		1000.0: 1000,
	})
	assert.Equal(t, h.Total(), uint64(1000))
}

func BenchmarkAdd(b *testing.B) {
	h := NewHistogram([]float64{100, 200, 300, 400, 500, 600, 700, 800, 900, 1000})
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		for j := 0; j < 1000; j++ {
			h.Add(float64(j))
		}
	}
}
