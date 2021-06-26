package cmetrics

/*
#cgo LDFLAGS: -L/usr/local/lib -lcmetrics -lmpack -lxxhash
#cgo CFLAGS: -I/usr/local/include/ -w

#include <cmetrics/cmetrics.h>
#include <cmetrics/cmt_gauge.h>

*/
import "C"
import (
	"fmt"
	"time"
	"unsafe"
)

type CMTContext struct {
	context *C.struct_cmt
}

type CMTGauge struct {
	gauge *C.struct_cmt_gauge
}

func GoStringArrayToCptr(arr []string) **C.char {
	var b *C.char
	ptrSize := unsafe.Sizeof(b)

	// Allocate the char** list.
	ptr := C.malloc(C.size_t(len(arr)) * C.size_t(ptrSize))

	for i := 0; i < len(arr); i++ {
		element := (**C.char)(unsafe.Pointer(uintptr(ptr) + uintptr(i)*ptrSize))
		*element = (*C.char)(unsafe.Pointer(&arr[i]))
	}

	return (**C.char)(ptr)
}

func (g *CMTGauge) Increment(ts time.Time, labelValues []string) error {
	ret := C.cmt_gauge_inc(g.gauge, C.ulong(ts.UnixNano()), C.int(len(labelValues)), GoStringArrayToCptr(labelValues))
	if ret != 0 {
		return fmt.Errorf("cannot increment gauge value")
	}
	return nil
}

func (g *CMTGauge) Decrement(ts time.Time, labelValues []string) error {
	ret := C.cmt_gauge_dec(g.gauge, C.ulong(ts.UnixNano()), C.int(len(labelValues)), GoStringArrayToCptr(labelValues))
	if ret != 0 {
		return fmt.Errorf("cannot decrement gauge value")
	}
	return nil
}


func (g *CMTGauge) Subtract(ts time.Time, value float64, labelValues []string) error {
	ret := C.cmt_gauge_sub(g.gauge, C.ulong(ts.UnixNano()),  C.double(value), C.int(len(labelValues)), GoStringArrayToCptr(labelValues))
	if ret != 0 {
		return fmt.Errorf("cannot substract gauge value")
	}
	return nil
}


func (g *CMTGauge) GetValue(labelsCount int, labelValues []string) (float64, error) {
	var  value C.double
	ret := C.cmt_gauge_get_val(
		g.gauge,
		C.int(labelsCount),
		GoStringArrayToCptr(labelValues),
		&value)

	if ret != 0 {
		return -1, fmt.Errorf("cannot get value for gauge")
	}
	return float64(value), nil
}

func (g *CMTGauge) Set(ts time.Time, value float64, labelValues []string) error {
	ret := C.cmt_gauge_set(g.gauge, C.ulong(ts.UnixNano()), C.double(value), C.int(len(labelValues)), GoStringArrayToCptr(labelValues))
	if ret != 0 {
		return fmt.Errorf("cannot set gauge value")
	}
	return nil
}

func (ctx *CMTContext) NewGauge(namespace, subsystem, name, help string, labelKeys []string) (*CMTGauge, error) {
	gauge := C.cmt_gauge_create(ctx.context,
		C.CString(namespace),
		C.CString(subsystem),
		C.CString(name),
		C.CString(help),
		C.int(len(labelKeys)),
		GoStringArrayToCptr(labelKeys),
	)
	if gauge == nil {
		return nil, fmt.Errorf("cannot create gauge")
	}
	return &CMTGauge{gauge}, nil
}

func (ctx *CMTContext) Destroy() {
	C.cmt_destroy(ctx.context)
}

func NewCMTContext() (*CMTContext, error) {
	cmt := C.cmt_create()
	if cmt == nil {
		return nil, fmt.Errorf("cannot create cmt context")
	}
	return &CMTContext{context: cmt}, nil
}