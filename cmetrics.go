package cmetrics

/*
#cgo LDFLAGS: -L/usr/local/lib -lcmetrics -lmpack -lxxhash
#cgo CFLAGS: -I/usr/local/include/ -w

#include <cmetrics/cmetrics.h>
#include <cmetrics/cmt_gauge.h>
#include <cmetrics/cmt_encode_prometheus.h>
#include <cmetrics/cmt_encode_msgpack.h>
#include <cmetrics/cmt_decode_msgpack.h>
#include <cmetrics/cmt_encode_text.h>
#include <cmetrics/cmt_encode_influx.h>
#include <cmetrics/cmt_counter.h>
*/
import "C"
import (
	"errors"
	"time"
	"unsafe"
)

type Context struct {
	context *C.struct_cmt
}

type Gauge struct {
	gauge *C.struct_cmt_gauge
}

type Counter struct {
	counter *C.struct_cmt_counter
}

func GoStringArrayToCptr(arr []string) **C.char {
	size := C.size_t(unsafe.Sizeof((*C.char)(nil)))
	length := C.size_t(len(arr))
	ptr := C.malloc(length * size)

	for i := 0; i < len(arr); i++ {
		element := (**C.char)(unsafe.Pointer(uintptr(ptr) + uintptr(i)*unsafe.Sizeof((*C.char)(nil))))
		*element = C.CString(arr[i])
	}
	return (**C.char)(ptr)
}

func (g *Gauge) Add(ts time.Time, value float64, labels []string) error {
	ret := C.cmt_gauge_add(g.gauge, C.ulong(ts.UnixNano()), C.double(value), C.int(len(labels)), GoStringArrayToCptr(labels))
	if ret != 0 {
		return errors.New("cannot substract gauge value")
	}
	return nil
}

func (g *Gauge) Inc(ts time.Time, labels []string) error {
	ret := C.cmt_gauge_inc(g.gauge, C.ulong(ts.UnixNano()), C.int(len(labels)), GoStringArrayToCptr(labels))
	if ret != 0 {
		return errors.New("cannot increment gauge value")
	}
	return nil
}

func (g *Gauge) Dec(ts time.Time, labels []string) error {
	ret := C.cmt_gauge_dec(g.gauge, C.ulong(ts.UnixNano()), C.int(len(labels)), GoStringArrayToCptr(labels))
	if ret != 0 {
		return errors.New("cannot decrement gauge value")
	}
	return nil
}

func (g *Gauge) Sub(ts time.Time, value float64, labels []string) error {
	ret := C.cmt_gauge_sub(g.gauge, C.ulong(ts.UnixNano()), C.double(value), C.int(len(labels)), GoStringArrayToCptr(labels))
	if ret != 0 {
		return errors.New("cannot subtract gauge value")
	}
	return nil
}

func (g *Gauge) GetVal(labels []string) (float64, error) {
	var value C.double
	ret := C.cmt_gauge_get_val(
		g.gauge,
		C.int(len(labels)),
		GoStringArrayToCptr(labels),
		&value)

	if ret != 0 {
		return -1, errors.New("cannot get value for gauge")
	}
	return float64(value), nil
}

func (g *Gauge) Set(ts time.Time, value float64, labels []string) error {
	ret := C.cmt_gauge_set(g.gauge, C.ulong(ts.UnixNano()), C.double(value), C.int(len(labels)), GoStringArrayToCptr(labels))
	if ret != 0 {
		return errors.New("cannot set gauge value")
	}
	return nil
}

func (ctx *Context) LabelAdd(key, val string) error {
	ret := C.cmt_labels_add_kv(ctx.context.static_labels, C.CString(key), C.CString(val))
	if ret != 0 {
		return errors.New("cannot set label to context")
	}
	return nil
}

func (ctx *Context) EncodePrometheus() (string, error) {
	ret := C.cmt_encode_prometheus_create(ctx.context, 1)
	if ret == nil {
		return "", errors.New("error encoding to prometheus format")
	}
	return C.GoString(ret), nil
}

func (ctx *Context) EncodeText() (string, error) {
	buffer := C.cmt_encode_text_create(ctx.context)
	if buffer == nil {
		return "", errors.New("error encoding to text format")
	}
	var text string = C.GoString(buffer)
	C.cmt_sds_destroy(buffer)
	return text, nil
}

func (ctx *Context) EncodeInflux() (string, error) {
	buffer := C.cmt_encode_influx_create(ctx.context)
	if buffer == nil {
		return "", errors.New("error encoding to text format")
	}
	var text string = C.GoString(buffer)
	C.cmt_encode_influx_destroy(buffer)
	return text, nil
}

