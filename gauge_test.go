// SPDX-FileCopyrightText: 2023 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package mockitmetrics

import (
	"testing"

	kit "github.com/go-kit/kit/metrics"
	"github.com/stretchr/testify/assert"
)

func TestGauge(t *testing.T) {
	tests := []struct {
		description string
		fn          func(kit.Gauge)
		opt         Option
		opts        []Option
		expected    map[string]float64
		expectPanic bool
	}{
		{
			description: "use a gauge with no labels",
			fn: func(g kit.Gauge) {
				g.Add(1)
			},
			expected: map[string]float64{
				"": 1.0,
			},
		}, {
			description: "use a gauge with labels",
			fn: func(g kit.Gauge) {
				g.With("label1", "value1", "label2", "value2").Add(1)
			},
			expected: map[string]float64{
				"value1.value2": 1.0,
			},
		}, {
			description: "use a gauge with labels using 2 calls",
			fn: func(g kit.Gauge) {
				g.With("label1", "value1").With("label2", "value2").Add(1)
			},
			expected: map[string]float64{
				"value1.value2": 1.0,
			},
		}, {
			description: "use a gauge with labels using 3 calls",
			fn: func(g kit.Gauge) {
				g.With("label1", "value1").With("label2", "value2").With("label3", "value3").Add(1)
				g.With("label7", "value7").With("label2", "value2").With("label9", "value9").Set(99)
				g.With("label7", "value7").With("label2", "value2").With("label9", "value9").Add(1)
			},
			expected: map[string]float64{
				"value1.value2.value3": 1.0,
				"value7.value2.value9": 100.0,
			},
		}, {
			description: "use a different delimiter",
			fn: func(g kit.Gauge) {
				g.With("label1", "value1").With("label2", "value2").With("label3", "value3").Add(1)
				g.With("label7", "value7").With("label2", "value2").With("label9", "value9").Set(99)
				g.With("label7", "value7").With("label2", "value2").With("label9", "value9").Add(1)
			},
			opt: Delimiter("-"),
			expected: map[string]float64{
				"value1-value2-value3": 1.0,
				"value7-value2-value9": 100.0,
			},
		}, {
			description: "output an empty gauge",
			fn:          func(h kit.Gauge) {},
		}, {
			description: "check that custom panic is honored",
			fn: func(g kit.Gauge) {
				g.With("label1", "value1")
			},
			opts: []Option{PanicFunc(func(any) {}), ExpectLabels()},
		}, {
			description: "check that panic is honored",
			fn: func(g kit.Gauge) {
				g.With("label1", "value1").Add(-1)
			},
			opt:         ExpectLabels(),
			expectPanic: true,
		}, {
			description: "use a gauge with no labels, expecting no labels",
			fn: func(g kit.Gauge) {
				g.Add(1)
			},
			opt: ExpectLabels(),
			expected: map[string]float64{
				"": 1.0,
			},
		}, {
			description: "error when an unexpected label is sent",
			fn: func(g kit.Gauge) {
				g.With("invalid", "value1").Add(1)
			},
			opt:         ExpectLabels(),
			expectPanic: true,
		}, {
			description: "use a gauge with labels, and require 2",
			fn: func(g kit.Gauge) {
				g.With("one", "value1", "two", "value2").Add(1)
			},
			opt: ExpectLabels("one", "two"),
			expected: map[string]float64{
				"value1.value2": 1.0,
			},
		}, {
			description: "error when a missing label is sent",
			fn: func(g kit.Gauge) {
				g.With("one", "value1").Add(1)
			},
			opt:         ExpectLabels("one", "two"),
			expectPanic: true,
		}, {
			description: "error when an extra label is sent",
			fn: func(g kit.Gauge) {
				g.With("label1", "label2", "label3").Add(1)
			},
			opt:         ExpectLabels("one", "two"),
			expectPanic: true,
		}, {
			description: "use a gauge with the wrong label",
			fn: func(g kit.Gauge) {
				g.With("label1", "value1", "label2", "value2")
			},
			opt:         ExpectLabels("one", "two"),
			expectPanic: true,
		}, {
			description: "use a gauge missing a value",
			fn: func(g kit.Gauge) {
				g.With("one")
			},
			opt:         ExpectLabels("one", "two"),
			expectPanic: true,
		}, {
			description: "use a gauge with the label set to ''",
			fn: func(g kit.Gauge) {
				g.With("", "value")
			},
			opt:         ExpectLabels("one", "two"),
			expectPanic: true,
		}, {
			description: "use a gauge with the value set to ''",
			fn: func(g kit.Gauge) {
				g.With("one", "")
			},
			opt:         ExpectLabels("one", "two"),
			expectPanic: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			assert := assert.New(t)

			opts := append(tc.opts, tc.opt)
			g := NewGauge(opts...)
			if tc.expectPanic {
				assert.Panics(func() { tc.fn(g) })
				return
			}

			tc.fn(g)

			assert.Equal(tc.expected, g.Value())
		})
	}
}
