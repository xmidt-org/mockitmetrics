// SPDX-FileCopyrightText: 2023 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package mockitmetrics

import (
	"sync"

	kit "github.com/go-kit/kit/metrics"
)

// NewGauge creates a new gauge with the provided options.
func NewGauge(opts ...Option) *Gauge {
	g := Gauge{
		delimiter: DelimiterDefault,
		panic:     func(a any) { panic(a) },
	}

	for _, opt := range opts {
		if opt != nil {
			opt.gaugeApply(&g)
		}
	}

	return &g
}

// Gauge is a mock gauge.
type Gauge struct {
	value          map[string]float64
	delimiter      string
	panic          func(any)
	m              sync.Mutex
	root           *Gauge
	expectedLabels *[]string
	lvp            []tuple
}

var _ kit.Gauge = (*Gauge)(nil)

// With returns a new gauge with the provided label values.
func (g *Gauge) With(labelValues ...string) kit.Gauge {
	root := g
	if g.root != nil {
		root = g.root
	}

	lvp, err := convert(labelValues)
	if err != nil {
		goto failure
	}

	lvp = append(g.lvp, lvp...)

	err = validateLabels(root.expectedLabels, lvp, false)
	if err != nil {
		goto failure
	}

	return &Gauge{
		root: root,
		lvp:  lvp,
	}

failure:
	root.panic(err)
	return nil
}

func (g *Gauge) update(value float64, delta bool) {
	root := g.root
	if root == nil {
		root = g
	}

	if err := validateLabels(root.expectedLabels, g.lvp, true); err != nil {
		root.panic(err)
		return
	}

	label := joinValues(g.lvp, root.delimiter)

	root.m.Lock()
	defer root.m.Unlock()

	if root.value == nil {
		root.value = map[string]float64{}
	}

	if _, ok := root.value[label]; !ok {
		root.value[label] = 0.0
	}

	if delta {
		root.value[label] += value
	} else {
		root.value[label] = value
	}
}

// Set sets the gauge to the provided value.
func (g *Gauge) Set(value float64) {
	g.update(value, false)
}

// Add adds the provided delta to the gauge.
func (g *Gauge) Add(delta float64) {
	g.update(delta, true)
}

// Value returns the current value of the gauge.
func (g *Gauge) Value() map[string]float64 {
	root := g.root
	if root == nil {
		root = g
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