func NewContextSetFromMsgPack(msgPackBuffer []byte) ([]*Context, error) {
	var ret []*Context
	var cBuffer *C.char
	ct, err := NewContext()
	if err != nil {
		return nil, err
	}
	cBuffer = (*C.char)(unsafe.Pointer(&msgPackBuffer[0]))
	var x = 0
	for x < len(msgPackBuffer) {
		r := C.cmt_decode_msgpack_create(&ct.context, cBuffer, C.ulong(len(msgPackBuffer)), (*C.ulong)(unsafe.Pointer(&x)))
		if r == 0 {
			ret = append(ret, ct)
			ct, err = NewContext()
			if err != nil {
				return nil, err
			}
		}
	}

	if len(ret) == 0 {
		return nil, errors.New("error decoding msgpack")
	}

	return ret, nil
}

func NewContextFromMsgPack(msgPackBuffer []byte, offset int64) (*Context, error) {
	var cBuffer *C.char
	var cOffset *C.ulong
	ct, err := NewContext()
	if err != nil {
		return nil, err
	}
	cBuffer = (*C.char)(unsafe.Pointer(&msgPackBuffer[0]))
	cOffset = (*C.ulong)(unsafe.Pointer(&offset))
	ret := C.cmt_decode_msgpack_create(&ct.context, cBuffer, C.ulong(len(msgPackBuffer)), cOffset)
	if ret != 0 {
		return nil, errors.New("error decoding msgpack")
	}
	return ct, nil
}

type EncoderType string

const (
	MsgPackEncoder    EncoderType = "MsgPackEncoder"
	PrometheusEncoder EncoderType = "PrometheusEncoder"
	InfluxEncoder     EncoderType = "InfluxEncoder"
	TextEncoder       EncoderType = "TextEncoder"
)

func (ctx *Context) Encode(t EncoderType) (interface{}, error) {
	switch t {
	case MsgPackEncoder:
		{
			return ctx.EncodeMsgPack()
		}
	case PrometheusEncoder:
		{
			return ctx.EncodePrometheus()
		}
	case InfluxEncoder:
		{
			return ctx.EncodeInflux()
		}
	case TextEncoder:
		{
			return ctx.EncodeText()
		}
	}
	return nil, errors.New(string("not found encoder suiteable for type " + t))
}

func (ctx *Context) EncodeMsgPack() ([]byte, error) {
	var buffer *C.char
	var bufferSize C.size_t

	ret := C.cmt_encode_msgpack_create(ctx.context, &buffer, &bufferSize)
	if ret != 0 {
		return nil, errors.New("error encoding to msgpack format")
	}
	return C.GoBytes(unsafe.Pointer(buffer), C.int(bufferSize)), nil
}

func (ctx *Context) GaugeCreate(namespace, subsystem, name, help string, labelKeys []string) (*Gauge, error) {
	gauge := C.cmt_gauge_create(ctx.context,
		C.CString(namespace),
		C.CString(subsystem),
		C.CString(name),
		C.CString(help),
		C.int(len(labelKeys)),
		GoStringArrayToCptr(labelKeys),
	)
	if gauge == nil {
		return nil, errors.New("cannot create gauge")
	}
	return &Gauge{gauge}, nil
}

func (g *Counter) Add(ts time.Time, value float64, labels []string) error {
	ret := C.cmt_counter_add(g.counter, C.ulong(ts.UnixNano()), C.double(value), C.int(len(labels)), GoStringArrayToCptr(labels))
	if ret != 0 {
		return errors.New("cannot add counter value")
	}
	return nil
}

func (g *Counter) Inc(ts time.Time, labels []string) error {
	ret := C.cmt_counter_inc(g.counter, C.ulong(ts.UnixNano()), C.int(len(labels)), GoStringArrayToCptr(labels))
	if ret != 0 {
		return errors.New("cannot Inc counter value")
	}
	return nil
}

func (g *Counter) GetVal(labels []string) (float64, error) {
	var value C.double
	ret := C.cmt_counter_get_val(
		g.counter,
		C.int(len(labels)),
		GoStringArrayToCptr(labels),
		&value)

	if ret != 0 {
		return -1, errors.New("cannot get value for counter")
	}
	return float64(value), nil
}

func (g *Counter) Set(ts time.Time, value float64, labels []string) error {
	ret := C.cmt_counter_set(g.counter, C.ulong(ts.UnixNano()), C.double(value), C.int(len(labels)), GoStringArrayToCptr(labels))
	if ret != 0 {
		return errors.New("cannot set counter value")
	}
	return nil
}

func (ctx *Context) CounterCreate(namespace, subsystem, name, help string, labelKeys []string) (*Counter, error) {
	counter := C.cmt_counter_create(ctx.context,
		C.CString(namespace),
		C.CString(subsystem),
		C.CString(name),
		C.CString(help),
		C.int(len(labelKeys)),
		GoStringArrayToCptr(labelKeys),
	)
	if counter == nil {
		return nil, errors.New("cannot create counter")
	}
	return &Counter{counter}, nil
}

func (ctx *Context) Destroy() {
	C.cmt_destroy(ctx.context)
}

func NewContext() (*Context, error) {
	cmt := C.cmt_create()
	if cmt == nil {
		return nil, errors.New("cannot create cmt context")
	}
	return &Context{context: cmt}, nil
}
