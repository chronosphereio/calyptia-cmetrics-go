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
func (suite *TestLibSuite) TestLabels() {
	context, err := NewCMTContext()
	suite.Nil(err)
	suite.NotNil(context)

	ts := time.Now()

	gauge, err := context.NewGauge("kubernetes", "network", "load", "Network load", []string{"hostname", "app"})
	suite.Nil(err)
	suite.NotNil(gauge)

	/* Default value for hash zero */
	value, err := gauge.GetValue(0, nil)
	suite.Nil(err)
	suite.Equal(0.0, value)

	/* Increment hash zero by 1 */
	err = gauge.Increment(ts, nil)
	suite.Nil(err)

	err = gauge.Add(ts, 2.0, nil)
	suite.Nil(err)

	value, err = gauge.GetValue(0, nil)
	suite.Nil(err)
	suite.Equal(3.0, value)
	/*
	 * Test 2: custom labels
	 * ---------------------
	 */
	/* Increment custom metric */
	err = gauge.Increment(ts, []string{"localhost", "cmetrics"})
	suite.Nil(err)

	value, err = gauge.GetValue(2, []string{"localhost", "cmetrics"})
	suite.Nil(err)
	suite.Equal(1.0, value)

	/* Add 10 to another metric using a different second label */
	err = gauge.Add(ts, 10, []string{"localhost", "test"})
	suite.Nil(err)

	value, err = gauge.GetValue(2, []string{"localhost", "test"})
	suite.Nil(err)
	suite.Equal(10.00, value)

	err = gauge.Subtract(ts, 2.5, []string{"localhost", "test"})
	suite.Nil(err)

	value, err = gauge.GetValue(2, []string{"localhost", "test"})
	suite.Nil(err)
	suite.Equal(7.5, value)

	encoded, err := context.PrometheusEncode()
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

	gauge, err := context.NewGauge("kubernetes", "network", "load", "Network load", []string{"hostname", "app"})
	suite.Nil(err)
	suite.NotNil(gauge)

	err = gauge.Set(time.Now(), 1, nil)
	suite.Nil(err)

	val, err := gauge.GetValue(0, nil)
	suite.Nil(err)
	suite.Equal(1.0, val)

	err = gauge.Increment(time.Now(), nil)
	suite.Nil(err)

	val, err = gauge.GetValue(0, nil)
	suite.Nil(err)
	suite.Equal(2.0, val)

	err = gauge.Subtract(time.Now(), 1, nil)
	suite.Nil(err)

	val, err = gauge.GetValue(0, nil)
	suite.Nil(err)
	suite.Equal(1.0, val)

	err = gauge.Decrement(time.Now(), nil)
	suite.Nil(err)

	val, err = gauge.GetValue(0, nil)
	suite.Nil(err)
	suite.Zero(val)

	context.Destroy()
}

func TestCMetricsBindings(t *testing.T) {
	suite.Run(t, &TestLibSuite{})
}
