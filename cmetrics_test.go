package cmetrics

import (
	"fmt"
	"github.com/influxdata/influxdb/models"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type TestLibSuite struct {
	suite.Suite
}

func (suite *TestLibSuite) TestContext() {
	context, err := NewContext()
	suite.Nil(err)
	suite.NotNil(context)
	context.Destroy()
}

func (suite *TestLibSuite) TestMultiContextFromMsgPack() {
	ts := time.Now()

	context1, err := NewContext()
	suite.Nil(err)
	suite.NotNil(context1)

	gauge, err := context1.GaugeCreate("kubernetes", "network", "load", "Network load", []string{"hostname", "app"})
	suite.Nil(err)
	suite.NotNil(gauge)

	err = gauge.Add(ts, 1.0, nil)
	suite.Nil(err)

	context2, err := NewContext()
	suite.Nil(err)
	suite.NotNil(context2)

	context3, err := NewContext()
	suite.Nil(err)
	suite.NotNil(context2)

	context4, err := NewContext()
	suite.Nil(err)
	suite.NotNil(context2)

	gauge, err = context2.GaugeCreate("kubernetes", "network", "loadss", "Network load", []string{"hostname", "app"})
	suite.Nil(err)
	suite.NotNil(gauge)

	err = gauge.Add(ts, 100.0, nil)
	suite.Nil(err)


	gauge, err = context3.GaugeCreate("kubernetes", "network", "loadss", "Network load", []string{"hostname", "app"})
	suite.Nil(err)
	suite.NotNil(gauge)

	err = gauge.Add(ts, 103.0, nil)
	suite.Nil(err)

	counter, err := context3.CounterCreate("kubernetes", "network", "loads", "Network load", []string{"hostname", "app"})
	suite.Nil(err)
	suite.NotNil(counter)

	counter, err = context3.CounterCreate("kubernetes", "network", "loads", "Network load", []string{"hostname", "app"})
	suite.Nil(err)
	suite.NotNil(counter)

	err = counter.Set(ts, 100.0, nil)
	suite.Nil(err)

	counter, err = context3.CounterCreate("kubernetes", "network", "loads", "Network load", []string{"hostname", "app"})
	suite.Nil(err)
	suite.NotNil(counter)

	err = counter.Set(ts, 101.0, nil)
	suite.Nil(err)

	buffer1, err := context1.EncodeMsgPack()
	suite.Nil(err)
	suite.NotNil(buffer1)

	buffer2, err := context2.EncodeMsgPack()
	suite.Nil(err)
	suite.NotNil(buffer2)

	buffer := append(buffer1, buffer2...)

	buffer3, err := context3.EncodeMsgPack()
	suite.Nil(err)
	suite.NotNil(buffer3)
	buffer = append(buffer, buffer3...)

	buffer4, err := context4.EncodeMsgPack()
	suite.Nil(err)
	suite.NotNil(buffer4)
	buffer = append(buffer, buffer4...)

	contextSet, err := NewContextSetFromMsgPack(buffer, 0)
	suite.Nil(err)
	suite.NotNil(contextSet)
	suite.Equal(len(contextSet), 4)

	//for _, context := range contextSet {
	//	fmt.Println(string(context.EncodeMsgPack()))
	//}

	buffer5, err := contextSet[0].EncodeMsgPack()
	suite.Nil(err)
	suite.NotNil(buffer5)

	buffer6, err := contextSet[1].EncodeMsgPack()
	suite.Nil(err)
	suite.NotNil(buffer6)

	buffer7, err := contextSet[2].EncodeMsgPack()
	suite.Nil(err)
	suite.NotNil(buffer7)

	buffer8, err := contextSet[3].EncodeMsgPack()
	suite.Nil(err)
	suite.NotNil(buffer8)

	suite.Equal(buffer5, buffer1)
	suite.Equal(buffer6, buffer2)
	suite.Equal(buffer7, buffer3)
	suite.Equal(buffer8, buffer4)

	encoded, err := contextSet[2].EncodeInflux()
	suite.Nil(err)

	points, err := models.ParsePointsString(encoded)
	suite.Nil(err)
	suite.NotNil(points)
	suite.Len(points, 3)
}

func (suite *TestLibSuite) TestGaugeLabels() {
	context, err := NewContext()
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
	context, err := NewContext()
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
	context, err := NewContext()
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

	encodedb, err := context.EncodeMsgPack()
	suite.Nil(err)
	suite.NotNil(encodedb)

	encodedBytes, err := context.EncodeText()
	suite.Nil(err)
	suite.NotNil(encodedBytes)
	encodedBytes, err = context.EncodeInflux()
	suite.Nil(err)
	suite.NotNil(encodedBytes)
}

func (suite *TestLibSuite) TestInfluxEncoding() {
	context, err := NewContext()
	suite.Nil(err)
	suite.NotNil(context)

	ts := time.Now()
	counter, err := context.CounterCreate("kubernetes", "network", "load", "Network load", []string{"hostname", "app"})
	suite.Nil(err)
	suite.NotNil(counter)

	err = counter.Set(ts, 1, nil)
	suite.Nil(err)

	var TestTags = []string{"tag1", "tag2"}

	for _, tag := range TestTags {
		err = context.LabelAdd(tag, "value")
		suite.Nil(err)
	}

	encoded, err := context.EncodeInflux()
	suite.Nil(err)

	points, err := models.ParsePointsString(encoded)
	suite.Nil(err)
	for _, point := range points {
		fields, _ := point.Fields()
		suite.Equal(1.0, fields["load"])
		for _, tag := range TestTags {
			suite.True(point.HasTag([]byte(tag)))
		}
	}
}

func (suite *TestLibSuite) TestNewContextFromMsgPack() {
	context, err := NewContext()
	suite.Nil(err)
	suite.NotNil(context)

	ts := time.Now()
	counter, err := context.CounterCreate("test1", "wired", "buffer", "a test", []string{"calyptia1", "cmetrics1"})
	suite.Nil(err)
	suite.NotNil(counter)

	err = counter.Inc(ts, nil)
	suite.Nil(err)

	gauge, err := context.GaugeCreate("test1", "wired", "buffer", "a test", []string{"calyptia2", "cmetrics2"})
	suite.Nil(err)
	suite.NotNil(gauge)

	err = gauge.Set(ts, 2, nil)
	suite.Nil(err)

	gauge, err = context.GaugeCreate("test2", "wired", "buffer", "a test", []string{"calyptia2", "cmetrics2"})
	suite.Nil(err)
	suite.NotNil(gauge)

	err = gauge.Set(ts, 1, nil)
	suite.Nil(err)

	contextInfluxEncoded, err := context.Encode(InfluxEncoder)
	suite.Nil(err)
	suite.NotEmpty(contextInfluxEncoded)

	contextTextEncoded, err := context.Encode(TextEncoder)
	suite.Nil(err)
	suite.NotEmpty(contextTextEncoded)

	contextPrometheusEncoded, err := context.Encode(PrometheusEncoder)
	suite.Nil(err)
	suite.NotEmpty(contextPrometheusEncoded)

	msgPackBuffer, err := context.Encode(MsgPackEncoder)
	suite.Nil(err)
	suite.NotEmpty(msgPackBuffer)

	newContext, err := NewContextFromMsgPack(msgPackBuffer.([]byte), 0)
	suite.Nil(err)
	suite.NotNil(newContext)

	newContextInfluxEncoded, err := newContext.Encode(InfluxEncoder)
	suite.Nil(err)
	suite.NotEmpty(newContextInfluxEncoded)
	suite.Equal(newContextInfluxEncoded, contextInfluxEncoded)

	newContextTextEncoded, err := newContext.Encode(TextEncoder)
	suite.Nil(err)
	suite.NotEmpty(newContextTextEncoded)
	suite.Equal(newContextTextEncoded, contextTextEncoded)

	newContextPrometheusEncoded, err := newContext.Encode(PrometheusEncoder)
	suite.Nil(err)
	suite.NotEmpty(newContextPrometheusEncoded)
	suite.Equal(newContextPrometheusEncoded, contextPrometheusEncoded)

	newContextMsgPackEncoded, err := newContext.Encode(MsgPackEncoder)
	suite.Nil(err)
	suite.NotEmpty(newContextMsgPackEncoded)
	suite.Equal(newContextMsgPackEncoded.([]byte), msgPackBuffer)
}

func (suite *TestLibSuite) TestCounter() {
	context, err := NewContext()
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

	encodedb, err := context.EncodeMsgPack()
	suite.Nil(err)
	suite.NotNil(encodedb)

	encodedBytes, err := context.EncodeText()
	suite.Nil(err)
	suite.NotNil(encodedBytes)

	err = context.LabelAdd("key", "value")
	suite.Nil(err)

	encodedBytes, err = context.EncodeText()
	suite.Nil(err)
	suite.Contains(encodedBytes, "key")

	encodedBytes, err = context.EncodeInflux()
	suite.Nil(err)
	suite.NotNil(encodedBytes)
	context.Destroy()
}

func TestCMetricsBindings(t *testing.T) {
	suite.Run(t, &TestLibSuite{})
}
