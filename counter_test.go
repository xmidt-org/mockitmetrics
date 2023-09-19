// SPDX-FileCopyrightText: 2023 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package mockitmetrics

import (
	"testing"

	kit "github.com/go-kit/kit/metrics"
	"github.com/stretchr/testify/assert"
)

func TestCounter(t *testing.T) {
	tests := []struct {
		description string
		fn          func(kit.Counter)
		opt         Option
		opts        []Option
		expected    map[string]float64
		expectPanic bool
	}{
		{
			description: "use a counter with no labels",
			fn: func(c kit.Counter) {
				c.Add(1)
			},
			expected: map[string]float64{
				"": 1.0,
			},
		}, {
			description: "use a counter with labels",
			fn: func(c kit.Counter) {
				c.With("label1", "value1", "label2", "value2").Add(1)
			},
			expected: map[string]float64{
				"value1.value2": 1.0,
			},
		}, {
			description: "use a counter with labels using 2 calls",
			fn: func(c kit.Counter) {
				c.With("label1", "value1").With("label2", "value2").Add(1)
			},
			expected: map[string]float64{
				"value1.value2": 1.0,
			},
		}, {
			description: "use a counter with labels using 3 calls",
			fn: func(c kit.Counter) {
				c.With("label1", "value1").With("label2", "value2").With("label3", "value3").Add(1)
			},
			expected: map[string]float64{
				"value1.value2.value3": 1.0,
			},
		}, {
			description: "use a different delimiter",
			fn: func(c kit.Counter) {
				c.With("label1", "value1").With("label2", "value2").With("label3", "value3").Add(1)
			},
			opt: Delimiter("-"),
			expected: map[string]float64{
				"value1-value2-value3": 1.0,
			},
		}, {
			description: "check that custom panic is honored",
			fn: func(c kit.Counter) {
				c.With("label1", "value1").Add(-1)
			},
			opt: PanicFunc(func(any) {}),
		}, {
			description: "check that panic is honored",
			fn: func(c kit.Counter) {
				c.With("label1", "value1").Add(-1)
			},
			expectPanic: true,
		}, {
			description: "use a counter with no labels, expecting no labels",
			fn: func(c kit.Counter) {
				c.Add(1)
			},
			opt: ExpectLabels(),
			expected: map[string]float64{
				"": 1.0,
			},
		}, {
			description: "error when an unexpected label is sent",
			fn: func(c kit.Counter) {
				c.With("invalid").Add(1)
			},
			opt:         ExpectLabels(),
			expectPanic: true,
		}, {
			description: "use a counter with labels, and require 2",
			fn: func(c kit.Counter) {
				c.With("one", "value1", "two", "value2").Add(1)
			},
			opt: ExpectLabels("one", "two"),
			expected: map[string]float64{
				"value1.value2": 1.0,
			},
		}, {
			description: "error when a missing label is sent",
			fn: func(c kit.Counter) {
				c.With("one", "value1").Add(1)
			},
			opt:         ExpectLabels("one", "two"),
			expectPanic: true,
		}, {
			description: "error when an extra label is sent",
			fn: func(c kit.Counter) {
				c.With("one", "value1", "two", "value2", "label3", "value3")
			},
			opt:         ExpectLabels("one", "two"),
			expectPanic: true,
		}, {
			description: "use a gauge with the wrong label",
			fn: func(c kit.Counter) {
				c.With("label1", "value1", "label2", "value2")
			},
			opt:         ExpectLabels("one", "two"),
			expectPanic: true,
		}, {
			description: "use a gauge missing a value",
			fn: func(c kit.Counter) {
				c.With("one")
			},
			opt:         ExpectLabels("one", "two"),
			expectPanic: true,
		}, {
			description: "use a gauge with the label set to ''",
			fn: func(c kit.Counter) {
				c.With("", "value")
			},
			opt:         ExpectLabels("one", "two"),
			expectPanic: true,
		}, {
			description: "use a gauge with the value set to ''",
			fn: func(c kit.Counter) {
				c.With("one", "")
			},
			opt:         ExpectLabels("one", "two"),
			expectPanic: true,
		},
	}

	for _, tc := range tests {
		t.Run(tc.description, func(t *testing.T) {
			assert := assert.New(t)

			opts := append(tc.opts, tc.opt)
			c := NewCounter(opts...)
			if tc.expectPanic {
				assert.Panics(func() { tc.fn(c) })
				return
			}

			tc.fn(c)

			assert.Equal(tc.expected, c.Value())
		})
	}
}
