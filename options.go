// SPDX-FileCopyrightText: 2023 Comcast Cable Communications Management, LLC
// SPDX-License-Identifier: Apache-2.0

package mockitmetrics

const (
	DelimiterDefault = "."
	NoLabelDefault   = "none"
)

type Option interface {
	counterApply(*Counter)
	gaugeApply(*Gauge)
	histogramApply(*Histogram)
}

// Delimiter sets the delimiter used to join labels.
func Delimiter(d string) Option {
	return delimiter(d)
}

type delimiter string

func (d delimiter) counterApply(c *Counter) {
	c.delimiter = string(d)
}

func (d delimiter) gaugeApply(g *Gauge) {
	g.delimiter = string(d)
}

func (d delimiter) histogramApply(h *Histogram) {
	h.delimiter = string(d)
}

// PanicFunc sets the function to call when panic() would be called.
func PanicFunc(f func(any)) Option {
	return panicFunc(f)
}

type panicFunc func(any)

func (f panicFunc) counterApply(c *Counter) {
	c.panic = f
}

func (f panicFunc) gaugeApply(g *Gauge) {
	g.panic = f
}

func (f panicFunc) histogramApply(h *Histogram) {
	h.panic = f
}

// ExpectLabels sets the labels that are expected to be passed to the metric.
//
// The labels aren't validated against, but provide the number of labels that
// are expected to be passed to the metric.
//
// If the number of labels passed to the metric doesn't match the number of
// labels passed to ExpectLabels, the call to update the metric will panic.
func ExpectLabels(labels ...string) Option {
	return expectLabels{labels: labels}
}

type expectLabels struct {
	labels []string
}

func (e expectLabels) counterApply(c *Counter) {
	c.labels = &e.labels
}

func (e expectLabels) gaugeApply(g *Gauge) {
	g.labels = &e.labels
}

func (e expectLabels) histogramApply(h *Histogram) {
	h.labels = &e.labels
}
