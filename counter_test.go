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
			description: "send a counter with no labels",
			fn: func(c kit.Counter) {
				c.Add(1)
			},
			expected: map[string]float64{
				"": 1.0,
			},
		}, {
			description: "send a counter with labels",
			fn: func(c kit.Counter) {
				c.With("label1", "label2").Add(1)
			},
			expected: map[string]float64{
				"label1.label2": 1.0,
			},
		}, {
			description: "send a counter with labels using 2 calls",
			fn: func(c kit.Counter) {
				c.With("label1").With("label2").Add(1)
			},
			expected: map[string]float64{
				"label1.label2": 1.0,
			},
		}, {
			description: "send a counter with labels using 3 calls",
			fn: func(c kit.Counter) {
				c.With("label1").With("label2").With("label3").Add(1)
			},
			expected: map[string]float64{
				"label1.label2.label3": 1.0,
			},
		}, {
			description: "use a different delimiter",
			fn: func(c kit.Counter) {
				c.With("label1").With("label2").With("label3").Add(1)
			},
			opt: Delimiter("-"),
			expected: map[string]float64{
				"label1-label2-label3": 1.0,
			},
		}, {
			description: "check that custom panic is honored",
			fn: func(c kit.Counter) {
				c.With("label1").Add(-1)
			},
			opt: PanicFunc(func(any) {}),
		}, {
			description: "check that panic is honored",
			fn: func(c kit.Counter) {
				c.With("label1").Add(-1)
			},
			expectPanic: true,
		}, {
			description: "send a counter with no labels, expecting no labels",
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
			description: "send a counter with labels, and require 2",
			fn: func(c kit.Counter) {
				c.With("label1", "label2").Add(1)
			},
			opt: ExpectLabels("one", "two"),
			expected: map[string]float64{
				"label1.label2": 1.0,
			},
		}, {
			description: "error when a missing label is sent",
			fn: func(c kit.Counter) {
				c.With("label1").Add(1)
			},
			opt:         ExpectLabels("one", "two"),
			expectPanic: true,
		}, {
			description: "error when an extra label is sent",
			fn: func(c kit.Counter) {
				c.With("label1", "label2", "label3").Add(1)
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
