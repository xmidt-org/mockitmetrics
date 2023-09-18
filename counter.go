// SPDX-FileCopyrightText: 2023 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package mockitmetrics

import (
	"fmt"
	"strings"
	"sync"

	kit "github.com/go-kit/kit/metrics"
)

// NewCounter creates a new counter with the provided options.
func NewCounter(opts ...Option) *Counter {
	c := Counter{
		delimiter: DelimiterDefault,
		panic:     func(a any) { panic(a) },
	}

	for _, opt := range opts {
		if opt != nil {
			opt.counterApply(&c)
		}
	}

	return &c
}

// Counter is a mock counter.
type Counter struct {
	value       map[string]float64
	panic       func(any)
	delimiter   string
	m           sync.Mutex
	root        *Counter
	labels      *[]string
	labelValues []string
}

var _ kit.Counter = (*Counter)(nil)

// With returns a new counter with the provided label values.
func (c *Counter) With(labelValues ...string) kit.Counter {
	root := c
	if c.root != nil {
		root = c.root
	}

	return &Counter{
		root:        root,
		labelValues: append(c.labelValues, labelValues...),
	}
}

// Add adds the provided delta to the counter.
func (c *Counter) Add(delta float64) {
	root := c.root
	if root == nil {
		root = c
	}

	if delta < 0.0 {
		root.panic("delta must be non-negative")
		return
	}

	if root.labels != nil && len(c.labelValues) != len(*root.labels) &&
		!(len(*root.labels) == 0 && len(c.labelValues) == 1 && c.labelValues[0] == "") {

		s := fmt.Sprintf("incorrect number of label values. labels: '%s' (%d), values '%s' (%d)",
			strings.Join(*root.labels, "', '"), len(*root.labels),
			strings.Join(root.labelValues, "', '"), len(root.labelValues),
		)
		root.panic(s)
		return
	}

	label := strings.Join(c.labelValues, root.delimiter)

	root.m.Lock()
	defer root.m.Unlock()

	if root.value == nil {
		root.value = map[string]float64{}
	}

	if _, ok := root.value[label]; !ok {
		root.value[label] = 0.0
	}
	root.value[label] += delta
}

// Value returns the current value of the tree of counters.
func (c *Counter) Value() map[string]float64 {
	root := c.root
	if root == nil {
		root = c
	}

	root.m.Lock()
	defer root.m.Unlock()

	if len(root.value) == 0 {
		return nil
	}

	rv := map[string]float64{}

	for k, v := range root.value {
		rv[k] = v
	}
	return rv
}
