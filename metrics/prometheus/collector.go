// Copyright 2023 Bitnet
// This file is part of the Bitnet library.
//
// This software is provided "as is", without warranty of any kind,
// express or implied, including but not limited to the warranties
// of merchantability, fitness for a particular purpose and
// noninfringement. In no even shall the authors or copyright
// holders be liable for any claim, damages, or other liability,
// whether in an action of contract, tort or otherwise, arising
// from, out of or in connection with the software or the use or
// other dealings in the software.

package prometheus

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"

	"github.com/ethereum/go-ethereum/metrics"
)

var (
	typeGaugeTpl           = "# TYPE %s gauge\n"
	typeCounterTpl         = "# TYPE %s counter\n"
	typeSummaryTpl         = "# TYPE %s summary\n"
	keyValueTpl            = "%s %v\n\n"
	keyQuantileTagValueTpl = "%s {quantile=\"%s\"} %v\n"
)

// collector is a collection of byte buffers that aggregate Prometheus reports
// for different metric types.
type collector struct {
	buff *bytes.Buffer
}

// newCollector creates a new Prometheus metric aggregator.
func newCollector() *collector {
	return &collector{
		buff: &bytes.Buffer{},
	}
}

func (c *collector) addCounter(name string, m metrics.Counter) {
	c.writeGaugeCounter(name, m.Count())
}

func (c *collector) addCounterFloat64(name string, m metrics.CounterFloat64) {
	c.writeGaugeCounter(name, m.Count())
}

func (c *collector) addGauge(name string, m metrics.Gauge) {
	c.writeGaugeCounter(name, m.Value())
}

func (c *collector) addGaugeFloat64(name string, m metrics.GaugeFloat64) {
	c.writeGaugeCounter(name, m.Value())
}

func (c *collector) addHistogram(name string, m metrics.Histogram) {
	pv := []float64{0.5, 0.75, 0.95, 0.99, 0.999, 0.9999}
	ps := m.Percentiles(pv)
	c.writeSummaryCounter(name, m.Count())
	c.buff.WriteString(fmt.Sprintf(typeSummaryTpl, mutateKey(name)))
	for i := range pv {
		c.writeSummaryPercentile(name, strconv.FormatFloat(pv[i], 'f', -1, 64), ps[i])
	}
	c.buff.WriteRune('\n')
}

func (c *collector) addMeter(name string, m metrics.Meter) {
	c.writeGaugeCounter(name, m.Count())
}

func (c *collector) addTimer(name string, m metrics.Timer) {
	pv := []float64{0.5, 0.75, 0.95, 0.99, 0.999, 0.9999}
	ps := m.Percentiles(pv)
	c.writeSummaryCounter(name, m.Count())
	c.buff.WriteString(fmt.Sprintf(typeSummaryTpl, mutateKey(name)))
	for i := range pv {
		c.writeSummaryPercentile(name, strconv.FormatFloat(pv[i], 'f', -1, 64), ps[i])
	}
	c.buff.WriteRune('\n')
}

func (c *collector) addResettingTimer(name string, m metrics.ResettingTimer) {
	if len(m.Values()) <= 0 {
		return
	}
	ps := m.Percentiles([]float64{50, 95, 99})
	val := m.Values()
	c.writeSummaryCounter(name, len(val))
	c.buff.WriteString(fmt.Sprintf(typeSummaryTpl, mutateKey(name)))
	c.writeSummaryPercentile(name, "0.50", ps[0])
	c.writeSummaryPercentile(name, "0.95", ps[1])
	c.writeSummaryPercentile(name, "0.99", ps[2])
	c.buff.WriteRune('\n')
}

func (c *collector) writeGaugeCounter(name string, value interface{}) {
	name = mutateKey(name)
	c.buff.WriteString(fmt.Sprintf(typeGaugeTpl, name))
	c.buff.WriteString(fmt.Sprintf(keyValueTpl, name, value))
}

func (c *collector) writeSummaryCounter(name string, value interface{}) {
	name = mutateKey(name + "_count")
	c.buff.WriteString(fmt.Sprintf(typeCounterTpl, name))
	c.buff.WriteString(fmt.Sprintf(keyValueTpl, name, value))
}

func (c *collector) writeSummaryPercentile(name, p string, value interface{}) {
	name = mutateKey(name)
	c.buff.WriteString(fmt.Sprintf(keyQuantileTagValueTpl, name, p, value))
}

func mutateKey(key string) string {
	return strings.ReplaceAll(key, "/", "_")
}
