package cmetrics

import (
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
