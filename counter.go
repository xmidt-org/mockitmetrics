// SPDX-FileCopyrightText: 2023 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package mockitmetrics

import (
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
	value          map[string]float64
	panic          func(any)
	delimiter      string
	m              sync.Mutex
	root           *Counter
	expectedLabels *[]string
	lvp            []tuple
}

var _ kit.Counter = (*Counter)(nil)

// With returns a new counter with the provided label values.
func (c *Counter) With(labelValues ...string) kit.Counter {
	root := c
	if c.root != nil {
		root = c.root
	}

	lvp, err := convert(labelValues)
	if err != nil {
		goto failure
	}

	lvp = append(c.lvp, lvp...)

	err = validateLabels(root.expectedLabels, lvp, false)
	if err != nil {
		goto failure
	}

	return &Counter{
		root: root,
		lvp:  lvp,
	}

failure:
	root.panic(err)
	return nil
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

	if err := validateLabels(root.expectedLabels, c.lvp, true); err != nil {
		root.panic(err)
		return
	}

	label := joinValues(c.lvp, root.delimiter)

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
