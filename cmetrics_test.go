package cmetrics

import (
	"fmt"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type TestLibSuite struct {
	suite.Suite
}

func (suite *TestLibSuite) TestContext() {
	context, err := NewCMTContext()
	suite.Nil(err)
	suite.NotNil(context)
	context.Destroy()
}

func (suite *TestLibSuite) TestGaugeLabels() {
	context, err := NewCMTContext()
	suite.Nil(err)
	suite.NotNil(context)

	ts := time.Now()

	gauge, err := context.GaugeCreate("kubernetes", "network", "load", "Network load", []string{"hostname", "app"})
	suite.Nil(err)
	suite.NotNil(gauge)

	/* Default value for hash zero */
	_, err = gauge.GetVal(nil)
	suite.NotNil(err)

	/* Inc hash zero by 1 */
	err = gauge.Inc(ts, nil)
	suite.Nil(err)

	err = gauge.Add(ts, 2.0, nil)
	suite.Nil(err)

	value, err := gauge.GetVal(nil)
	suite.Nil(err)
	suite.Equal(3.0, value)
	/*
	 * Test 2: custom labels
	 * ---------------------
	 */
	/* Inc custom metric */
	err = gauge.Inc(ts, []string{"localhost", "cmetrics"})
	suite.Nil(err)

	value, err = gauge.GetVal([]string{"localhost", "cmetrics"})
	suite.Nil(err)
	suite.Equal(1.0, value)

	/* Add 10 to another metric using a different second label */
	err = gauge.Add(ts, 10, []string{"localhost", "test"})
	suite.Nil(err)

	value, err = gauge.GetVal([]string{"localhost", "test"})
	suite.Nil(err)
	suite.Equal(10.00, value)

	err = gauge.Sub(ts, 2.5, []string{"localhost", "test"})
	suite.Nil(err)

	value, err = gauge.GetVal([]string{"localhost", "test"})
	suite.Nil(err)
	suite.Equal(7.5, value)

	encoded, err := context.EncodePrometheus()
	suite.Nil(err)

	metricsTemplate := fmt.Sprintf(`# HELP kubernetes_network_load Network load
# TYPE kubernetes_network_load gauge
kubernetes_network_load 3 %[1]v
kubernetes_network_load{hostname="localhost",app="cmetrics"} 1 %[1]v
kubernetes_network_load{hostname="localhost",app="test"} 7.5 %[1]v
`, ts.UnixNano()/int64(time.Millisecond))

	suite.Equal(metricsTemplate, encoded)
	suite.NotNil(encoded)
}

func (suite *TestLibSuite) TestGauge() {
	context, err := NewCMTContext()
	suite.Nil(err)
	suite.NotNil(context)

	gauge, err := context.GaugeCreate("kubernetes", "network", "load", "Network load", []string{"hostname", "app"})
	suite.Nil(err)
	suite.NotNil(gauge)

	err = gauge.Set(time.Now(), 1, nil)
	suite.Nil(err)

	val, err := gauge.GetVal(nil)
	suite.Nil(err)
	suite.Equal(1.0, val)

	err = gauge.Inc(time.Now(), nil)
	suite.Nil(err)

	val, err = gauge.GetVal(nil)
	suite.Nil(err)
	suite.Equal(2.0, val)

	err = gauge.Sub(time.Now(), 1, nil)
	suite.Nil(err)

	val, err = gauge.GetVal(nil)
	suite.Nil(err)
	suite.Equal(1.0, val)

	err = gauge.Dec(time.Now(), nil)
	suite.Nil(err)

	val, err = gauge.GetVal(nil)
	suite.Nil(err)
	suite.Zero(val)

	context.Destroy()
}

func (suite *TestLibSuite) TestCounterLabels() {
	context, err := NewCMTContext()
	suite.Nil(err)
	suite.NotNil(context)

	ts := time.Now()

	counter, err := context.CounterCreate("kubernetes", "network", "load", "Network load", []string{"hostname", "app"})
	suite.Nil(err)
	suite.NotNil(counter)

	/* Default value for hash zero */
	_, err = counter.GetVal(nil)
	suite.NotNil(err)

	/* Inc hash zero by 1 */
	err = counter.Inc(ts, nil)
	suite.Nil(err)

	err = counter.Add(ts, 2.0, nil)
	suite.Nil(err)

	value, err := counter.GetVal(nil)
	suite.Nil(err)
	suite.Equal(3.0, value)
	/*
	 * Test 2: custom labels
	 * ---------------------
	 */
	/* Inc custom metric */
	err = counter.Inc(ts, []string{"localhost", "cmetrics"})
	suite.Nil(err)

	value, err = counter.GetVal([]string{"localhost", "cmetrics"})
	suite.Nil(err)
	suite.Equal(1.0, value)

	/* Add 10 to another metric using a different second label */
	err = counter.Add(ts, 10, []string{"localhost", "test"})
	suite.Nil(err)

	value, err = counter.GetVal([]string{"localhost", "test"})
	suite.Nil(err)
	suite.Equal(10.00, value)

	encoded, err := context.EncodePrometheus()
	suite.Nil(err)

	metricsTemplate := fmt.Sprintf(`# HELP kubernetes_network_load Network load
# TYPE kubernetes_network_load counter
kubernetes_network_load 3 %[1]v
kubernetes_network_load{hostname="localhost",app="cmetrics"} 1 %[1]v
kubernetes_network_load{hostname="localhost",app="test"} 10 %[1]v
`, ts.UnixNano()/int64(time.Millisecond))

	suite.Equal(metricsTemplate, encoded)
	suite.NotNil(encoded)

	encoded, err = context.EncodeMsgPack()
	suite.Nil(err)
	suite.NotNil(encoded)

	encoded, err = context.EncodeText()
	suite.Nil(err)
	suite.NotNil(encoded)

}

func (suite *TestLibSuite) TestCounter() {
	context, err := NewCMTContext()
	suite.Nil(err)
	suite.NotNil(context)

	ts := time.Now()
	counter, err := context.CounterCreate("kubernetes", "network", "load", "Network load", []string{"hostname", "app"})
	suite.Nil(err)
	suite.NotNil(counter)

	err = counter.Set(ts, 1, nil)
	suite.Nil(err)

	val, err := counter.GetVal(nil)
	suite.Nil(err)
	suite.Equal(1.0, val)

	err = counter.Inc(ts, nil)
	suite.Nil(err)

	val, err = counter.GetVal(nil)
	suite.Nil(err)
	suite.Equal(2.0, val)

	encoded, err := context.EncodePrometheus()
	suite.Nil(err)

	metricsTemplate := fmt.Sprintf(`# HELP kubernetes_network_load Network load
# TYPE kubernetes_network_load counter
kubernetes_network_load 2 %[1]v
`, ts.UnixNano()/int64(time.Millisecond))

	suite.Equal(metricsTemplate, encoded)
	suite.NotNil(encoded)

	encoded, err = context.EncodeMsgPack()
	suite.Nil(err)
	suite.NotNil(encoded)

	encoded, err = context.EncodeText()
	suite.Nil(err)
	suite.NotNil(encoded)

	context.Destroy()
}

func TestCMetricsBindings(t *testing.T) {
	suite.Run(t, &TestLibSuite{})
}
