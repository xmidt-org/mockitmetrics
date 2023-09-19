// SPDX-FileCopyrightText: 2023 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package mockitmetrics

import (
	"testing"

	kit "github.com/go-kit/kit/metrics"
	"github.com/stretchr/testify/assert"
)

func TestHistogram(t *testing.T) {
	tests := []struct {
		description string
		fn          func(kit.Histogram)
		opt         Option
		opts        []Option
		expected    map[string][]float64
		expectPanic bool
	}{
		{
			description: "use a histogram with no labels",
			fn: func(h kit.Histogram) {
				h.Observe(1)
			},
			expected: map[string][]float64{
				"": {1.0},
			},
		}, {
			description: "use a histogram with labels",
			fn: func(h kit.Histogram) {
				h.With("label1", "value1", "label2", "value2").Observe(1)
			},
			expected: map[string][]float64{
				"value1.value2": {1.0},
			},
		}, {
			description: "use two histograms with labels",
			fn: func(h kit.Histogram) {
				h.With("label1", "value1", "label2", "value2").Observe(1)
				h.With("label1", "value1", "label2", "value2").Observe(10)
			},
			expected: map[string][]float64{
				"value1.value2": {1.0, 10.0},
			},
		}, {
			description: "use a histogram with labels using 2 calls",
			fn: func(h kit.Histogram) {
				h.With("label1", "value1").With("label2", "value2").Observe(1)
			},
			expected: map[string][]float64{
				"value1.value2": {1.0},
			},
		}, {
			description: "use a histogram with labels using 3 calls",
			fn: func(h kit.Histogram) {
				h.With("label1", "value1").With("label2", "value2").With("label3", "value3").Observe(1)
			},
			expected: map[string][]float64{
				"value1.value2.value3": {1.0},
			},
		}, {
			description: "use a different delimiter",
			fn: func(h kit.Histogram) {
				h.With("label1", "value1").With("label2", "value2").With("label3", "value3").Observe(1)
			},
			opt: Delimiter("-"),
			expected: map[string][]float64{
				"value1-value2-value3": {1.0},
			},
		}, {
			description: "output an empty histogram",
			fn:          func(h kit.Histogram) {},
		}, {
			description: "use a counter with no labels, expecting no labels",
			fn: func(h kit.Histogram) {
				h.Observe(1)
			},
			opt: ExpectLabels(),
			expected: map[string][]float64{
				"": {1.0},
			},
		}, {
			description: "error when an unexpected label is sent",
			fn: func(h kit.Histogram) {
				h.With("invalid")
			},
			opt:         ExpectLabels(),
			expectPanic: true,
		}, {
			description: "honor the custom panic function",
			fn: func(h kit.Histogram) {
				h.With("invalid", "value1")
			},
			opts: []Option{ExpectLabels(), PanicFunc(func(any) {})},
		}, {
			description: "use a counter with labels, and require 2",
			fn: func(h kit.Histogram) {
				h.With("one", "value1", "two", "value2").Observe(1)
			},
			opt: ExpectLabels("one", "two"),
			expected: map[string][]float64{
				"value1.value2": {1.0},
			},
		}, {
			description: "error when a missing label is sent",
			fn: func(h kit.Histogram) {
				h.With("one", "value1").Observe(1)
			},
			opt:         ExpectLabels("one", "two"),
			expectPanic: true,
		}, {
			description: "error when an extra label is sent",
			fn: func(h kit.Histogram) {
				h.With("label1", "label2", "label3").Observe(1)
			},
			opt:         ExpectLabels("one", "two"),
			expectPanic: true,
		}, {
			description: "use a gauge with the wrong label",
			fn: func(h kit.Histogram) {
				h.With("label1", "value1", "label2", "value2")
			},
			opt:         ExpectLabels("one", "two"),
			expectPanic: true,
		}, {
			description: "use a gauge missing a value",
			fn: func(h kit.Histogram) {
				h.With("one")
			},
			opt:         ExpectLabels("one", "two"),
			expectPanic: true,
		}, {
			description: "use a gauge with the label set to ''",
			fn: func(h kit.Histogram) {
				h.With("", "value")
			},
			opt:         ExpectLabels("one", "two"),
			expectPanic: true,
		}, {
			description: "use a gauge with the value set to ''",
			fn: func(h kit.Histogram) {
				h.With("one", "")
			},
			opt:         ExpectLabels("one", "two"),
			expectPanic: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			assert := assert.New(t)

			opts := append(tc.opts, tc.opt)
			h := NewHistogram(opts...)
			if tc.expectPanic {
				assert.Panics(func() { tc.fn(h) })
				return
			}
			tc.fn(h)

			assert.Equal(tc.expected, h.Value())
		})
	}
}
