// SPDX-FileCopyrightText: 2023 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package mockitmetrics

import (
	"fmt"
	"strings"
	"sync"

	kit "github.com/go-kit/kit/metrics"
)

// NewHistogram creates a new histogram with the provided options.
func NewHistogram(opts ...Option) *Histogram {
	h := Histogram{
		delimiter: DelimiterDefault,
		panic:     func(a any) { panic(a) },
	}

	for _, opt := range opts {
		if opt != nil {
			opt.histogramApply(&h)
		}
	}

	return &h
}

// Histogram is a mock histogram.
type Histogram struct {
	value       map[string][]float64
	delimiter   string
	panic       func(any)
	m           sync.Mutex
	root        *Histogram
	labels      *[]string
	labelValues []string
}

var _ kit.Histogram = (*Histogram)(nil)

// With returns a new histogram with the provided label values.
func (h *Histogram) With(labelValues ...string) kit.Histogram {
	root := h
	if h.root != nil {
		root = h.root
	}

	return &Histogram{
		root:        root,
		labelValues: append(h.labelValues, labelValues...),
	}
}

// Observe adds the provided value to the histogram.
func (h *Histogram) Observe(value float64) {
	root := h.root
	if root == nil {
		root = h
	}

	if root.labels != nil && len(h.labelValues) != len(*root.labels) &&
		!(len(*root.labels) == 0 && len(h.labelValues) == 1 && h.labelValues[0] == "") {

		s := fmt.Sprintf("incorrect number of label values. labels: '%s' (%d), values '%s' (%d)",
			strings.Join(*root.labels, "', '"), len(*root.labels),
			strings.Join(root.labelValues, "', '"), len(root.labelValues),
		)
		root.panic(s)
		return
	}

	label := strings.Join(h.labelValues, root.delimiter)

	root.m.Lock()
	defer root.m.Unlock()

	if root.value == nil {
		root.value = map[string][]float64{}
	}

	if _, ok := root.value[label]; !ok {
		root.value[label] = []float64{}
	}
	root.value[label] = append(root.value[label], value)
}

// Value returns the current value of the histogram.
func (h *Histogram) Value() map[string][]float64 {
	root := h.root
	if root == nil {
		root = h
	}

	root.m.Lock()
	defer root.m.Unlock()

	if len(root.value) == 0 {
		return nil
	}

	rv := map[string][]float64{}

	for k, v := range root.value {
		rv[k] = v
	}
	return rv
}
